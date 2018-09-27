package main

import (
	"fmt"
	"github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	//"os/signal"
	"strconv"
	"strings"
	//"syscall"
	"time"
)

type Resource struct {
	NODO string `csv:"nodo"`
	CPU  string `csv:"cpu"`
	MEM  string `csv:"mem"`
	X    int    `csv:"x"`
}

var Resources = []*Resource{}
var Medir bool

type serverResources struct {
	Info []float64 `json:"Info"`
	Host string    `json:"Host"`
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func medirRecursos() {
	Resources = []*Resource{}
	x := 1
	for Medir {
		idle0, total0 := getCPUSample()
		time.Sleep(1 * time.Second)
		idle1, total1 := getCPUSample()
		cmd := exec.Command("sh", "-c", `free -m | awk 'NR==2{printf "%s", $3}'`)
		out, _ := cmd.Output()
		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
		//now := time.Now()
		Resources = append(Resources, &Resource{CPU: getFormattedValue(cpuUsage), MEM: string(out), NODO: nodo, X: x})
		x++
		//	fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
	}
}

var router *gin.Engine
var nodo string
var folder string

func main() {
	fmt.Println("uso: monitor puerto")
	var port = os.Args[1]
	folder = os.Args[2]
	name, err := os.Hostname()
	if err != nil {
		nodo = "undefined"
	} else {
		nodo = name
	}

	gocsv.TagSeparator = ";"
	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:         "views",
		Extension:    ".tpl",
		Master:       "layouts/base",
		DisableCache: true,
	})

	router.GET("/start", Start)
	router.GET("/stop", Stop)
	router.GET("/save", SaveResultados)
	router.GET("/status", Status)
	router.GET("/exit", Exit)
	router.GET("/serverResources", GetServerResources)

	// Start serving the application
	router.Run(":" + port)
}
func SaveResultados(c *gin.Context) {

	f, _ := os.OpenFile(folder+"/"+nodo+"_datos.csv", os.O_RDWR|os.O_CREATE, 0755)

	err := gocsv.MarshalFile(&Resources, f) // Get all clients as CSV string
	//err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
	f.Close()
	c.JSON(http.StatusOK, "ok")
}
func Status(c *gin.Context) {

	if Medir == true {
		c.JSON(http.StatusOK, "runing")
	} else {
		c.JSON(http.StatusOK, "stoped")
	}
}
func Start(c *gin.Context) {

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0777)
	}

	if Medir != true {

		Medir = true

		go func() {
			medirRecursos()
		}()
	}
	c.JSON(http.StatusOK, "ok")
}
func Stop(c *gin.Context) {
	Medir = false
	c.JSON(http.StatusOK, "ok")
}
func Exit(c *gin.Context) {
	Medir = false
	//SaveResultados(c)
	c.JSON(http.StatusOK, "ok")
	os.Exit(3)
}

func GetServerResources(c *gin.Context) {
	idle0, total0 := getCPUSample()
	time.Sleep(500 * time.Millisecond)
	idle1, total1 := getCPUSample()
	cmd := exec.Command("sh", "-c", `free -m | awk 'NR==2{printf "%s", $3}'`)
	out, _ := cmd.Output()
	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	value, _ := strconv.ParseFloat(string(out), 64)

	//fmt.Println(out)
	res := serverResources{Info: []float64{value / 1024, cpuUsage}, Host: c.Query("host")}

	c.JSON(http.StatusOK, res)
}

func getFormattedValue(percentageValue float64) string {
	value := fmt.Sprintf("%.2f", percentageValue)
	return strings.Replace(value, ".", ",", -1)
}
