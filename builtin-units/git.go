package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os/exec"
	"strings"
)

func gitUnits() []cli.Command {
	return []cli.Command{
		{Name: "init", Usage: "Initial a new repository", Action: gitInit},
		{Name: "add", Usage: "Add the changes", Action: gitAdd},
		{Name: "commit", Usage: "Commit the changes", Action: gitCommit},
		{Name: "push", Usage: "Push to remote server", Action: gitPush},
		{Name: "checkout", Usage: "Change the branch", Action: gitCheckout},
		{Name: "clone", Usage: "Clone a remote repository", Action: gitClone},
		{Name: "merge", Usage: "Merge a branch into the current branch", Action: gitMerge},
		{Name: "pull", Usage: "Pull upstream changes into the local repository", Action: gitPull},
		{Name: "remote_add", Usage: "Pull upstream changes into the local repository", Action: gitRemoteAdd},
	}
}

// initial a new repository
func gitInit(_ *cli.Context) error {

	err := exec.Command("git", "init").Run()
	checkErr(err)
	return nil
}

// add the changes
func gitAdd(c *cli.Context) error {
	// inputs
	files := c.Args().Get(0)
	expected(files)

	// execute
	err := exec.Command("git", "add", files).Run()
	checkErr(err)

	return nil
}

// commit the changes
func gitCommit(c *cli.Context) error {
	// inputs
	message := c.Args().Get(0)
	expected(message)

	// execute
	err := exec.Command("git", "commit", "-m", message).Run()
	checkErr(err)

	return nil
}

// push to remote repository
func gitPush(c *cli.Context) error {
	// inputs
	remote, branch := c.Args().Get(0), c.Args().Get(1)
	expected(remote, branch)

	// execute
	err := exec.Command("git", "push", remote, branch).Run()
	checkErr(err)

	// outputs
	fmt.Println("remote:" + remote)
	fmt.Println("branch:" + branch)
	return nil
}

// add the changes
func gitCheckout(c *cli.Context) error {
	// inputs
	branch := c.Args().Get(0)
	expected(branch)

	// execute
	err := exec.Command("git", "checkout", branch).Run()
	checkErr(err)

	fmt.Println("branch:" + branch)
	return nil
}

// clone repository
func gitClone(c *cli.Context) error {
	// inputs
	repo := c.Args().Get(0)
	expected(repo, repo)

	// execute
	err := exec.Command("git", "clone", repo).Run()
	checkErr(err)

	a := strings.Split(repo, "/")
	repoName := strings.Replace(a[len(a)-1], ".git", "", 1)

	// outputs
	fmt.Println("repository:" + repo)
	fmt.Println("local:" + pwd() + "/" + repoName)
	return nil
}

// merge a branch
func gitMerge(c *cli.Context) error {
	// inputs
	branch := c.Args().Get(0)
	expected(branch)

	// execute
	err := exec.Command("git", "merge", branch).Run()
	checkErr(err)

	// outputs
	fmt.Println("branch:" + branch)
	return nil
}

// pull changes
func gitPull(c *cli.Context) error {
	// inputs
	remote := c.Args().Get(0)
	expected(remote)

	// execute
	err := exec.Command("git", "pull", remote).Run()
	checkErr(err)

	// outputs
	fmt.Println("remote:" + remote)
	return nil
}

// add new remote
func gitRemoteAdd(c *cli.Context) error {
	// inputs
	name, url := c.Args().Get(0), c.Args().Get(1)
	expected(name, url)

	// execute
	err := exec.Command("git", "remote", "add", name, url).Run()
	checkErr(err)

	// outputs
	fmt.Println("name:" + name)
	fmt.Println("url:" + url)
	return nil
}
