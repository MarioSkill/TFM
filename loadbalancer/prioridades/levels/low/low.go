package low

import (
	"container/list"
	"fmt"
	"gitlab.com/marioskill/configuration"
	"gitlab.com/marioskill/global"
	"gitlab.com/marioskill/loadbalancer/prioridades/levels"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	ID            string `json:"Id"`
	Host          string `json:"Host"`
	IP            string `json:"IP"`
	Version       string `json:"Version"`
	Port          string `json:"Port"`
	Priority      int    `json:"Priority"`
	QAL           int    `json:"QAL"`
	LastClientSub time.Time
	Heartbeat     int
}
type Client struct {
	ID       string    `json:"Id"`
	Host     string    `json:"Host"`
	IP       string    `json:"IP"`
	Priority int       `json:"Priority"`
	QAL      int       `json:"QAL"`
	TIME     time.Time `json:"Time"`
}

var Servers []Server
var Clients []Client
var nServers = 0
var nClients = 0
var maxClientsPerSeconds = 64 - 10 // 64 clientes concurrentes - 10 para tener tiempo de levantar servidor
var maxClientsSupported = nServers * maxClientsPerSeconds
var clientsPerSubcription = 25          // para poder causar sobrecarga, cada registro de cliente se va a considerar como 20 clientes
var averageTimeClient = 4 * time.Second // segundos

var index *list.Element //Indice: servidor atiende que cliente
var roundRobin = list.New()

var mTick = &sync.Mutex{}

func RoundRobin() {
	mTick.Lock()
	var e = index.Next()
	if e == nil {
		index = roundRobin.Front()
	} else {
		index = e
	}
	mTick.Unlock()
}

func SubClient(id string, host string, ip string, version string, p int, qal int) levels.ServerInformation {
	mTick.Lock()
	nClients = (len(Clients) + 1) * clientsPerSubcription
	//fmt.Println(">> Clientes >>", nClients)
	if nClients > maxClientsSupported {
		//fmt.Println(nServers)
		//mTick.Unlock()
		//SubServer("1", "localhost", "111", "22", "8080")
		//mTick.Lock()
		/*
			raw, err := ioutil.ReadFile(configuration.LoafBalancerServers + "/eIDAS_Server_lb.json")
			if err != nil {
				fmt.Println(err.Error())

				//return err
			}*/
		data := url.Values{}
		//	configuration.GUI_port = configuration.GUI_port
		data.Set("testConf", "eIDAS_Server_lb_low.json")

		switch global.GetSchedulerName() {
		case "fcfs-priorities":
			sendResult("http://localhost:"+configuration.GUI_port+"/schedulers/fcfs-priorities/balancer", data)
		case "fcfsd-priorities":
			fmt.Printf("LB2")
			sendResult("http://localhost:"+configuration.GUI_port+"/schedulers/fcfsd-priorities/balancer", data)
		case "priorities":
			//	sendResult("http://localhost:"+configuration.GUI_port+"/schedulers/fcfs-priorities/balancer", data)
		}

	}

	data := Client{ID: id, Host: host, IP: ip, TIME: time.Now(), Priority: p, QAL: qal}
	Clients = append(Clients, data)
	if nServers > 0 {
		mTick.Unlock()
		RoundRobin()
		mTick.Lock()
		var sInfo = index.Value.(Server)
		aux := levels.ServerInformation{ID: sInfo.ID, Host: sInfo.Host, Port: sInfo.Port, Lambda: 0.1, Timestamp: time.Now(), Nclients: clientsPerSubcription}
		mTick.Unlock()
		return aux
	} else {
		mTick.Unlock()
		return levels.ServerInformation{}
	}

}

func SubServer(id string, host string, ip string, version string, port string, p int, qal int) {
	mTick.Lock()
	nServers++
	data := Server{ID: id, Host: host, IP: ip, Version: version, Port: port, Heartbeat: 0, LastClientSub: time.Time{}, Priority: p, QAL: qal}
	index = roundRobin.PushBack(data)
	maxClientsSupported = nServers * maxClientsPerSeconds
	mTick.Unlock()
	//fmt.Println(">> ", maxClientsSupported)
}

func deregisterClients() {
	for {
		mTick.Lock()
		for i := 0; i < len(Clients); i++ {
			if time.Now().Sub(Clients[i].TIME) > averageTimeClient {
				//fmt.Println(time.Now().Sub(Clients[i].TIME))
				Clients = Clients[:i+copy(Clients[i:], Clients[i+1:])]
			}
		}
		nClients = len(Clients) * clientsPerSubcription
		mTick.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func sendResult(resource string, aux interface{}) {

	data := aux.(url.Values)
	//var u *url.URL
	//u, _ = url.ParseRequestURI(resource)
	urlStr := resource //u.String() // 'https://urlServer'

	//fmt.Println(urlStr, data)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client.Do(r)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println((resp))

}
