package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Pipeline builtin units cli"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{Name: "fs", Usage: "Filesystem units", Subcommands: fsUnits()},
		{Name: "git", Usage: "Git units", Subcommands: gitUnits()},
		{Name: "func", Usage: "Functions units", Subcommands: funcUnits()},
		{Name: "ftp", Usage: "Ftp units", Subcommands: ftpUnits()},
	}

	app.Run(os.Args)
}

func expected(n ...string) {
	for _, v := range n {
		if v == "" {
			log.Fatalf("Error: Inputs not enough, expected %d", len(n))
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
}

func pwd() string {
	path, err := os.Getwd()
	checkErr(err)

	return path
}
