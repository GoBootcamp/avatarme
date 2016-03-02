package main

import (
  "fmt"
  "os"
  "github.com/codegangsta/cli"
  "github.com/duranmla/avatarme/avatar"
  "github.com/duranmla/clirescue/cmdutil"
)

var (
  Stdout *os.File = os.Stdout
)

func main() {
	app := cli.NewApp()

	app.Name = "avatarme"
	app.Usage = "CLI tool to generate hashes from your email"

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "prints out a User struct with hash",
			Action: func(c *cli.Context) {
        name, email := requestCredentials();
        user := avatar.New(name, email)
				fmt.Println(user)
			},
		},
	}

	app.Run(os.Args)
}

func requestCredentials() (name, email string){
  fmt.Fprint(Stdout, "name: ")
  name = cmdutil.ReadLine()
  fmt.Fprint(Stdout, "email: ")
  email = cmdutil.ReadLine()

  return name, email
}
