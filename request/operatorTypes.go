package request

//const SERVER = "http://localhost:5050/api/v1"

//var SERVER = ""

type MesosQuery struct {
	Type string `json:"type"`
}

type MesosQuery2 struct {
	Type       string `json:"type"`
	GetMetrics struct {
		Timeout struct {
			Nanoseconds int64 `json:"nanoseconds"`
		} `json:"timeout"`
	} `json:"get_metrics"`
}

type Result struct {
	Http_code int
	Http_msg  string
	Response  []byte
}

type GET_HEALTH struct {
	Type      string `json:"type"`
	GetHealth struct {
		Healthy bool `json:"healthy"`
	} `json:"get_health"`
}

type GET_FLAGS struct {
	Type     string `json:"type"`
	GetFlags struct {
		Flags []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"flags"`
	} `json:"get_flags"`
}

type GET_VERSION struct {
	Type       string `json:"type"`
	GetVersion struct {
		VersionInfo struct {
			Version   string `json:"version"`
			BuildDate string `json:"build_date"`
			BuildTime int    `json:"build_time"`
			BuildUser string `json:"build_user"`
		} `json:"version_info"`
	} `json:"get_version"`
}

type GET_METRICS struct {
	Type       string `json:"type"`
	GetMetrics struct {
		Metrics []struct {
			Name  string  `json:"name"`
			Value float64 `json:"value"`
		} `json:"metrics"`
	} `json:"get_metrics"`
}

type GET_LOGGING_LEVEL struct {
	Type            string `json:"type"`
	GetLoggingLevel struct {
		Level int `json:"level"`
	} `json:"get_logging_level"`
}

type SET_LOGGING_LEVEL struct {
	Type            string `json:"type"`
	SetLoggingLevel struct {
		Duration struct {
			Nanoseconds int64 `json:"nanoseconds"`
		} `json:"duration"`
		Level int `json:"level"`
	} `json:"set_logging_level"`
}

type LIST_FILES struct {
	Type      string `json:"type"`
	ListFiles struct {
		FileInfos []struct {
			Gid   string `json:"gid"`
			Mode  int    `json:"mode"`
			Mtime struct {
				Nanoseconds int64 `json:"nanoseconds"`
			} `json:"mtime"`
			Nlink int    `json:"nlink"`
			Path  string `json:"path"`
			Size  int    `json:"size"`
			UID   string `json:"uid"`
		} `json:"file_infos"`
	} `json:"list_files"`
}

type READ_FILE struct {
	Type     string `json:"type"`
	ReadFile struct {
		Data string `json:"data"`
		Size int    `json:"size"`
	} `json:"read_file"`
}

type GET_STATE struct {
	Type     string `json:"type"`
	GetState struct {
		GetAgents struct {
			Agents []struct {
				Active    bool `json:"active"`
				AgentInfo struct {
					Hostname string `json:"hostname"`
					ID       struct {
						Value string `json:"value"`
					} `json:"id"`
					Port      int `json:"port"`
					Resources []struct {
						Name   string `json:"name"`
						Role   string `json:"role"`
						Scalar struct {
							Value float64 `json:"value"`
						} `json:"scalar,omitempty"`
						Type   string `json:"type"`
						Ranges struct {
							Range []struct {
								Begin int `json:"begin"`
								End   int `json:"end"`
							} `json:"range"`
						} `json:"ranges,omitempty"`
					} `json:"resources"`
				} `json:"agent_info"`
				Pid            string `json:"pid"`
				RegisteredTime struct {
					Nanoseconds int64 `json:"nanoseconds"`
				} `json:"registered_time"`
				TotalResources []struct {
					Name   string `json:"name"`
					Role   string `json:"role"`
					Scalar struct {
						Value float64 `json:"value"`
					} `json:"scalar,omitempty"`
					Type   string `json:"type"`
					Ranges struct {
						Range []struct {
							Begin int `json:"begin"`
							End   int `json:"end"`
						} `json:"range"`
					} `json:"ranges,omitempty"`
				} `json:"total_resources"`
				Version string `json:"version"`
			} `json:"agents"`
		} `json:"get_agents"`
		GetExecutors struct {
			Executors []struct {
				AgentID struct {
					Value string `json:"value"`
				} `json:"agent_id"`
				ExecutorInfo struct {
					Command struct {
						Shell bool   `json:"shell"`
						Value string `json:"value"`
					} `json:"command"`
					ExecutorID struct {
						Value string `json:"value"`
					} `json:"executor_id"`
					FrameworkID struct {
						Value string `json:"value"`
					} `json:"framework_id"`
				} `json:"executor_info"`
			} `json:"executors"`
		} `json:"get_executors"`
		GetFrameworks struct {
			Frameworks []struct {
				Active        bool `json:"active"`
				Connected     bool `json:"connected"`
				FrameworkInfo struct {
					Checkpoint      bool    `json:"checkpoint"`
					FailoverTimeout float64 `json:"failover_timeout"`
					Hostname        string  `json:"hostname"`
					ID              struct {
						Value string `json:"value"`
					} `json:"id"`
					Name      string `json:"name"`
					Principal string `json:"principal"`
					Role      string `json:"role"`
					User      string `json:"user"`
				} `json:"framework_info"`
				RegisteredTime struct {
					Nanoseconds int64 `json:"nanoseconds"`
				} `json:"registered_time"`
				ReregisteredTime struct {
					Nanoseconds int64 `json:"nanoseconds"`
				} `json:"reregistered_time"`
			} `json:"frameworks"`
		} `json:"get_frameworks"`
		GetTasks struct {
			CompletedTasks []struct {
				AgentID struct {
					Value string `json:"value"`
				} `json:"agent_id"`
				ExecutorID struct {
					Value string `json:"value"`
				} `json:"executor_id"`
				FrameworkID struct {
					Value string `json:"value"`
				} `json:"framework_id"`
				Name      string `json:"name"`
				Resources []struct {
					Name   string `json:"name"`
					Role   string `json:"role"`
					Scalar struct {
						Value float64 `json:"value"`
					} `json:"scalar,omitempty"`
					Type   string `json:"type"`
					Ranges struct {
						Range []struct {
							Begin int `json:"begin"`
							End   int `json:"end"`
						} `json:"range"`
					} `json:"ranges,omitempty"`
				} `json:"resources"`
				State             string `json:"state"`
				StatusUpdateState string `json:"status_update_state"`
				StatusUpdateUUID  string `json:"status_update_uuid"`
				Statuses          []struct {
					AgentID struct {
						Value string `json:"value"`
					} `json:"agent_id"`
					ContainerStatus struct {
						NetworkInfos []struct {
							IPAddresses []struct {
								IPAddress string `json:"ip_address"`
							} `json:"ip_addresses"`
						} `json:"network_infos"`
					} `json:"container_status"`
					ExecutorID struct {
						Value string `json:"value"`
					} `json:"executor_id"`
					Source string `json:"source"`
					State  string `json:"state"`
					TaskID struct {
						Value string `json:"value"`
					} `json:"task_id"`
					Timestamp float64 `json:"timestamp"`
					UUID      string  `json:"uuid"`
				} `json:"statuses"`
				TaskID struct {
					Value string `json:"value"`
				} `json:"task_id"`
			} `json:"completed_tasks"`
		} `json:"get_tasks"`
	} `json:"get_state"`
}
type GET_AGENTS struct {
	Type      string `json:"type"`
	GetAgents struct {
		Agents []struct {
			Active    bool `json:"active"`
			AgentInfo struct {
				Hostname string `json:"hostname"`
				ID       struct {
					Value string `json:"value"`
				} `json:"id"`
				Port      int `json:"port"`
				Resources []struct {
					Name   string `json:"name"`
					Role   string `json:"role"`
					Scalar struct {
						Value float64 `json:"value"`
					} `json:"scalar,omitempty"`
					Type   string `json:"type"`
					Ranges struct {
						Range []struct {
							Begin int `json:"begin"`
							End   int `json:"end"`
						} `json:"range"`
					} `json:"ranges,omitempty"`
				} `json:"resources"`
			} `json:"agent_info"`
			Pid            string `json:"pid"`
			RegisteredTime struct {
				Nanoseconds int64 `json:"nanoseconds"`
			} `json:"registered_time"`
			TotalResources []struct {
				Name   string `json:"name"`
				Role   string `json:"role"`
				Scalar struct {
					Value float64 `json:"value"`
				} `json:"scalar,omitempty"`
				Type   string `json:"type"`
				Ranges struct {
					Range []struct {
						Begin int `json:"begin"`
						End   int `json:"end"`
					} `json:"range"`
				} `json:"ranges,omitempty"`
			} `json:"total_resources"`
			Version string `json:"version"`
		} `json:"agents"`
	} `json:"get_agents"`
}

type GET_FRAMEWORKS struct {
	Type          string `json:"type"`
	GetFrameworks struct {
		Frameworks []struct {
			Active        bool `json:"active"`
			Connected     bool `json:"connected"`
			FrameworkInfo struct {
				Checkpoint      bool    `json:"checkpoint"`
				FailoverTimeout float64 `json:"failover_timeout"`
				Hostname        string  `json:"hostname"`
				ID              struct {
					Value string `json:"value"`
				} `json:"id"`
				Name      string `json:"name"`
				Principal string `json:"principal"`
				Role      string `json:"role"`
				User      string `json:"user"`
			} `json:"framework_info"`
			RegisteredTime struct {
				Nanoseconds int64 `json:"nanoseconds"`
			} `json:"registered_time"`
			ReregisteredTime struct {
				Nanoseconds int64 `json:"nanoseconds"`
			} `json:"reregistered_time"`
		} `json:"frameworks"`
	} `json:"get_frameworks"`
}

type GET_EXECUTORS struct {
	Type         string `json:"type"`
	GetExecutors struct {
		Executors []struct {
			AgentID struct {
				Value string `json:"value"`
			} `json:"agent_id"`
			ExecutorInfo struct {
				Command struct {
					Shell bool   `json:"shell"`
					Value string `json:"value"`
				} `json:"command"`
				ExecutorID struct {
					Value string `json:"value"`
				} `json:"executor_id"`
				FrameworkID struct {
					Value string `json:"value"`
				} `json:"framework_id"`
			} `json:"executor_info"`
		} `json:"executors"`
	} `json:"get_executors"`
}

type GET_TASKS struct {
	Type     string `json:"type"`
	GetTasks struct {
		Tasks []struct {
			AgentID struct {
				Value string `json:"value"`
			} `json:"agent_id"`
			ExecutorID struct {
				Value string `json:"value"`
			} `json:"executor_id"`
			FrameworkID struct {
				Value string `json:"value"`
			} `json:"framework_id"`
			Name      string `json:"name"`
			Resources []struct {
				Name   string `json:"name"`
				Role   string `json:"role"`
				Scalar struct {
					Value float64 `json:"value"`
				} `json:"scalar,omitempty"`
				Type   string `json:"type"`
				Ranges struct {
					Range []struct {
						Begin int `json:"begin"`
						End   int `json:"end"`
					} `json:"range"`
				} `json:"ranges,omitempty"`
			} `json:"resources"`
			State             string `json:"state"`
			StatusUpdateState string `json:"status_update_state"`
			StatusUpdateUUID  string `json:"status_update_uuid"`
			Statuses          []struct {
				AgentID struct {
					Value string `json:"value"`
				} `json:"agent_id"`
				ContainerStatus struct {
					NetworkInfos []struct {
						IPAddresses []struct {
							IPAddress string `json:"ip_address"`
						} `json:"ip_addresses"`
					} `json:"network_infos"`
				} `json:"container_status"`
				ExecutorID struct {
					Value string `json:"value"`
				} `json:"executor_id"`
				Source string `json:"source"`
				State  string `json:"state"`
				TaskID struct {
					Value string `json:"value"`
				} `json:"task_id"`
				Timestamp float64 `json:"timestamp"`
				UUID      string  `json:"uuid"`
			} `json:"statuses"`
			TaskID struct {
				Value string `json:"value"`
			} `json:"task_id"`
		} `json:"tasks"`
	} `json:"get_tasks"`
}

type GET_ROLES struct {
	Type     string `json:"type"`
	GetRoles struct {
		Roles []struct {
			Name       string  `json:"name"`
			Weight     float64 `json:"weight"`
			Frameworks []struct {
				Value string `json:"value"`
			} `json:"frameworks,omitempty"`
			Resources []struct {
				Name   string `json:"name"`
				Role   string `json:"role"`
				Scalar struct {
					Value float64 `json:"value"`
				} `json:"scalar,omitempty"`
				Type   string `json:"type"`
				Ranges struct {
					Range []struct {
						Begin int `json:"begin"`
						End   int `json:"end"`
					} `json:"range"`
				} `json:"ranges,omitempty"`
			} `json:"resources,omitempty"`
		} `json:"roles"`
	} `json:"get_roles"`
}

type GET_WEIGHTS struct {
	Type       string `json:"type"`
	GetWeights struct {
		WeightInfos []struct {
			Role   string  `json:"role"`
			Weight float64 `json:"weight"`
		} `json:"weight_infos"`
	} `json:"get_weights"`
}

type GET_MASTER struct {
	Type      string `json:"type"`
	GetMaster struct {
		MasterInfo struct {
			Address struct {
				Hostname string `json:"hostname"`
				IP       string `json:"ip"`
				Port     int    `json:"port"`
			} `json:"address"`
			Hostname string `json:"hostname"`
			ID       string `json:"id"`
			IP       int    `json:"ip"`
			Pid      string `json:"pid"`
			Port     int    `json:"port"`
			Version  string `json:"version"`
		} `json:"master_info"`
	} `json:"get_master"`
}

type RESERVE_RESOURCES struct {
	Type             string `json:"type"`
	ReserveResources struct {
		AgentID struct {
			Value string `json:"value"`
		} `json:"agent_id"`
		Resources []struct {
			Type        string `json:"type"`
			Name        string `json:"name"`
			Reservation struct {
				Principal string `json:"principal"`
			} `json:"reservation"`
			Role   string `json:"role"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar"`
		} `json:"resources"`
	} `json:"reserve_resources"`
}

