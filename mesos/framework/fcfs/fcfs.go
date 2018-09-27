package fcfs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"gitlab.com/marioskill/configuration"
	loadbalancer "gitlab.com/marioskill/loadbalancer/estandar"
	"gitlab.com/marioskill/mesos/api/rpc"
	"gitlab.com/marioskill/mesos/framework/data"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type RRserver struct {
	ID   *rpc.OfferID
	AID  *rpc.AgentID
	Host string
	IP   string
}

var currentServer = 0 // indica la posición del servicdor actual

var cmdMutex = &sync.Mutex{} //mutex para los comandos
var appMutex = &sync.Mutex{} // para distribuir apps

var master string //"http://localhost:5050"
var path string   //"/api/v1/scheduler"

var frameworkInfoFile string
var stateFile string

var marshaller = jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      "  ",
	OrigName:    true,
}
var mesosStreamID string
var frameworkInfo rpc.FrameworkInfo
var commandChan = make(chan string, 100)
var taskId uint64
var tasksState = make(map[string]*rpc.TaskStatus)
var Runing = true

var frameworkID string
var Tasks []data.Task

const sendResults = "http://10.0.10.1:6060/eidasResult"
const sendTimes = "http://10.0.10.1:6060/TimeClient"

func Start(frameworkConf data.FrameworkConf) {
	taskId = 0
	log.Println("[FCFS] ----------> SCHEDULER LOADING < ---------- [FCFS]")
	master = frameworkConf.Framework.Master.URL     // url al frontend
	path = frameworkConf.Framework.Master.Scheduler //url de la api del planificador

	//Ficheros con la configuracón previa de planificador asi como su estado
	frameworkInfoFile = configuration.FrameworksConf + "/fcfs/" + frameworkConf.Framework.FrameworkInfo
	stateFile = configuration.FrameworksConf + "/fcfs/" + frameworkConf.Framework.State

	log.Printf("[FCFS] Try to load previous FrameworkInfo from %s", frameworkInfoFile)
	frameworkJson, err := ioutil.ReadFile(frameworkInfoFile)

	//log.Print(string(frameworkJson))

	if err == nil {
		err := jsonpb.UnmarshalString(string(frameworkJson), &frameworkInfo)
		if err != nil {
			log.Printf("[FCFS] Error %v. Please delete %s and restart", err, frameworkInfoFile)
		}
	} else {
		log.Printf("[FCFS] Fallback to defaults due to error [%v]", err)
		frameworkInfo = rpc.FrameworkInfo{
			User:            &frameworkConf.Framework.Conf.User,
			Name:            &frameworkConf.Framework.Conf.AppName,
			Hostname:        &frameworkConf.Framework.Conf.Hostname,
			WebuiUrl:        &frameworkConf.Framework.Conf.WEBurl,
			FailoverTimeout: &frameworkConf.Framework.Conf.FailoverTimeout,
			Checkpoint:      &frameworkConf.Framework.Conf.Checkpoint,
		}
	}
	Runing = true

	sub, errorsub := subscribe()
	if errorsub == nil {
		go func() {
			//log.Println()
			scheduling(sub)
		}()

	} else {
		log.Printf("[FCFS] Problems to subreibe against mesos")
		log.Println(errorsub)
	}

}

func subscribe() (*http.Response, error) {
	subscribeCall := &rpc.Call{
		FrameworkId: frameworkInfo.Id,
		Type:        rpc.Call_SUBSCRIBE.Enum(),
		Subscribe:   &rpc.Call_Subscribe{FrameworkInfo: &frameworkInfo},
	}

	body, err := marshaller.MarshalToString(subscribeCall)

	if err != nil {
		log.Printf("SUBCRIBE-> %v Error: %v", log.Llongfile, err)
		return nil, err
	}
	//log.Print(body)
	res, err := http.Post(master+path, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Printf("SUBCRIBE-2-> %v Error: %v", log.Llongfile, err)
		return nil, err
	}

	return res, nil
}

