package logic

import (
	"bufio"
	"fmt"
	"github.com/crazy-me/ops_dial/cli"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	pkg          = "telegraf.zip"     // telegraf包名
	originalPath = "D:\\telegraf.zip" // telegraf 本地路径
	remotePath   = "/opt"             // telegraf 上传的远端路径，修改此路径时还需同时修改telegraf相关配置
)

func Run(fileName string) {
	var (
		err    error
		port   int
		client *cli.SSHClient
		scp    *cli.TransferInfo
		exec   *cli.ExecInfo
	)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Open %s Failed:%s\n", fileName, err.Error())
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		accountSlice := strings.Split(scanner.Text(), "|")
		port, _ = strconv.Atoi(accountSlice[3])
		// 建立SSH连接
		client, err = cli.NewSSHClient(accountSlice[2], port, cli.AuthConfig{
			User:     accountSlice[0],
			Password: accountSlice[1],
			Timeout:  10 * time.Second,
		})

		if err != nil { // SSH连接建立失败跳过继续执行下一个
			log.Printf("SSHClient Connect Failed host:%s, err: %s\n", accountSlice[2], err.Error())
			continue
		}

		// 开始安装 telegraf
		fmt.Printf("%s Connect Success Upload telegraf package...\n", accountSlice[2])
		scp, err = client.Upload(originalPath, remotePath+"/"+pkg)
		if err != nil {
			log.Printf("%s Upload %s Failed:%s\n", accountSlice[2], pkg, err.Error())
			continue
		}

		fmt.Printf("Upload %s Success Size:%d\n", pkg, scp.TransferByte)

		exec, err = client.Exec(fmt.Sprintf("unzip %s/telegraf.zip -d %s", remotePath, remotePath))
		if err != nil {
			log.Printf("%s unzip %s Failed:%s\n", accountSlice[2], pkg, err.Error())
			continue
		}
		fmt.Printf("unzip %s Success:%s\n", pkg, exec.OutputString())

		_, _ = client.Exec(fmt.Sprintf("chmod -R +x %s/telegraf", remotePath))
		exec, err = client.Exec(fmt.Sprintf("%s/telegraf/entry.sh", remotePath))
		if err != nil {
			log.Printf("%s shell %s/telegraf/entry.sh Failed:%s\n", accountSlice[2], remotePath, err.Error())
			continue
		}

		fmt.Printf("%s install telegraf %s\n", accountSlice[2], exec.OutputString())
	}

}