type UNRESERVE_RESOURCES struct {
	Type               string `json:"type"`
	UnreserveResources struct {
		AgentID struct {
			Value string `json:"value"`
		} `json:"agent_id"`
		Resources []struct {
			Type        string `json:"type"`
			Name        string `json:"name"`
			Reservation struct {
				Principal string `json:"principal"`
			} `json:"reservation"`
			Role   string `json:"role"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar"`
		} `json:"resources"`
	} `json:"unreserve_resources"`
}

type CREATE_VOLUMES struct {
	Type          string `json:"type"`
	CreateVolumes struct {
		AgentID struct {
			Value string `json:"value"`
		} `json:"agent_id"`
		Volumes []struct {
			Type string `json:"type"`
			Disk struct {
				Persistence struct {
					ID        string `json:"id"`
					Principal string `json:"principal"`
				} `json:"persistence"`
				Volume struct {
					ContainerPath string `json:"container_path"`
					Mode          string `json:"mode"`
				} `json:"volume"`
			} `json:"disk"`
			Name   string `json:"name"`
			Role   string `json:"role"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar"`
		} `json:"volumes"`
	} `json:"create_volumes"`
}

type DESTROY_VOLUMES struct {
	Type           string `json:"type"`
	DestroyVolumes struct {
		AgentID struct {
			Value string `json:"value"`
		} `json:"agent_id"`
		Volumes []struct {
			Disk struct {
				Persistence struct {
					ID        string `json:"id"`
					Principal string `json:"principal"`
				} `json:"persistence"`
				Volume struct {
					ContainerPath string `json:"container_path"`
					Mode          string `json:"mode"`
				} `json:"volume"`
			} `json:"disk"`
			Name   string `json:"name"`
			Role   string `json:"role"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar"`
			Type string `json:"type"`
		} `json:"volumes"`
	} `json:"destroy_volumes"`
}