func scheduling(res *http.Response) error {

	defer res.Body.Close()
	reader := bufio.NewReader(res.Body)
	line, err := reader.ReadString('\n')
	bytesCount, err := strconv.Atoi(strings.Trim(line, "\n"))
	if err != nil {
		log.Printf("SCHEDULING->%v Error: %v", log.Llongfile, err)
		return err
	}

	for Runing {
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, "\n")
		data := line[:bytesCount]
		bytesCount, err = strconv.Atoi((line[bytesCount:]))

		if err != nil {
			log.Printf("SCHEDULING-2->%v Error: %v", log.Llongfile, err)
			return err
		}
		var sub rpc.Event
		json.Unmarshal([]byte(data), &sub)

		switch *sub.Type {
		case rpc.Event_SUBSCRIBED:

			frameworkInfo.Id = sub.Subscribed.FrameworkId
			mesosStreamID = res.Header.Get("Mesos-Stream-Id")
			json, err := marshaller.MarshalToString(&frameworkInfo)
			if err != nil {
				log.Printf("SCHEDULING-3->")
				return err
			}
			ioutil.WriteFile(frameworkInfoFile, []byte(json), 0644)

			//reconcile()
		case rpc.Event_HEARTBEAT:
			log.Print("PING")
		case rpc.Event_OFFERS:
			log.Printf("Handle OFFERS returns")
			handleOffers(sub.Offers)
		case rpc.Event_UPDATE:
			log.Printf("Handle UPDATE returns")
			handleUpdate(sub.Update)

		}

	}
	return nil
}

func reconcile() {
	oldState, err := ioutil.ReadFile(stateFile)
	if err == nil {
		err := json.Unmarshal(oldState, &tasksState)
		if err != nil {
			log.Printf("RECONCILE -> Error on loading previous state %v", err)
		}
	}

	oldTasks := make([]*rpc.Call_Reconcile_Task, 0)
	maxId := 0
	for _, t := range tasksState {
		oldTasks = append(oldTasks, &rpc.Call_Reconcile_Task{
			TaskId:  t.TaskId,
			AgentId: t.AgentId,
		})
		numericId, err := strconv.Atoi(t.TaskId.GetValue())
		if err == nil && numericId > maxId {
			maxId = numericId
		}
	}
	atomic.StoreUint64(&taskId, uint64(maxId))
	call(&rpc.Call{
		Type:      rpc.Call_RECONCILE.Enum(),
		Reconcile: &rpc.Call_Reconcile{Tasks: oldTasks},
	})
}

func kill(id string) error {
	update, ok := tasksState[id]
	log.Printf("Kill task %s [%#v]", id, update)
	if !ok {
		return fmt.Errorf("KILL -> Unknown task %s", id)
	}
	return call(&rpc.Call{
		Type: rpc.Call_KILL.Enum(),
		Kill: &rpc.Call_Kill{
			TaskId:  update.TaskId,
			AgentId: update.AgentId,
		},
	})
}

func handleUpdate(update *rpc.Event_Update) error {
	tasksState[update.Status.TaskId.GetValue()] = update.Status
	stateJson, _ := json.Marshal(tasksState)
	ioutil.WriteFile(stateFile, stateJson, 0644)

	return call(&rpc.Call{
		Type: rpc.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &rpc.Call_Acknowledge{
			AgentId: update.Status.AgentId,
			TaskId:  update.Status.TaskId,
			Uuid:    update.Status.Uuid,
		},
	})
}

func getTaskMaxPriority() int {
	priority := -1
	idTask := 0
	for i := 0; i < len(Tasks); i++ {
		if Tasks[i].Priority > priority {
			idTask = i
			priority = Tasks[i].Priority
		}
	}
	return idTask
}

func Difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
func DesployTask(t data.Task) string {
	fmt.Println("[NEW TASK] --> " + t.Name)

	cmdMutex.Lock() //CONTROLAMOS TASK
	t.Time = time.Now().Format("15:04:05")
	Tasks = append(Tasks, t)
	cmdMutex.Unlock()

	//cmd := r.Form["cmd"][0]
	commandChan <- "RUN"
	return t.Command
}
func handleOffers(offers *rpc.Event_Offers) error {

	select {
	case <-commandChan: //Llega una peticion

		//Buscar elemento con mayor prioridad

		var UsedServers []int
		//var OfferedServers []int

		cmdMutex.Lock()
		priority := getTaskMaxPriority()
		var tasklen = len(Tasks)
		cmdMutex.Unlock()

		if priority > -1 && tasklen > 0 {

			cmdMutex.Lock()
			t := Tasks[priority]
			Tasks = Tasks[:priority+copy(Tasks[priority:], Tasks[priority+1:])]
			cmdMutex.Unlock()

			/////appMutex.Lock()
			appMutex.Lock()
			var selectedServer []RRserver
			//var j = 0

			for _, offer := range offers.Offers {
				selectedServer = append(selectedServer, RRserver{ID: offer.Id, AID: offer.AgentId, Host: *offer.Url.Address.Hostname, IP: *offer.Url.Address.Ip})
				//OfferedServers = append(OfferedServers, j)
				//j++
			}
			//currentServer = -1
			for i := 0; i < len(selectedServer) && i < t.Instances; i++ {

				offerIdsDeploy := []*rpc.OfferID{}
				fmt.Println(currentServer, len(selectedServer))
				if currentServer+1 >= len(selectedServer) {
					currentServer = 0
				} else {
					currentServer++
				}
				fmt.Println(currentServer)
				//	UsedServers = append(UsedServers, selectedServer[currentServer].ID)
				UsedServers = append(UsedServers, currentServer)

				data := selectedServer[currentServer]
				offerIdsDeploy = append(offerIdsDeploy, data.ID)

				myAgentId := data.AID
				TRUE := true
				newTaskId := fmt.Sprint(atomic.AddUint64(&taskId, 1))
				Portaux, _ := strconv.Atoi(newTaskId)
				localCMD := t.Command

				if t.Name == "eidasServer" {
					var porServer = strconv.Itoa(8000 + Portaux)

					loadbalancer.SubServer(newTaskId, data.Host, data.IP, "2.0", porServer, t.Priority, t.QAL)

					localCMD = strings.Replace(localCMD, "#PublicPort", porServer, -1)
					localCMD = strings.Replace(localCMD, "#Bin", t.Bin, -1)
					localCMD = strings.Replace(localCMD, "#Program", t.Option[0].Program, -1)
					localCMD = strings.Replace(localCMD, "#IPServer", data.IP, -1)
					localCMD = strings.Replace(localCMD, "#PortServer", porServer, -1)

				} else if t.Name == "eidasClient" {
					subcrition := loadbalancer.SubClient(strconv.Itoa(Portaux), data.Host, data.IP, "1", t.Priority, t.QAL)
					localCMD = localCMD + " " + subcrition.Host + "_" + strconv.Itoa(Portaux) + " http://" + subcrition.Host + " " + subcrition.Port + " " + sendResults + " " + strconv.FormatFloat(subcrition.Lambda, 'E', -1, 64) + " " + strconv.Itoa(subcrition.Nclients) + " " + subcrition.Timestamp.Format("15:04:05") + " " + sendTimes + " " + t.Time
				}
				var nameAPP = t.Name + t.Option[0].Type
				accept := &rpc.Call{
					Type: rpc.Call_ACCEPT.Enum(),
					Accept: &rpc.Call_Accept{
						OfferIds: offerIdsDeploy,
						Operations: []*rpc.Offer_Operation{{
							Type: rpc.Offer_Operation_LAUNCH.Enum(),
							Launch: &rpc.Offer_Operation_Launch{TaskInfos: []*rpc.TaskInfo{{
								Name:      &nameAPP,
								TaskId:    &rpc.TaskID{Value: &newTaskId},
								AgentId:   myAgentId,
								Resources: getResources(t.Cpu, t.Mem),
								Command: &rpc.CommandInfo{
									Shell: &TRUE,
									Value: &localCMD,
								},
							}},
							},
						},
						},
					},
				}

				err := call(accept)
				if err != nil {
					fmt.Println("tarea no aceptada")
					//return err
				}
				//localCMD = cmd

			}
			var aux = currentServer
			for i := t.Instances; i < len(selectedServer); i++ {

				if aux+1 >= len(selectedServer) {
					aux = 0
				} else {
					aux++
				}
				newTaskId := fmt.Sprint(atomic.AddUint64(&taskId, 1))
				var c = ""
				var nameAPP = "SYNC#RR"
				offerIdsDeploy := []*rpc.OfferID{}
				data := selectedServer[aux]
				offerIdsDeploy = append(offerIdsDeploy, data.ID)
				TRUE := true
				myAgentId := data.AID
				accept := &rpc.Call{
					Type: rpc.Call_ACCEPT.Enum(),
					Accept: &rpc.Call_Accept{
						OfferIds: offerIdsDeploy,
						Operations: []*rpc.Offer_Operation{{
							Type: rpc.Offer_Operation_LAUNCH.Enum(),
							Launch: &rpc.Offer_Operation_Launch{TaskInfos: []*rpc.TaskInfo{{
								Name:      &nameAPP,
								TaskId:    &rpc.TaskID{Value: &newTaskId},
								AgentId:   myAgentId,
								Resources: getResources(1.0, 1.0),
								Command: &rpc.CommandInfo{
									Shell: &TRUE,
									Value: &c,
								},
							}},
							},
						},
						},
					},
				}

				err := call(accept)
				if err != nil {
					fmt.Println("tarea no aceptada")
					//return err
				}
			}
			/*
				offerIds := []*rpc.OfferID{}
				for _, offer := range offers.Offers {
					offerIds = append(offerIds, offer.Id)
					//break
				}

				decline := &rpc.Call{
					Type:    rpc.Call_DECLINE.Enum(),
					Decline: &rpc.Call_Decline{OfferIds: offerIds},
				}
				err := call(decline)

				if err != nil {
					fmt.Println(err)
				}*/
			t.Instances = t.Instances - len(selectedServer)
			appMutex.Unlock()

			cmdMutex.Lock()

			numTask := 0
			if t.Instances > 0 {
				Tasks = append(Tasks, t)
				numTask = 1
			} else {
				numTask = len(Tasks)
			}
			cmdMutex.Unlock()

			if numTask > 0 {
				commandChan <- "RUN"
			}
		}

	default:
		offerIds := []*rpc.OfferID{}
		for _, offer := range offers.Offers {
			offerIds = append(offerIds, offer.Id)
			//break
		}

		decline := &rpc.Call{
			Type:    rpc.Call_DECLINE.Enum(),
			Decline: &rpc.Call_Decline{OfferIds: offerIds},
		}
		err := call(decline)

		if err != nil {
			fmt.Println(err)
		}

	}
	//appMutex.UnLock()
	return nil
}

