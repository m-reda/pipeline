package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os/exec"
	"strconv"
	"time"
)

func funcUnits() []cli.Command {

	return []cli.Command{
		{Name: "sleep", Usage: "Sleep for seconds", Action: funcSleep},
		{Name: "url", Usage: "Sleep for seconds", Action: funcURL},
		{Name: "command", Usage: "Run a command", Action: funcCommand},
		{Name: "zip", Usage: "Compress files", Action: funcZip},
		{Name: "unzip", Usage: "Unzip file", Action: funcUnzip},
		{
			Name:   "remove",
			Usage:  "Remove file or directory",
			Action: funcEmail,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "smtp-host"},
				cli.StringFlag{Name: "smtp-username"},
				cli.StringFlag{Name: "smtp-password"},
				cli.StringFlag{Name: "smtp-port", Value: "25"},
			},
		},
		{
			Name:   "ssh",
			Usage:  "Run SSH command",
			Action: funcSSH,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "password", EnvVar: "PIPELINE_SSH_PASSWORD"},
				cli.StringFlag{Name: "key", EnvVar: "PIPELINE_SSH_KEY"},
			},
		},
	}
}

func funcSleep(c *cli.Context) error {
	// inputs
	seconds := c.Args().Get(0)
	expected(seconds)

	// execute
	i, err := strconv.Atoi(seconds)
	checkErr(err)
	time.Sleep(time.Duration(i) * time.Second)

	return nil
}

func funcURL(c *cli.Context) error {
	// inputs
	method, url := c.Args().Get(0), c.Args().Get(1)
	expected(method, url)

	// execute
	r, err := http.NewRequest(method, url, nil)
	checkErr(err)

	client := &http.Client{}
	resp, err := client.Do(r)
	checkErr(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	// outputs
	fmt.Println("url:" + url)
	fmt.Println("method:" + method)
	fmt.Println("body:" + string(body))
	return nil
}

func funcEmail(c *cli.Context) error {
	// inputs
	to, msg := c.Args().Get(0), c.Args().Get(1)
	expected(to, msg)

	// execute
	err := smtp.SendMail(
		c.String("smtp-host")+":"+c.String("smtp-port"),
		smtp.PlainAuth("", c.String("smtp-username"), c.String("smtp-password"), c.String("smtp-host")),
		c.String("smtp-username"),
		[]string{to},
		[]byte(msg),
	)
	checkErr(err)

	return nil
}

func funcCommand(c *cli.Context) error {
	// inputs
	command := c.Args().Get(0)
	expected(command)

	// execute
	console, err := exec.Command("/bin/sh", "-c", command).Output()
	checkErr(err)

	// outputs
	fmt.Println("output:" + base64.StdEncoding.EncodeToString(console))
	return nil
}

func funcZip(c *cli.Context) error {
	// inputs
	zipName, files := c.Args().Get(0), c.Args().Get(1)
	expected(zipName, files)

	// execute
	err := exec.Command("zip", "-r", "-X", zipName, files).Run()
	checkErr(err)

	// outputs
	fmt.Println("zip_name:" + zipName)
	fmt.Println("files:" + files)
	return nil
}

func funcUnzip(c *cli.Context) error {
	// inputs
	zipName := c.Args().Get(0)
	expected(zipName)

	// execute
	err := exec.Command("unzip", zipName).Run()
	checkErr(err)

	// outputs
	fmt.Println("zip_name:" + zipName)
	return nil
}

func funcSSH(c *cli.Context) error {
	// inputs
	host, user, command := c.Args().Get(0), c.Args().Get(1), c.Args().Get(2)
	expected(user, host, command)

	// auth method
	var authMethod ssh.AuthMethod
	if keyPath := c.String("key"); keyPath != "" {
		buf, err := ioutil.ReadFile(keyPath)
		checkErr(err)

		key, err := ssh.ParsePrivateKey(buf)
		checkErr(err)

		authMethod = ssh.PublicKeys(key)
	} else if password := c.String("password"); password != "" {
		authMethod = ssh.Password(password)
	} else {
		log.Fatal("Password or key is requird")
	}

	// config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{authMethod},
	}

	// connect
	conn, err := ssh.Dial("tcp", host, config)
	checkErr(err)
	defer conn.Close()

	session, err := conn.NewSession()
	checkErr(err)
	defer session.Close()

	// execute
	var output bytes.Buffer
	session.Stdout = &output
	err = session.Run(command)
	checkErr(err)

	// outputs
	fmt.Println("host:" + host)
	fmt.Println("username:" + user)
	fmt.Println("output:" + base64.StdEncoding.EncodeToString(output.Bytes()))
	return nil
}
