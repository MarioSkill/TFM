package prioridades

import (
	"gitlab.com/marioskill/loadbalancer/prioridades/levels"
	"gitlab.com/marioskill/loadbalancer/prioridades/levels/high"
	"gitlab.com/marioskill/loadbalancer/prioridades/levels/low"
	"gitlab.com/marioskill/loadbalancer/prioridades/levels/medium"
)

func SubClient(id string, host string, ip string, version string, p int, qal int) levels.ServerInformation {
	var info levels.ServerInformation
	switch qal {
	case 0: //low
		info = low.SubClient(id, host, ip, version, p, qal)
	case 1: //medium
		info = medium.SubClient(id, host, ip, version, p, qal)
	case 2: //high
		info = high.SubClient(id, host, ip, version, p, qal)
	}
	return info
}

func SubServer(id string, host string, ip string, version string, port string, p int, qal int) {

	switch qal {
	case 0: //low
		low.SubServer(id, host, ip, version, port, p, qal)
	case 1: //medium
		medium.SubServer(id, host, ip, version, port, p, qal)
	case 2: //high
		high.SubServer(id, host, ip, version, port, p, qal)
	}
}