//cmd string, cpu float64, mem float64, name string, bin string, instances int, priority int

func getResources(cpu float64, mem float64) []*rpc.Resource {
	CPU := "cpus"
	MEM := "mem"
	//cpu := float64(0.1)
	return []*rpc.Resource{
		{
			Name:   &CPU,
			Type:   rpc.Value_SCALAR.Enum(),
			Scalar: &rpc.Value_Scalar{Value: &cpu},
		},
		{
			Name:   &MEM,
			Type:   rpc.Value_SCALAR.Enum(),
			Scalar: &rpc.Value_Scalar{Value: &mem},
		},
	}

}

// call Sends request to Mesos.
// Returns nil if request was accepted in other case returns error
func call(message *rpc.Call) error {
	message.FrameworkId = frameworkInfo.Id
	body, err := marshaller.MarshalToString(message)
	if err != nil {
		log.Printf("CALL -> %v Error: %v", log.Llongfile, err)
		return err
	}
	req, err := http.NewRequest("POST", master+path, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Mesos-Stream-Id", mesosStreamID)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Call %s %s", message.Type, string(body))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("CALL2-> %v Error: %v", log.Llongfile, err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 202 {
		io.Copy(os.Stderr, res.Body)
		return fmt.Errorf("CALL3-> Error %d", res.StatusCode)
	}

	return nil
}

func Stop() {
	Runing = false
}

func GetFrameworkID() string {
	return frameworkID
}

func GetTaskStatus() string {
	stateJson, _ := json.Marshal(tasksState) //Can't use proto.Marshal because there is no definition for map

	return string(stateJson)
}