type GROW_VOLUME struct {
	Type       string `json:"type"`
	GrowVolume struct {
		AgentID struct {
			Value string `json:"value"`
		} `json:"agent_id"`
		Volume struct {
			Disk struct {
				Persistence struct {
					ID        string `json:"id"`
					Principal string `json:"principal"`
				} `json:"persistence"`
				Volume struct {
					ContainerPath string `json:"container_path"`
					Mode          string `json:"mode"`
				} `json:"volume"`
			} `json:"disk"`
			Name   string `json:"name"`
			Role   string `json:"role"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar"`
			Type string `json:"type"`
		} `json:"volume"`
		Addition struct {
			Name   string `json:"name"`
			Role   string `json:"role"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar"`
			Type string `json:"type"`
		} `json:"addition"`
	} `json:"grow_volume"`
}

type SHRINK_VOLUME struct {
	Type         string `json:"type"`
	ShrinkVolume struct {
		AgentID struct {
			Value string `json:"value"`
		} `json:"agent_id"`
		Volume struct {
			Disk struct {
				Persistence struct {
					ID        string `json:"id"`
					Principal string `json:"principal"`
				} `json:"persistence"`
				Volume struct {
					ContainerPath string `json:"container_path"`
					Mode          string `json:"mode"`
				} `json:"volume"`
			} `json:"disk"`
			Name   string `json:"name"`
			Role   string `json:"role"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar"`
			Type string `json:"type"`
		} `json:"volume"`
		Subtract struct {
			Value float64 `json:"value"`
		} `json:"subtract"`
	} `json:"shrink_volume"`
}

