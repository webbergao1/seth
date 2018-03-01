package main

import (
	"fmt"
	"os"
	"runtime"

	"seth/log"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Info("Processing. Please wait....")
	args := append(os.Args, "newaccount")

	app := cli.NewApp()
	app.Name = "seth"
	app.Version = "1.0.1"
	app.Author = "seth"
	app.Email = ""
	app.Usage = "the seth command line interface"

	app.Commands = Commands()

	err := app.Run(args) //os.Args

	if err != nil {
		fmt.Printf(err.Error())
	}

}

// Commands is
func Commands() []cli.Command {

	n := &NodeCli{}

	return []cli.Command{
		cli.Command{
			Name:      "help",
			Before:    n.init,
			Action:    n.Help,
			ShortName: "h",
			Usage:     "Help for cmd",
		},
		cli.Command{
			Name:      "newaccount",
			Before:    n.init,
			Action:    n.NewAccount,
			ShortName: "n",
			Usage:     "new account return the account address&privatekey",
		},
		cli.Command{
			Name:      "start",
			Before:    n.init,
			Action:    n.Start,
			ShortName: "s",
			Usage:     "Start the seth node",
		},
		cli.Command{
			Name:      "clear",
			Before:    n.init,
			Action:    n.Clear,
			ShortName: "c",
			Usage:     "Clear the seth node",
		},
	}
}
