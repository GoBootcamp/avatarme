package main

import (
  "fmt"
  "os"
  "github.com/codegangsta/cli"
  person "github.com/duranmla/avatarme/user"
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
        name, email := person.RequestCredentials();
        user := person.New(name, email)
				fmt.Println(user)
			},
		},
	}

	app.Run(os.Args)
}