type GET_MAINTENANCE_SCHEDULE struct {
	Type                   string `json:"type"`
	GetMaintenanceSchedule struct {
		Schedule struct {
			Windows []struct {
				MachineIds []struct {
					Hostname string `json:"hostname,omitempty"`
					IP       string `json:"ip,omitempty"`
				} `json:"machine_ids"`
				Unavailability struct {
					Start struct {
						Nanoseconds int64 `json:"nanoseconds"`
					} `json:"start"`
				} `json:"unavailability"`
			} `json:"windows"`
		} `json:"schedule"`
	} `json:"get_maintenance_schedule"`
}

type UPDATE_MAINTENANCE_SCHEDULE struct {
	Type                      string `json:"type"`
	UpdateMaintenanceSchedule struct {
		Schedule struct {
			Windows []struct {
				MachineIds []struct {
					Hostname string `json:"hostname,omitempty"`
					IP       string `json:"ip,omitempty"`
				} `json:"machine_ids"`
				Unavailability struct {
					Start struct {
						Nanoseconds int64 `json:"nanoseconds"`
					} `json:"start"`
				} `json:"unavailability"`
			} `json:"windows"`
		} `json:"schedule"`
	} `json:"update_maintenance_schedule"`
}

type START_MAINTENANCE struct {
	Type             string `json:"type"`
	StartMaintenance struct {
		Machines []struct {
			Hostname string `json:"hostname"`
			IP       string `json:"ip"`
		} `json:"machines"`
	} `json:"start_maintenance"`
}

