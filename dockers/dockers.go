package dockers

import (
	"encoding/json"
	"fmt"
	"gitlab.com/marioskill/configuration"
	"gitlab.com/marioskill/global"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type server struct {
	Name  string
	Value int
	Icon  string
}
type ClusterConf struct {
	BackEnd struct {
		Master struct {
			N       int `json:"N"`
			QUORUM  int `json:"QUORUM"`
			Servers []struct {
				Host string `json:"Host"`
				IP   string `json:"IP"`
			} `json:"Servers"`
		} `json:"Master"`
		Agents struct {
			N       int `json:"N"`
			Servers []struct {
				Host string `json:"Host"`
				IP   string `json:"IP"`
			} `json:"Servers"`
		} `json:"Agents"`
	} `json:"Back-End"`
	FrontEnd struct {
		URL  string `json:"URL"`
		PORT string `json:"PORT"`
		USER string `json:"USER"`
	} `json:"Front-End"`
}

type Containers struct {
	ID      string
	IMAGE   string
	COMMAND string
	CREATED string
	STATUS  string
	PORTS   string
	NAMES   string
	HOST    string
}

var ClusterFile string

func loadConf() ClusterConf {

	if ClusterFile == "" {
		ClusterFile = configuration.NodesConfs + "/" + global.GetClusterConfiguration()
	}
	raw, err := ioutil.ReadFile(ClusterFile)

	if err != nil {
		fmt.Println("No se pudo cargar la configuración (dockers.loadconf)")
		os.Exit(1)
	}

	var c ClusterConf
	json.Unmarshal(raw, &c)
	return c
}

func GetClusterURL() string {

	conf := loadConf()
	return conf.FrontEnd.URL
}
func GetClusterUser() string {
	conf := loadConf()
	return conf.FrontEnd.USER
}

func GetAgentsInfo() ClusterConf {
	conf := loadConf()
	return conf
}

func GetNumServer(s string) []server {
	var master = 0
	var agents = 0

	if s == "Running" {
		conf := loadConf()
		master = conf.BackEnd.Master.N
		agents = conf.BackEnd.Agents.N
	}

	var serverList = []server{
		server{Name: "MESOS MASTERS", Value: master, Icon: "dns"},
		server{Name: "MESOS AGENTS", Value: agents, Icon: "line_weight"},
		server{Name: "ZOOKEEPERS", Value: master, Icon: "donut_large"},
	}

	return serverList
}

func Master() map[string]string {
	conf := loadConf()
	var zk = `-e "MESOS_ZK=zk://`
	var i = 1
	for _, n := range conf.BackEnd.Master.Servers {

		if i == conf.BackEnd.Master.N || conf.BackEnd.Master.N == 1 {
			zk += n.IP + `:2181/mesos" `
		} else {
			zk += n.IP + `:2181,`
		}

		i++
	}
	var master = ""
	dmaster := make(map[string]string)
	i = 1
	for _, n := range conf.BackEnd.Master.Servers {

		master = `docker run --net="host"  --name mesos_master_` + strconv.Itoa(i) + ` -p 5050:5050 -e MESOS_CLUSTER="TUCAN" -e "MESOS_HOSTNAME=` + n.Host + `" ` +
			`-e "MESOS_IP=` + n.IP + `" ` + zk + `-e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" ` +
			`-e "MESOS_QUORUM=` + strconv.Itoa(conf.BackEnd.Master.QUORUM) + `" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" -d marioskill/mesos`
		dmaster[n.Host] = master
		//fmt.Println(master)
		i++
	}
	return dmaster
}

func Agents() map[string]string {

	conf := loadConf()
	var zk = `-e "MESOS_MASTER=zk://`
	var i = 1
	for _, n := range conf.BackEnd.Master.Servers {

		if i == conf.BackEnd.Master.N || conf.BackEnd.Master.N == 1 {
			zk += n.IP + `:2181/mesos" `
		} else {
			zk += n.IP + `:2181,`
		}

		i++
	}
	var agents = ""
	dagents := make(map[string]string)
	i = 1
	for _, n := range conf.BackEnd.Agents.Servers {

		agents = `docker run --net="host" -d -v /var/run/docker.sock:/var/run/docker.sock -v /home/mvasile/bin:/programs --name mesos_slave_` + strconv.Itoa(i) + ` ` +
			`-p 5051:5051 -e "MESOS_HOSTNAME=` + n.Host + `" ` + zk + ` -e "MESOS_LOG_DIR=/var/log/mesos" ` +
			`-e "MESOS_LOGGING_LEVEL=INFO" -e "MESOS_WORK_DIR=/var/lib/mesos" --entrypoint="/usr/local/src/mesos/1.4.1/build/bin/mesos-slave.sh" marioskill/mesos --launcher=posix --containerizers=docker,mesos`
		dagents[n.Host] = agents
		//fmt.Println(agents)
		i++
	}

	return dagents
}

func Zookeepers() map[string]string {
	conf := loadConf()
	var addServer = ""
	var i = 1
	for _, n := range conf.BackEnd.Master.Servers {
		addServer += `-e ADDITIONAL_ZOOKEEPER_` + strconv.Itoa(i) + `=server.` + strconv.Itoa(i) + `=` + n.IP + `:2888:3888 `
		i++
	}
	var zk = ""
	dzK := make(map[string]string)
	i = 1
	for _, n := range conf.BackEnd.Master.Servers {
		zk = `docker run -d --net="host" --name zk_` + strconv.Itoa(i) + ` -e SERVER_ID=` + strconv.Itoa(i) + ` ` + addServer + `marioskill/zookeeper`
		dzK[n.Host] = zk
		//fmt.Println(zk)
		i++
	}
	return dzK
}

func GetListNodes() map[string]string {
	conf := loadConf()
	cluster := make(map[string]string)
	for _, n := range conf.BackEnd.Master.Servers {
		cluster[n.Host] = n.IP

	}
	for _, n := range conf.BackEnd.Agents.Servers {
		cluster[n.Host] = n.IP

	}
	return cluster
}

func GetStatusNodes(nodes map[string]string) []Containers {
	//cmdStr := "sudo docker run -v ~/exp/a.out:/a.out ubuntu:14.04 /a.out -m 10m"
	conf := loadConf()

	//+ k + ` '` + v + `''`
	//nodes := GetListNodes()
	cmdStr := `docker ps -a --format \"table {{.ID}}\\\t{{.Image}}\\\t{{.Command}}\\\t{{.CreatedAt}}\\\t{{.Status}}\\\t\\\"{{.Ports}}\\\"\\\t{{.Names}}\"`
	var nodeList []Containers
	for host, _ := range nodes {

		re := regexp.MustCompile(`\s\s+`)
		ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + host + ` '` + cmdStr + `'`
		//fmt.Println(ssh)
		fmt.Println("get docker info from: " + host)
		out, _ := exec.Command("/bin/sh", "-c", ssh).Output()

		fmt.Println(string(out))
		//		break
		lines := strings.Split(string(out), "\n")

		for i, line := range lines {
			if i > 0 {
				j := re.Split(line, -1)
				if len(j) > 1 {
					nodeList = append(nodeList, Containers{j[0], j[1], j[2], j[3], j[4], j[5], j[6], host})
				}
			}

		}

	}

	return nodeList

}

//manage docker on cluster

func StopContainer(host string, container string) string {
	conf := loadConf()
	cmdStr := `docker stop ` + container
	ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + host + ` '` + cmdStr + `'`
	out, _ := exec.Command("/bin/sh", "-c", ssh).Output()
	return string(out)
}

func RemoveContainer(host string, container string) string {
	conf := loadConf()
	cmdStr := `docker rm ` + container
	ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + host + ` '` + cmdStr + `'`
	out, _ := exec.Command("/bin/sh", "-c", ssh).Output()
	return string(out)
}

func StopAndRemoveContainer(host string, container string) string {
	conf := loadConf()
	cmdStr := `docker stop ` + container + ` \&\& docker rm ` + container
	ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + host + ` '` + cmdStr + `'`
	out, _ := exec.Command("/bin/sh", "-c", ssh).Output()
	return string(out)
}

func StartContainer(host string, container string) string {
	conf := loadConf()
	cmdStr := `docker start ` + container
	ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + host + ` '` + cmdStr + `'`
	out, _ := exec.Command("/bin/sh", "-c", ssh).Output()
	return string(out)
}

func ShutdownCluster() string {
	conf := loadConf()
	out := ""
	for _, n := range conf.BackEnd.Master.Servers {
		cmdStr := `docker stop \$\(docker ps -a -q\) \&\& docker rm \$\(docker ps -a -q\)`
		//fmt.Println(cmdStr)
		ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + n.Host + ` '` + cmdStr + `'`
		outAux, _ := exec.Command("/bin/sh", "-c", ssh).Output()
		out += string(outAux)
	}
	for _, n := range conf.BackEnd.Agents.Servers {
		cmdStr := `docker stop \$\(docker ps -a -q\) \&\& docker rm \$\(docker ps -a -q\)`
		ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + n.Host + ` '` + cmdStr + `'`
		outAux, _ := exec.Command("/bin/sh", "-c", ssh).Output()
		//fmt.Println(cmdStr)
		out += string(outAux)
	}
	return string(out)
}

// 1 copiamos el binario desde mi PC al front del cluster
// 2 copiamos el binario al docker
func DeployApptoAgents(app string) string {

	conf := loadConf()

	//fmt.Println(conf)

	var scp = `scp ` + configuration.Apps + "/" + app + ` ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + `:` + configuration.AppClusterPath + app
	fmt.Println(scp)
	scp_out, e := exec.Command("/bin/sh", "-c", scp).Output()

	var exit = ""
	if e == nil {
		var i = 1

		for _, n := range conf.BackEnd.Agents.Servers {

			getSlaveName := `docker ps --format \"{{.Names}}\" | grep mesos`
			sshN := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + n.Host + ` ` + getSlaveName
			mesosSlave, _ := exec.Command("/bin/sh", "-c", sshN).Output()
			//fmt.Println(string(mesosSlave))
			mesosAgent := strings.Trim(string(mesosSlave), "\n")

			cmdStr := "docker cp " + configuration.AppClusterPath + app + " " + mesosAgent + ":/usr/local/bin/" + app
			ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + n.Host + ` ` + cmdStr
			i++
			out, _ := exec.Command("/bin/sh", "-c", ssh).Output()
			fmt.Println(ssh)
			exit += string(out)
		}
	}
	return exit + " " + string(scp_out)
}

