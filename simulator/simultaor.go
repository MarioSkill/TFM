package main

import (
	"fmt"
	"gitlab.com/marioskill/simulator/eidas"
	//"log"
	"math"
	"math/rand"
	"os"
	//"os/exec"
	"strconv"
	"sync"
	"time"
)

var timer = 0.0
var endInterval = 0.0
var Lai = -1

//./eidas 99 http://163.117.148.105 6060 http://163.117.148.105:8081/eidasResult

func main() {

	var clientID = os.Args[1]
	var eidasServer = os.Args[2] //"http://163.117.148.105:8080"
	var eidasPort = os.Args[3]
	var resultServer = os.Args[4] //"http://163.117.148.105:8081/eidasResult"
	var tasaLambda, _ = strconv.ParseFloat(os.Args[5], 64)
	var maxClients, _ = strconv.Atoi(os.Args[6])
	var timestampClientRegisted = os.Args[7]
	var serverTime = os.Args[8]
	var timetaskDeploy = os.Args[9]

	var wg sync.WaitGroup

	λ := float64(1) / tasaLambda
	var elapsed float64
	fmt.Println("Tasa de llegada: ", λ, " por segundo")

	var current = 0
	for Lai < maxClients-1 {

		start := time.Now()
		//fmt.Println(float64(elapsed))
		check(float64(elapsed), λ)
		if Lai == current {
			wg.Add(1)
			go func() {
				defer wg.Done()
				eidas.StartTest(strconv.Itoa(current)+"_"+clientID, eidasServer, eidasPort, resultServer, timestampClientRegisted, serverTime, timetaskDeploy)
			}()
			current++
		}
		//time.Sleep(1 * time.Second)
		elapsed = time.Since(start).Seconds()

	}
	wg.Wait()
	//fmt.Println(Lai)
}

func nextArrive(λ float64) float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return -math.Log(float64(1.0)-rand.Float64()) / λ
}

func check(dt_sec float64, λ float64) {
	timer += dt_sec
	if timer > endInterval {
		timer = 0.0
		endInterval = nextArrive(λ)
		Lai++
	}
}
