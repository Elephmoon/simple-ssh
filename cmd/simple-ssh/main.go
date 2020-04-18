package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	defaultSSHKey  = ".ssh/id_rsa"
	defaultHomeDir = "HOME"
	inputSpeed     = 14400
	outputSpeed    = 14400
	exitCommand    = "exit\n"
)

var (
	user       = flag.String("u", "", "User name")
	privateKey = flag.String("pk", getDefaultPrivateKey(), "Private key file")
	host       = flag.String("h", "", "Host")
	port       = flag.Int("p", 22, "Port")
)

func main() {
	flag.Parse()

	key, err := ioutil.ReadFile(*privateKey)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}

	clientConfig := &ssh.ClientConfig{
		User:            *user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	pipeIn, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: inputSpeed,
		ssh.TTY_OP_OSPEED: outputSpeed,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	sendCommandToServer(session, pipeIn)
	log.Print("session end")
}

func sendCommandToServer(session *ssh.Session, in io.WriteCloser) {
	for {
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Fprint(in, str)
		if err != nil {
			log.Fatal(err)
		}
		if str == exitCommand {
			session.Close()
			return
		}
	}
}

func getDefaultPrivateKey() string {
	home := os.Getenv(defaultHomeDir)
	if len(home) > 0 {
		return path.Join(home, defaultSSHKey)
	}
	return ""
}
