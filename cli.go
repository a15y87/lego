package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/xenolf/lego/acme"
)

// Logger is used to log errors; if nil, the default log.Logger is used.
var Logger *log.Logger

// logger is an helper function to retrieve the available logger
func logger() *log.Logger {
	if Logger == nil {
		Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return Logger
}

func main() {

	app := cli.NewApp()
	app.Name = "lego"
	app.Usage = "Let's encrypt client to go!"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Create and install a certificate",
			Action: run,
		},
		{
			Name:  "auth",
			Usage: "Create a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
		},
		{
			Name:  "install",
			Usage: "Install a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
		},
		{
			Name:  "revoke",
			Usage: "Revoke a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "certificate",
					Usage: "Revoke a specific certificate",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "Revoke all certs generated by the provided authorized key.",
				},
			},
		},
		{
			Name:  "rollback",
			Usage: "Rollback a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "checkpoints",
					Usage: "Revert configuration N number of checkpoints",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "domains, d",
			Usage: "Add domains to the process",
		},
		cli.StringFlag{
			Name:  "server, s",
			Value: "https://www.letsencrypt-demo.org/acme/new-reg",
			Usage: "CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client.",
		},
		cli.StringFlag{
			Name:  "authkey, k",
			Usage: "Path to the authorized key file",
		},
		cli.StringFlag{
			Name:  "email, m",
			Usage: "Email used for registration and recovery contact.",
		},
		cli.IntFlag{
			Name:  "rsa-key-size, B",
			Value: 2048,
			Usage: "Size of the RSA key.",
		},
		cli.BoolFlag{
			Name:  "no-confirm",
			Usage: "Turn off confirmation screens.",
		},
		cli.BoolFlag{
			Name:  "agree-tos, e",
			Usage: "Skip the end user license agreement screen.",
		},
		cli.StringFlag{
			Name:  "config-dir",
			Value: configDir,
			Usage: "Configuration directory.",
		},
		cli.StringFlag{
			Name:  "work-dir",
			Value: workDir,
			Usage: "Working directory.",
		},
		cli.StringFlag{
			Name:  "backup-dir",
			Value: backupDir,
			Usage: "Configuration backups directory.",
		},
		cli.StringFlag{
			Name:  "key-dir",
			Value: keyDir,
			Usage: "Keys storage.",
		},
		cli.StringFlag{
			Name:  "cert-dir",
			Value: certDir,
			Usage: "Certificates storage.",
		},
	}

	app.Run(os.Args)
}

func checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return nil
}

func run(c *cli.Context) {
	err := checkFolder(c.GlobalString("config-dir"))
	if err != nil {
		logger().Fatalf("Cound not check/create path: %v", err)
	}

	conf := NewConfiguration(c)

	//TODO: move to account struct? Currently MUST pass email.
	if !c.GlobalIsSet("email") {
		logger().Fatal("You have to pass an account (email address) to the program using --email or -m")
	}

	acc := NewAccount(c.GlobalString("email"), conf)
	client := acme.NewClient(c.GlobalString("server"), acc)
	if acc.Registration == nil {
		reg, err := client.Register()
		if err != nil {
			logger().Fatalf("Could not complete registration -> %v", err)
		}

		acc.Registration = reg
		acc.Save()

		logger().Print("!!!! HEADS UP !!!!")
		logger().Printf(`
			Your account credentials have been saved in your Let's Encrypt
			configuration directory at "%s".
			You should make a secure backup	of this folder now. This
			configuration directory will also contain certificates and
			private keys obtained from Let's Encrypt so making regular
			backups of this folder is ideal.

			If you lose your account credentials, you can recover
			them using the token
			"%s".
			You must write that down and put it in a safe place.`, c.GlobalString("config-dir"), reg.Body.Recoverytoken)
	}

	if !c.GlobalIsSet("domains") {
		logger().Fatal("Please specify --domains")
	}

}