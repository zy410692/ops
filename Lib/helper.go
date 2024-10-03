package Lib

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

func clearConsole() {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")

		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		log.Fatal("目前clear操作只写了Linux、darwin、Windows3个。其他没做，累的")
	}
}
func MustFlag(name, t string, cmd *cobra.Command) interface{} {
	switch t {
	case "string":
		if v, err := cmd.Flags().GetString(name); err == nil && v != "" {
			return v
		}
	case "int":
		if v, err := cmd.Flags().GetInt(name); err == nil && v != 0 {
			return v
		}
	}

	log.Fatal(name, " is required")
	return nil
}

var ShellModes = ssh.TerminalModes{
	ssh.ECHO:          1,     // enable echoing
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

func getDefaultPrivateKeyPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir + "/.ssh/id_rsa", nil
}

func SSHConnectKey(user, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	privateKeyPath, err := getDefaultPrivateKeyPath()
	if err != nil {
		log.Fatalf("Failed to get default private key path: %v", err)
	}
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}

	// 解析私钥
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.PublicKeys(signer))
	hostKeyCallbk := ssh.InsecureIgnoreHostKey()

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}
