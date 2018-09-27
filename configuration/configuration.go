package configuration

var Root = "/TFM/mesos-framework"
var AppClusterPath = "/home/mvasile/bin/"
var Mediciones = "mediciones"
var ClusetConf = Root + "/assets/cluster-conf"
var TestConf = ClusetConf + "/tests"
var LoafBalancerServers = ClusetConf + "/lb_servers"

var Apps = ClusetConf + "/bin"
var FrameworksConf = ClusetConf + "/frameworks"
var NodesConfs = ClusetConf + "/nodeConf"
var Results = ClusetConf + "/results"

var GUI_port = "8081"
var Front_port = "6060"
var Monitor_port = "9999"

var ClusterUser = "mvasile"
var ClusterURL = "tucan.arcos.inf.uc3m.es"

var Mesos_Server_url_API = "http://localhost:5050/api/v1"

var EidasResults = "http://192.168.1.10:8081/eidasResult"

var Scheduler_baseline_ws_port = ":8082"
var Scheduler_edf_ws_port = ":8083"
var Scheduler_prioridades_ws_port = ":8084"

var FrameworksNames = []string{"baseline", "d"}

var Heartbeat_window = 30 //minutos
var Heartbeat_protocol = "udp"
var Heartbeat_port = ":8694"

const QAL_VERY_HIGH int = 5
const QAL_HIGH int = 4
const QAL_MEDIUM int = 3
const QAL_MEDIUM_LOW int = 2
const QAL_LOW int = 1
