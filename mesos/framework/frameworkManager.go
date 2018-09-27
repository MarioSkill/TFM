package framework

import (
	"encoding/json"
	"fmt"
	"gitlab.com/marioskill/configuration"
	"gitlab.com/marioskill/global"
	"gitlab.com/marioskill/mesos/framework/data"
	"gitlab.com/marioskill/mesos/framework/fcfs"
	"gitlab.com/marioskill/mesos/framework/fcfsd"
	"gitlab.com/marioskill/mesos/framework/fcfsdpriorities"
	"gitlab.com/marioskill/mesos/framework/fcfspriorities"
	"io/ioutil"
)

var conf data.FrameworkConf
var Framework interface{}

func StartScheduler(scheduler string) error {
	var name = "InitConf.json"
	raw, err := ioutil.ReadFile(configuration.FrameworksConf + "/" + scheduler + "/" + name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var c data.FrameworkConf
	json.Unmarshal(raw, &c)
	conf = c

	switch conf.Framework.Name {
	case "fcfs":
		global.SetSchedulerName("fcfs")
		fcfs.Start(conf)
	case "fcfs-d":
		global.SetSchedulerName("fcfs-d")
		fcfsd.Start(conf)
	case "fcfs-priorities":
		global.SetSchedulerName("fcfs-priorities")
		fcfspriorities.Start(conf)
	case "fcfsd-priorities":
		global.SetSchedulerName("fcfsd-priorities")
		fcfsdpriorities.Start(conf)

	}
	global.SetSchedulerStatus("Started")
	return nil
}

func Kill(framework string, id string) {
	switch framework {
	/*
		case "fcfs":
			resultado = fcfs.DesployTask(t)
		case "fcfs-d":
			resultado = fcfsd.DesployTask(t)*/
	case "fcfs-priorities":
		fcfspriorities.Kill(id)
		//case "fcfs-d-priorities":
		//	resultado = fcfsdpriorities.DesployTask(t)
	}
}

func RunTask(framework string, t data.Task) string {

	var resultado string

	switch framework {
	case "fcfs":
		resultado = fcfs.DesployTask(t)
	case "fcfs-d":
		resultado = fcfsd.DesployTask(t)
	case "fcfs-priorities":
		resultado = fcfspriorities.DesployTask(t)
	case "fcfsd-priorities":
		resultado = fcfsdpriorities.DesployTask(t)
	}
	return resultado
}

func GetTaskStatus(framework string) string {
	var resultado string
	switch framework {
	case "fcfs":
		resultado = fcfs.GetTaskStatus()
	case "fcfs-d":
		resultado = fcfsd.GetTaskStatus()
	case "fcfs-priorities":
		resultado = fcfspriorities.GetTaskStatus()
	case "fcfsd-priorities":
		resultado = fcfsdpriorities.GetTaskStatus()
	}
	return resultado

}

func StopScheduler(scheduler string) {
	switch conf.Framework.Name {
	case "fcfs":
		fcfs.Stop()
	case "fcfs-d":
		fcfsd.Stop()
	case "fcfs-priorities":
		fcfspriorities.Stop()
	case "fcfs-d-priorities":
		fcfsdpriorities.Stop()
	}

	global.SetSchedulerName("")
	global.SetSchedulerStatus("OFF")
	fmt.Println("scheduler parado")
}
