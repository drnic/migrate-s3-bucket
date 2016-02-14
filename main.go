package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
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
				// toBucket, err := bucketWithName(c.Args()[1])
				// if err != nil {
				// 	fmt.Println(err)
				// 	os.Exit(1)
				// }
				fromBucketAWS, _ := aws.GetAuth(fromBucket.Credentials["access_key_id"].(string), fromBucket.Credentials["secret_access_key"].(string))
				// toBucketAWS, _ := aws.GetAuth(toBucket.Credentials["access_key_id"].(string), toBucket.Credentials["secret_access_key"].(string))
				fmt.Println("From bucket:", fromBucketAWS)
				client := s3.New(fromBucketAWS, aws.USEast)
				bucket := client.Bucket(fromBucket.Credentials["bucket"].(string))
				items, err := bucket.GetBucketContents()
				if err != nil {
					log.Fatal(err)
				}

				for key, val := range *items {
					fmt.Printf("'%s': %#v\n", key, val)
				}

				// fmt.Println("To bucket:", toBucketAWS)
				// client = s3.New(toBucketAWS, aws.USEast)
				// bucket, err = client.Bucket()
				//
				// if err != nil {
				// 	log.Fatal(err)
				// }
				//
				// log.Print(fmt.Sprintf("%T %+v", resp.Buckets[0], resp.Buckets[0]))
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
