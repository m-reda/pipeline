package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/urfave/cli"
	"os"
	"strings"
)

var ftpClient *ftp.ServerConn

func ftpConnect(host, port, username, password string) {
	var err error
	ftpClient, err = ftp.Dial(host + ":" + port)
	checkErr(err)

	err = ftpClient.Login(username, password)
	checkErr(err)
}

func ftpUnits() []cli.Command {
	flags := []cli.Flag{
		cli.StringFlag{Name: "host", EnvVar: "PIPELINE_FTP_HOST"},
		cli.StringFlag{Name: "port", EnvVar: "PIPELINE_FTP_PORT"},
		cli.StringFlag{Name: "username", EnvVar: "PIPELINE_FTP_USER"},
		cli.StringFlag{Name: "password", EnvVar: "PIPELINE_FTP_PASS"},
	}

	return []cli.Command{
		{Name: "list", Usage: "List directory contents", Action: ftpList, Flags: flags},
		{Name: "rm", Usage: "Delete file", Action: ftpRemoveDir, Flags: flags},
		{Name: "rmdir", Usage: "Delete directory", Action: ftpRemove, Flags: flags},
		{Name: "mkdir", Usage: "Make new directory", Action: ftpMakeDir, Flags: flags},
		{Name: "rename", Usage: "Rename file or directory", Action: ftpRename, Flags: flags},
		{Name: "upload", Usage: "Upload file", Action: ftpUpload, Flags: flags},
	}
}

func ftpList(c *cli.Context) error {
	// connect
	ftpConnect(c.String("host"), c.String("port"), c.String("username"), c.String("password"))
	defer ftpClient.Quit()

	// inputs
	dir := c.Args().Get(0)
	expected(dir)

	// execute
	entries, err := ftpClient.NameList(dir)
	checkErr(err)

	// outputs
	fmt.Println("directory:" + dir)
	fmt.Println("list:" + strings.Join(entries[2:], "/"))
	return nil
}

func ftpRemove(c *cli.Context) error {
	// connect
	ftpConnect(c.String("host"), c.String("port"), c.String("username"), c.String("password"))
	defer ftpClient.Quit()

	// inputs
	path := c.Args().Get(0)
	expected(path)

	// execute
	err := ftpClient.Delete(path)
	checkErr(err)

	// outputs
	fmt.Println("path:" + path)
	return nil
}

func ftpMakeDir(c *cli.Context) error {
	// connect
	ftpConnect(c.String("host"), c.String("port"), c.String("username"), c.String("password"))
	defer ftpClient.Quit()

	// inputs
	path := c.Args().Get(0)
	expected(path)

	// execute
	err := ftpClient.MakeDir(path)
	checkErr(err)

	// outputs
	fmt.Println("path:" + path)
	return nil
}

func ftpRename(c *cli.Context) error {
	// connect
	ftpConnect(c.String("host"), c.String("port"), c.String("username"), c.String("password"))
	defer ftpClient.Quit()

	// inputs
	from, to := c.Args().Get(0), c.Args().Get(1)
	expected(from, to)

	// execute
	err := ftpClient.Rename(from, to)
	checkErr(err)

	// outputs
	fmt.Println("from:" + from)
	fmt.Println("to:" + to)
	return nil
}

func ftpUpload(c *cli.Context) error {
	// connect
	ftpConnect(c.String("host"), c.String("port"), c.String("username"), c.String("password"))
	defer ftpClient.Quit()

	// inputs
	source, dest := c.Args().Get(0), c.Args().Get(1)
	expected(source, dest)

	// execute
	file, err := os.Open(source)
	checkErr(err)

	err = ftpClient.Stor(dest, file)
	checkErr(err)

	// outputs
	fmt.Println("source:" + source)
	fmt.Println("destination:" + dest)
	return nil
}

func ftpRemoveDir(c *cli.Context) error {
	// connect
	ftpConnect(c.String("host"), c.String("port"), c.String("username"), c.String("password"))
	defer ftpClient.Quit()

	// inputs
	path := c.Args().Get(0)
	expected(path)

	// execute
	err := ftpClient.RemoveDir(path)
	checkErr(err)

	// outputs
	fmt.Println("path:" + path)
	return nil
}
