package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
)

func buckets() (buckets []cfenv.Service, err error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		return
	}
	services := appEnv.Services
	buckets, err = services.WithTag("s3")
	return
}

func bucketWithName(name string) (bucket *cfenv.Service, err error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		return
	}
	services := appEnv.Services
	bucket, err = services.WithName(name)
	return
}

func main() {

	app := cli.NewApp()
	app.Name = "migrate-s3-bucket"
	app.Usage = "Migrate objects from one S3 bucket to another"
	app.Action = func(c *cli.Context) {
	}

	app.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "perform migration inside CF container",
			Action: func(c *cli.Context) {
				buckets, err := buckets()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if len(c.Args()) != 2 {
					fmt.Println("USAGE: migrate-s3-bucket migrate <from> <to>")
					fmt.Println("  Using available buckets:")
					for _, bucket := range buckets {
						fmt.Println("    ", bucket.Name)
					}
					os.Exit(1)
				}

				fromBucket, err := bucketWithName(c.Args()[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				toBucket, err := bucketWithName(c.Args()[1])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("From bucket:", fromBucket.Credentials)
				fmt.Println("To bucket:", toBucket.Credentials)
			},
		},
		{
			Name:  "webserver",
			Usage: "simple dummy webserver",
			Action: func(c *cli.Context) {
				m := martini.Classic()
				m.Get("/", func() string {
					return "Now open `cf ssh` container and run `migrate-s3-bucket migrate`"
				})
				m.Run()
			},
		},
	}

	app.Run(os.Args)
}
