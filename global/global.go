package global

import (
	"fmt"
	"gitlab.com/marioskill/configuration"
	"gitlab.com/marioskill/mesos/api/operator"
	//"gitlab.com/marioskill/ssh"
	"io/ioutil"
	"sync"
)

var mutexClusterStatus = &sync.Mutex{}
var mutexClusterConf = &sync.Mutex{}
var mutexSchedulerStatus = &sync.Mutex{}
var mutexSchedulerName = &sync.Mutex{}

var clusterStatus string
var deployedClusterConfiguration string
var schedulerStatus string
var schedulerName string

func Init() {

	//clusterStatus = "Shutdown"
	deployedClusterConfiguration = ""
	schedulerStatus = "OFF"
	schedulerName = "null"
}

func GetClusterStatus() string {

	var s = operator.Get_health()

	if s == true {
		return "Running"
	} else {
		return "Shutdown"
	}
}

func GetClusterConfiguration() string {
	mutexClusterConf.Lock()
	if deployedClusterConfiguration == "" {
		data, e := ioutil.ReadFile(configuration.ClusetConf + "/dockerconf.txt")
		fmt.Println(string(data), e)
		mutexClusterConf.Unlock()
		return string(data)
	}
	var s = deployedClusterConfiguration
	mutexClusterConf.Unlock()
	return s
}

func GetSchedulerStatus() string {
	mutexSchedulerStatus.Lock()
	var s = schedulerStatus
	mutexSchedulerStatus.Unlock()
	return s
}

func GetSchedulerName() string {
	mutexSchedulerName.Lock()
	var s = schedulerName
	mutexSchedulerName.Unlock()
	return s
}

func SetClusterStatus(s string) {
	mutexClusterStatus.Lock()
	clusterStatus = s
	mutexClusterStatus.Unlock()
}

func SetClusterConfiguration(s string) {
	mutexClusterConf.Lock()
	deployedClusterConfiguration = s
	mydata := []byte(s)
	ioutil.WriteFile(configuration.ClusetConf+"/dockerconf.txt", mydata, 0777)
	mutexClusterConf.Unlock()

}

func SetSchedulerStatus(s string) {
	mutexSchedulerStatus.Lock()
	schedulerStatus = s
	mutexSchedulerStatus.Unlock()
}

func SetSchedulerName(s string) {
	mutexSchedulerName.Lock()
	schedulerName = s
	mutexSchedulerName.Unlock()
}
