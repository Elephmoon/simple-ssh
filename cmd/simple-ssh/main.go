package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
)

const (
	defaultSSHKey  = ".ssh/id_rsa"
	defaultHomeDir = "HOME"
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
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}

	data, err := session.CombinedOutput("uname -a")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(data))
}

func getDefaultPrivateKey() string {
	home := os.Getenv(defaultHomeDir)
	if len(home) > 0 {
		return path.Join(home, defaultSSHKey)
	}
	return ""
}