type STOP_MAINTENANCE struct {
	Type            string `json:"type"`
	StopMaintenance struct {
		Machines []struct {
			Hostname string `json:"hostname"`
			IP       string `json:"ip"`
		} `json:"machines"`
	} `json:"stop_maintenance"`
}

type GET_QUOTA struct {
	Type     string `json:"type"`
	GetQuota struct {
		Status struct {
			Infos []struct {
				Guarantee []struct {
					Name   string `json:"name"`
					Role   string `json:"role"`
					Scalar struct {
						Value float64 `json:"value"`
					} `json:"scalar"`
					Type string `json:"type"`
				} `json:"guarantee"`
				Principal string `json:"principal"`
				Role      string `json:"role"`
			} `json:"infos"`
		} `json:"status"`
	} `json:"get_quota"`
}

type SET_QUOTA struct {
	Type     string `json:"type"`
	SetQuota struct {
		QuotaRequest struct {
			Force     bool `json:"force"`
			Guarantee []struct {
				Name   string `json:"name"`
				Role   string `json:"role"`
				Scalar struct {
					Value float64 `json:"value"`
				} `json:"scalar"`
				Type string `json:"type"`
			} `json:"guarantee"`
			Role string `json:"role"`
		} `json:"quota_request"`
	} `json:"set_quota"`
}

type REMOVE_QUOTA struct {
	Type        string `json:"type"`
	RemoveQuota struct {
		Role string `json:"role"`
	} `json:"remove_quota"`
}

type MARK_AGENT_GONE struct {
	Type          string `json:"type"`
	MarkAgentGone struct {
		AgentID struct {
			Value string `json:"value"`
		} `json:"agent_id"`
	} `json:"mark_agent_gone"`
}
