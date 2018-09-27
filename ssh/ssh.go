package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	//	"os"
)

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(1)
		log.Fatalln(fmt.Sprintf("Cannot read SSH public key file %s", file))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Println(2)
		log.Fatalln(fmt.Sprintf("Cannot parse SSH public key file %s", file))
		return nil
	}
	return ssh.PublicKeys(key)
}

// Get default location of a private key
func privateKeyPath() string {
	return "/home/mario" + "/.ssh/id_rsa"
}

// Get private key for ssh authentication
func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, _ := ioutil.ReadFile(keyPath)
	return ssh.ParsePrivateKey(buff)
}

// Get ssh client config for our connection
// SSH config will use 2 authentication strategies: by key and by password
func makeSshConfig() (*ssh.ClientConfig, error) {
	/*
		key, err := parsePrivateKey(privateKeyPath())
		if err != nil {
			return nil, err
		}
	*/
	sshConfig := ssh.ClientConfig{
		// SSH connection username
		User: "mvasile",
		Auth: []ssh.AuthMethod{publicKeyFile("/home/mario/.ssh/id_rsa")}, // put here your private key path

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &sshConfig, nil
}

// Handle local client connections and tunnel data to the remote serverq
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			fmt.Println(3)
			log.Println("error while copy remote->local:", err)
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			fmt.Println(4)
			log.Println(err)
		}
		chDone <- true
	}()

	<-chDone
}

func Start(l, r string) error {
	//func main() {
	// Connection settings
	sshAddr := "tucan.arcos.inf.uc3m.es:22"
	localAddr := l  //os.Args[1]
	remoteAddr := r //os.Args[2]
	//localAddr := "localhost:6060"
	//remoteAddr := "localhost:6060"

	// Build SSH client configuration
	cfg, err := makeSshConfig()
	if err != nil {
		//log.Fatalln(err)
		//return err
		fmt.Println(5)
		fmt.Println(err)
	}

	// Establish connection with SSH server
	conn, err := ssh.Dial("tcp", sshAddr, cfg)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(9)
		fmt.Println(err)
	}
	defer conn.Close()

	// Establish connection with remote server
	remote, err := conn.Dial("tcp", remoteAddr)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(6)
		fmt.Println(err)
	}

	// Start local server to forward traffic to remote connection
	local, err := net.Listen("tcp", localAddr)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(7)
		fmt.Println(err)
	}
	defer local.Close()

	// Handle incoming connections
	for {
		client, err := local.Accept()
		if err != nil {
			//log.Fatalln(err)
			fmt.Println(8)
			fmt.Println(err)
		}

		handleClient(client, remote)
	}
}

/*
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH public key file %s", file))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH public key file %s", file))
		return nil
	}
	return ssh.PublicKeys(key)
}

// local service to be forwarded
var localEndpoint = Endpoint{
	Host: "localhost",
	Port: 5050,
}

// remote SSH server
var serverEndpoint = Endpoint{
	Host: "tucan.arcos.inf.uc3m.es",
	Port: 22,
}

// remote forwarding port (on remote SSH server network)
var remoteEndpoint = Endpoint{
	Host: "localhost",
	Port: 6060,
}

func main() {

	// refer to https://godoc.org/golang.org/x/crypto/ssh for other authentication types
	sshConfig := &ssh.ClientConfig{
		// SSH connection username
		User: "mvasile",
		Auth: []ssh.AuthMethod{publicKeyFile("/home/mvasile/.ssh/id_rsa")}, // put here your private key path

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
	}
	defer listener.Close()

	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localEndpoint.String())
		if err != nil {
			log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
		}

		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		handleClient(client, local)
	}

}
*/
