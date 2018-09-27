package operator

//communication between operators and Mesos master/agent.
import (
	"encoding/json"
	"fmt"
	"gitlab.com/marioskill/request"
)

func Get_health() bool {
	query := request.MesosQuery{Type: "GET_HEALTH"}

	aux, err := request.Do(query)
	if err == nil {
		var result request.GET_HEALTH
		json.Unmarshal(aux.Response, &result)
		fmt.Println(aux.Http_code)
		return result.GetHealth.Healthy
	} else {
		return false
	}
}

func Get_flags() request.GET_FLAGS {
	query := request.MesosQuery{Type: "GET_FLAGS"}

	aux, _ := request.Do(query)
	var result request.GET_FLAGS
	json.Unmarshal(aux.Response, &result)
	return result

}

func Get_version() request.GET_VERSION {
	query := request.MesosQuery{Type: "GET_VERSION"}

	aux, _ := request.Do(query)
	var result request.GET_VERSION
	json.Unmarshal(aux.Response, &result)

	//fmt.Println("GET_VERSION")
	//fmt.Println(aux.Http_code)
	//fmt.Println(result)
	return result
}

func Get_metric() request.GET_METRICS {
	query := request.MesosQuery2{
		Type: "GET_METRICS",
	}
	query.GetMetrics.Timeout.Nanoseconds = 5000000000

	aux, _ := request.Do(query)
	var result request.GET_METRICS
	json.Unmarshal(aux.Response, &result)

	return result
}

func Get_logginglevel() {
	query := request.MesosQuery{Type: "GET_LOGGING_LEVEL"}

	aux, _ := request.Do(query)
	var result request.GET_LOGGING_LEVEL
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_LOGGING_LEVEL")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

/***
(by default it’s 0, libprocess uses levels 1, 2, and 3).
*/
func Set_logginglevel(level int) {
	query := request.SET_LOGGING_LEVEL{Type: "SET_LOGGING_LEVEL"}

	query.SetLoggingLevel.Level = level
	aux, _ := request.Do(query)
	var result request.GET_LOGGING_LEVEL
	json.Unmarshal(aux.Response, &result)

	fmt.Println("SET_LOGGING_LEVEL")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func List_files() request.LIST_FILES {

	in := []byte(`{ "type": "LIST_FILES", "list_files": { "path":"/" } }`)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)

	aux, _ := request.Do(raw)
	var result request.LIST_FILES
	json.Unmarshal(aux.Response, &result)
	return result
}

func Read_file() {

	in := []byte(`{
		  "type": "READ_FILE",
		  "read_file": {
		    "length": 6,
		    "offset": 1,
		    "path": "myname"
		  }
		}`)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)

	aux, _ := request.Do(raw)
	var result request.READ_FILE
	json.Unmarshal(aux.Response, &result)

	fmt.Println("READ_FILE")
	fmt.Println(aux.Http_code)
	fmt.Println(aux)
}

func Get_state() {
	query := request.MesosQuery{Type: "GET_STATE"}
	aux, _ := request.Do(query)
	var result request.GET_STATE
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_STATE")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_agents() request.GET_AGENTS {
	query := request.MesosQuery{Type: "GET_AGENTS"}
	aux, err := request.Do(query)
	var result request.GET_AGENTS
	if err == nil {

		json.Unmarshal(aux.Response, &result)
	}
	return result
}

func Get_frameworks() {
	query := request.MesosQuery{Type: "GET_FRAMEWORKS"}
	aux, _ := request.Do(query)
	var result request.GET_FRAMEWORKS
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_FRAMEWORKS")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_executors() {
	query := request.MesosQuery{Type: "GET_EXECUTORS"}
	aux, _ := request.Do(query)
	var result request.GET_EXECUTORS
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_EXECUTORS")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_tasks() {
	query := request.MesosQuery{Type: "GET_TASKS"}
	aux, _ := request.Do(query)
	var result request.GET_TASKS
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_TASKS")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_roles() {
	query := request.MesosQuery{Type: "GET_ROLES"}
	aux, _ := request.Do(query)
	var result request.GET_ROLES
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_ROLES")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_weights() {
	query := request.MesosQuery{Type: "GET_WEIGHTS"}
	aux, _ := request.Do(query)
	var result request.GET_WEIGHTS
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_WEIGHTS")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Update_weights(role string, weight string) {

	in := []byte(`{
		  "type": "UPDATE_WEIGHTS",
		  "update_weights": {
		    "weight_infos": [
		      {
		        "role": ` + role + `,
		        "weight":` + (weight) + `
		      }
		    ]
		  }
		}`)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)

	aux, _ := request.Do(raw)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_WEIGHTS")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_master() {
	query := request.MesosQuery{Type: "GET_MASTER"}
	aux, _ := request.Do(query)
	var result request.GET_MASTER
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_MASTER")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Reserve_resource(query request.RESERVE_RESOURCES) {

	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("RESERVE_RESOURCES")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Unreserve_resource(query request.UNRESERVE_RESOURCES) {

	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("UNRESERVE_RESOURCES")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Create_volumes(query request.CREATE_VOLUMES) {

	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("CREATE_VOLUMES")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Destroy_volumes(query request.DESTROY_VOLUMES) {

	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("DESTROY_VOLUMES")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Grow_volumes(query request.GROW_VOLUME) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GROW_VOLUME")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Shrink_volumes(query request.SHRINK_VOLUME) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("SHRINK_VOLUME")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_maintenance_schedule() {
	query := request.MesosQuery{Type: "GET_MAINTENANCE_SCHEDULE"}
	aux, _ := request.Do(query)
	var result request.GET_MAINTENANCE_SCHEDULE
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_MAINTENANCE_SCHEDULE")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Update_maintenance_schedule(query request.UPDATE_MAINTENANCE_SCHEDULE) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("UPDATE_MAINTENANCE_SCHEDULE")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Start_maintenance(query request.START_MAINTENANCE) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("START_MAINTENANCE")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Stop_maintenance(query request.STOP_MAINTENANCE) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("STOP_MAINTENANCE")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Get_cuota() {
	query := request.MesosQuery{Type: "GET_QUOTA"}
	aux, _ := request.Do(query)
	var result request.GET_QUOTA
	json.Unmarshal(aux.Response, &result)

	fmt.Println("GET_QUOTA")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Set_cuota(query request.SET_QUOTA) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("SET_QUOTA")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

func Remove_cuota(query request.REMOVE_QUOTA) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("REMOVE_QUOTA")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

/**
This call can be used by operators to assert that an agent instance has failed and is never coming back
(e.g., ephemeral instance from cloud provider). The master would shutdown the agent and send TASK_GONE_BY_OPERATOR updates
for all the running tasks. This signal can be used by stateful frameworks to re-schedule their workloads
 (volumes, reservations etc.) to other agent instances. It is possible that the tasks might still be running if the operator’s assertion
  was wrong and the agent was partitioned away from the master. The agent would be shutdown when it tries to reregister with
  the master when the partition heals. This call is idempotent.
*/

func Mark_agent_gone(query request.MARK_AGENT_GONE) {
	aux, _ := request.Do(query)
	var result interface{}
	json.Unmarshal(aux.Response, &result)

	fmt.Println("MARK_AGENT_GONE")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}

/*
EVENTS to-do
func Subscribe() {
	query := request.MesosQuery{Type: "SUBSCRIBE"}
	aux,_ := request.Do(query)
	var result request.SUBSCRIBE
	json.Unmarshal(aux.Response, &result)

	fmt.Println("SUBSCRIBE")
	fmt.Println(aux.Http_code)
	fmt.Println(result)
}
*/
