package main

import (
	"fmt"
	"bytes"
	"golang.org/x/crypto/ssh"
	"net"
	"io/ioutil"
	"crypto/x509"
	"encoding/pem"
	"flag"
)

func decrypt(key []byte, password []byte) []byte {
	block, rest := pem.Decode(key)
	if len(rest) > 0 {
		panic("Extra data included in key")
	}

	der, err := x509.DecryptPEMBlock(block, password)
	if err != nil {
		panic(err.Error())
	}
	return der
}

func PublicKeyFile(file string, password string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
		return nil
	}

	der := decrypt(buffer, []byte(password))
	key, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		panic(err.Error())
	}

	singer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return ssh.PublicKeys(singer)
}

func PublicKeyFile2(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func main() {
	var keyPath = flag.String("key", "", "private key path")
	var keyPwd = flag.String("pwd", "", "private key password")
	var hostIP = flag.String("host", "", "remote ip")
	var hostPort = flag.String("port", "", "remote port")
	flag.Parse()

	if *keyPath == "" {
		fmt.Println("key is empty")
		return
	}

	if *hostIP == "" {
		fmt.Println("host is empty")
		return
	}

	if *hostPort == "" {
		fmt.Println("port is empty")
		return
	}

	var config *ssh.ClientConfig
	if len(*keyPwd) != 0 {
		config = &ssh.ClientConfig{
			User: "lyg",
			Auth: []ssh.AuthMethod{
				PublicKeyFile(*keyPath, *keyPwd),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	} else {
		config = &ssh.ClientConfig{
			User: "lyg",
			Auth: []ssh.AuthMethod{
				PublicKeyFile2(*keyPath),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", *hostIP, *hostPort), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