func DeployApptoNode(app string, deployonlytofront bool) {
	conf := loadConf()

	//fmt.Println(conf)
	if deployonlytofront == true {

		var scp = `scp ` + configuration.Apps + "/" + app + ` ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + `:` + configuration.AppClusterPath + app
		exec.Command("/bin/sh", "-c", scp).Output()

		cmdStr := "chmod 777 " + configuration.AppClusterPath + app + " && " + configuration.AppClusterPath + app + " " + configuration.Front_port + " " + configuration.AppClusterPath
		ssh_txt := `ssh ` + configuration.ClusterUser + "@" + configuration.ClusterURL + ` '` + cmdStr + `'`
		exec.Command("/bin/sh", "-c", ssh_txt).Start()
	} else {

		var scp = `scp ` + configuration.Apps + "/" + app + ` ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + `:` + configuration.AppClusterPath + app
		_, e := exec.Command("/bin/sh", "-c", scp).Output()

		fmt.Println(scp)

		if e == nil {
			current_time := time.Now().Local()
			folder := configuration.AppClusterPath + "/" + configuration.Mediciones + "/" + current_time.Format("2006_01_02_15_04_05")

			for _, n := range conf.BackEnd.Agents.Servers {

				cmdStr := configuration.AppClusterPath + app + " " + configuration.Monitor_port + " " + folder
				ssh := `ssh ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + ` ssh ` + n.Host + ` '` + cmdStr + `'`
				fmt.Println(ssh)
				exec.Command("/bin/sh", "-c", ssh).Start()

			}
		}

	}
}
func GetBackendNodes() string {
	conf := loadConf()
	var host = "host="
	var i = 1
	for _, n := range conf.BackEnd.Agents.Servers {
		if i == conf.BackEnd.Agents.N || conf.BackEnd.Agents.N == 1 {
			host = host + n.Host
		} else {
			host = host + n.Host + ","
		}

		i++
	}

	return host
}

func GetResultsFromServer() {
	conf := loadConf()
	scp := `scp -r ` + conf.FrontEnd.USER + "@" + conf.FrontEnd.URL + `:` + configuration.AppClusterPath + configuration.Mediciones + "/ " + configuration.Results
	exec.Command("/bin/sh", "-c", scp).Output()
	fmt.Println(scp)
}

/*
func main() {
	ClusterFile = configuration.NodesConfs + "/3_5_mesos.json"
	fmt.Println(DeployApptoAgents("eidas"))
}
*/

func ClearEidas() string {
	conf := loadConf()
	exit := ""
	for _, n := range conf.BackEnd.Agents.Servers {

		eidasServers := `docker ps -a --format \"{{.Image}}#{{.Names}}\" | grep eidas2`
		sshN := `ssh ` + n.Host + ` ` + eidasServers
		eidasS, _ := exec.Command("ssh", conf.FrontEnd.USER+"@"+conf.FrontEnd.URL, sshN).Output()

		Servers := strings.Split(string(eidasS), "\n")
		fmt.Println((Servers))

		for i := range Servers {
			//fmt.Println(Servers[i], len(Servers[i]))
			if len(Servers[i]) > 0 {
				app := strings.TrimPrefix(Servers[i], "eidas2#")
				cmdStr := `'` + "docker stop " + app + " && docker rm " + app + `'`
				ssh := `ssh ` + n.Host + ` ` + cmdStr
				fmt.Println(ssh)
				out, _ := exec.Command("ssh", conf.FrontEnd.USER+"@"+conf.FrontEnd.URL, ssh).Output()

				exit += string(out)
			}
			/*

			 */
		}

	}
	return "ok"
}
