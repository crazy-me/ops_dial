package cli

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

// Cli 交互式Terminal结构体
type Cli struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Address  string      `json:"address"`
	Port     int         `json:"port"`
	Client   *ssh.Client `json:"client"`
	Last     string      `json:"last"`
}

func (c *Cli) connect() (err error) {
	conf := ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.Address, c.Port)
	terminalClient, err := ssh.Dial("tcp", addr, &conf)
	if err != nil {
		return err
	}

	c.Client = terminalClient
	return
}

func (c *Cli) Run(shell string) (string, error) {
	if c.Client == nil {
		if err := c.connect(); err != nil {
			return "", err
		}
	}

	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	output, err := session.CombinedOutput(shell)

	c.Last = string(output)
	return c.Last, err
}

// New 初始化 Struct
func New(username, password, address string, port int) *Cli {
	return &Cli{
		Username: username,
		Password: password,
		Address:  address,
		Port:     port,
	}
}
