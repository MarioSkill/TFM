package main

import (
	"encoding/json"
	"fmt"
	"github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

var tiempos = &sync.Mutex{}
var resultados = &sync.Mutex{}

type serverResources struct {
	Info []float64 `json:"Info"`
	Host string    `json:"Host"`
}

type Tiempos struct {
	ID                string `csv:"id"`
	InicioEjecucion   string `csv:"tiempoInicioEjecucion"`    // cuando empieza a ejecutar
	ClienteRegistrado string `csv:"tiempoRegistroDelCliente"` // Cuando se registra en lb
	TareaPlanificada  string `csv:"tiempoTareaPlanificada"`   //Cuando entra la tarea en el planificador

}
type Login struct {
	ID                string `csv:"id"`
	InicioEjecucion   string `csv:"tiempoInicioEjecucion"`    // cuando empieza a ejecutar
	ClienteRegistrado string `csv:"tiempoRegistroDelCliente"` // Cuando se registra en lb
	TareaPlanificada  string `csv:"tiempoTareaPlanificada"`   //Cuando entra la tarea en el planificador
	Login             string `csv:"Login"`                    // cuando empieza a ejecutar
	Error             string `csv:"Error"`                    // Cuando se registra en lb
	LoginTime         string `csv:"LoginTime"`                //Cuando entra la tarea en el planificador
}

var resourcesPath string

func EidasGetResults(c *gin.Context) {

	/*r := "error=" + c.PostForm("error") + ", Login=" + c.PostForm("login") + ", ClientID= " +
		c.PostForm("ClientID") + ", time=" + c.PostForm("time") + "\r\n"
	/*
		f, err := os.OpenFile(resourcesPath+"/resultados.txt", os.O_APPEND|os.O_WRONLY, 0600)

		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(r); err != nil {
			panic(err)
		}*/

	//scp eidas mvasile@tucan.arocs.inf.uc3m.es:/home/mvasile/bin/eidas

	var aux = []*Login{}
	aux = append(aux, &Login{
		ID:                c.PostForm("ID"),
		Login:             c.PostForm("login"),
		Error:             c.PostForm("error"),
		LoginTime:         c.PostForm("LoginTime"),
		InicioEjecucion:   c.PostForm("InicioEjecucion"),
		ClienteRegistrado: c.PostForm("ClienteRegistrado"),
		TareaPlanificada:  c.PostForm("TareaPlanificada")})

	datacsv, err := gocsv.MarshalString(&aux)
	//datacsv, err := gocsv.MarshalString(&aux)
	csvContent := strings.Split(datacsv, "\n")
	csvContent = csvContent[:0+copy(csvContent[0:], csvContent[0+1:])]
	s := strings.Join(csvContent, "\n")
	fmt.Println(csvContent)
	resultados.Lock()
	f, _ := os.OpenFile(resourcesPath+"/resultados.csv", os.O_APPEND|os.O_WRONLY, 0755)
	if _, err = f.WriteString(s); err != nil {
		resultados.Unlock()
		panic(err)
	}
	f.Close()
	resultados.Unlock()
	c.JSON(http.StatusOK, "ok")

}
func TiempoDeEspera(c *gin.Context) {

	inicioEjecucion := c.PostForm("InicioEjecucion")
	clienteRegistrado := c.PostForm("ClienteRegistrado")
	tareaPlanificada := c.PostForm("TareaPlanificada")
	id := c.PostForm("ID")

	var aux = []*Tiempos{}
	aux = append(aux, &Tiempos{ID: id, InicioEjecucion: inicioEjecucion, ClienteRegistrado: clienteRegistrado, TareaPlanificada: tareaPlanificada})

	datacsv, err := gocsv.MarshalString(&aux)
	csvContent := strings.Split(datacsv, "\n")
	csvContent = csvContent[:0+copy(csvContent[0:], csvContent[0+1:])]
	s := strings.Join(csvContent, "\n")
	//fmt.Println(s)
	resultados.Lock()
	f, _ := os.OpenFile(resourcesPath+"/tiempoServicio.csv", os.O_APPEND|os.O_WRONLY, 0755)
	if _, err = f.WriteString(s); err != nil {
		resultados.Unlock()
		panic(err)
	}
	f.Close()
	resultados.Unlock()
}

var router *gin.Engine
var serverport string
var port string

func main() {
	gocsv.TagSeparator = ";"
	fmt.Println("puerto ruta a los rescursos")
	port = os.Args[1]
	resourcesPath = os.Args[2]
	serverport = "9999"
	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:         "views",
		Extension:    ".tpl",
		Master:       "layouts/base",
		DisableCache: true,
	})
	router.POST("/eidasResult", EidasGetResults)
	router.POST("/TimeClient", TiempoDeEspera)
	router.GET("/start", Start)
	router.GET("/stop", Stop)
	router.GET("/status", Status)
	router.GET("/status/monitor", StatusMonitor)
	router.GET("/status/monitor/host", GetDataByHost)

	// Start serving the application
	router.Run(":" + port)

}

func Start(c *gin.Context) {
	server := c.Query("node")
	//folder := resourcesPath + "/" + c.PostForm("folder") //carpeta donde guardar los resultados

	//cmdStr := resourcesPath + "/monitor 9999 " + folder
	//	cmdStr := "/TFM/mesos-framework/go/src/gitlab.com/marioskill/monitor/monitor 9999 " + folder

	//ssh := `ssh ` + server + ` '` + cmdStr + `'`
	//ssh := cmdStr
	//exec.Command("/bin/sh", "-c", ssh).Start()
	//fmt.Println(ssh)

	//time.Sleep(1 * time.Second)
	MakeRequest("http://"+server, "/start")
	c.JSON(http.StatusOK, "ok")
}

func Stop(c *gin.Context) {
	server := c.Query("node")

	MakeRequest("http://"+server, "/stop")
	MakeRequest("http://"+server, "/save")
	///MakeRequest("http://"+server, "/exit")
}

func MakeRequest(server string, acction string) ([]byte, error) {
	resp, err := http.Get(server + ":" + serverport + acction)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(body))
	return body, nil
}

func Status(c *gin.Context) {
	c.JSON(http.StatusOK, "{Port:"+port+",Path:"+resourcesPath+"}")
}

func StatusMonitor(c *gin.Context) {
	server := c.Query("host")
	servers := strings.Split(server, ",")
	var result = ""
	for i := 0; i < len(servers); i++ {
		response, _ := MakeRequest("http://"+servers[i], "/status")
		result = result + servers[i] + ":" + string(response) + ":Q:" + servers[i] + ":" + serverport + "/status" + "\n"
	}
	c.JSON(http.StatusOK, result)
}

func GetDataByHost(c *gin.Context) {

	result, _ := MakeRequest("http://"+c.Query("host"), "/serverResources?host="+c.Query("host"))
	var r serverResources
	json.Unmarshal(result, &r)

	c.JSON(http.StatusOK, r)
}
