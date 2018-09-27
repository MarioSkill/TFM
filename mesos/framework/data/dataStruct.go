package data

type FrameworkConf struct {
	Framework struct {
		Name          string `json:"Name"`
		FrameworkInfo string `json:"FrameworkInfo"`
		State         string `json:"State"`
		Master        struct {
			URL       string `json:"URL"`
			Scheduler string `json:"Scheduler"`
		} `json:"Master"`
		Conf struct {
			User            string  `json:"User"`
			AppName         string  `json:"AppName"`
			Hostname        string  `json:"hostname"`
			WEBurl          string  `json:"WEBurl"`
			WEBPort         string  `json:"WEBPort"`
			FailoverTimeout float64 `json:"FailoverTimeout"`
			Checkpoint      bool    `json:"Checkpoint"`
		} `json:"Conf"`
	} `json:"Framework"`
}

type Task struct {
	TobeCopied bool   `json:"TobeCopied"`
	Name       string `json:"Name"`
	Bin        string `json:"Bin"`
	Command    string `json:"Command"`
	Option     []struct {
		PublicPort string `json:"PublicPort"`
		Program    string `json:"Program"`
		Type       string `json:"Type"`
	} `json:"Option"`
	Instances int     `json:"instances"`
	Cpu       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Priority  int     `json:"priority"`
	QAL       int     `json:"QAL"`
	Time      string  `json:"Time"`
}

type Tests struct {
	Test struct {
		Name     string `json:"Name"`
		GetTasks struct {
			Tasks []Task `json:"Tasks"`
		} `json:"Get_Tasks"`
		URLResponse string `json:"URLResponse"`
		User        string `json:"User"`
	} `json:"Test"`
}
