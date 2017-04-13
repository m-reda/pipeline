package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
)

func fsUnits() []cli.Command {

	return []cli.Command{
		{Name: "copy", Usage: "Copy file or directory", Action: fsCopy},
		{Name: "move", Usage: "Move file or directory", Action: fsMove},
		{Name: "remove", Usage: "Remove file or directory", Action: fsRemove, Flags: []cli.Flag{cli.BoolFlag{Name: "rf"}}},
		{Name: "mkdir", Usage: "Make new directory", Action: fsMakeDir, Flags: []cli.Flag{cli.IntFlag{Name: "mode", Value: 0777}}},
		{Name: "mkfile", Usage: "Make new file", Action: fsMakeFile},
	}
}

// copy file or directory
func fsCopy(c *cli.Context) error {
	// inputs
	src, dest := c.Args().Get(0), c.Args().Get(1)
	expected(src, dest)

	// execute
	err := exec.Command("cp", "-rf", src, dest).Run()
	checkErr(err)

	// outputs
	fmt.Println("destination:" + dest)
	return nil
}

// move file or directory
func fsMove(c *cli.Context) error {
	// inputs
	src, dest := c.Args().Get(0), c.Args().Get(1)
	expected(src, dest)

	// execute
	err := exec.Command("mv", src, dest).Run()
	checkErr(err)

	// outputs
	fmt.Println("destination:" + dest)
	return nil
}

// delete file or directory
func fsRemove(c *cli.Context) error {
	// inputs
	path := c.Args().Get(0)
	expected(path)

	// execute
	args := []string{path}
	if c.Bool("rf") {
		args = []string{"-rf", path}
	}

	err := exec.Command("rm", args...).Run()
	checkErr(err)

	// outputs
	fmt.Println("path:" + path)
	return nil
}

// make directory
func fsMakeDir(c *cli.Context) error {
	// inputs
	path := c.Args().Get(0)
	expected(path)

	// execute
	err := os.MkdirAll(path, os.FileMode(c.Int("mode")))
	checkErr(err)

	// outputs
	fmt.Println("path:" + path)
	return nil
}

// make file
func fsMakeFile(c *cli.Context) error {
	// inputs
	path, content := c.Args().Get(0), c.Args().Get(1)
	expected(path)

	// execute
	var file, err = os.Create(path)
	checkErr(err)

	defer file.Close()
	file.Write([]byte(content))

	// outputs
	fmt.Println("path:" + path)
	return nil
}
