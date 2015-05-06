package main

import (
	"encoding/json"
	"fmt"

	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/kelseyhightower/envconfig"

	parser "github.com/sentinel-tools/sconf-parser"
)

type LaunchConfig struct {
	SentinelConfigFile string
	RedSkullAddress    string
	ValidateNodes      bool
	UseRedSkull        bool
	UseSentinelConfig  bool
}

var (
	pod      parser.PodConfig
	err      error
	config   LaunchConfig
	enc      *json.Encoder
	podname  string
	cfile    string
	app      *cli.App
	sentinel parser.SentinelConfig
)

// cli flags

func SetConfigFile(c *cli.Context) error {
	cfile = c.GlobalString("sentinelconfig")
	sentinel, err = parser.ParseSentinelConfig(config.SentinelConfigFile)
	return err
}
func init() {
	err := envconfig.Process("podconnector", &config)
	if err != nil {
		log.Fatal(err)
	}
	// If we specify a source of pod info, set that source as what we want to use.
	if config.RedSkullAddress > "" {
		config.UseRedSkull = true
	}
	if config.SentinelConfigFile > "" {
		config.UseSentinelConfig = true
	}

	// now, set defaults for the source selected
	if config.UseRedSkull {
		if config.RedSkullAddress == "" {
			if config.UseRedSkull {
				config.RedSkullAddress = "localhost:8001"
			}
		}
	}

	if config.UseSentinelConfig {
		if config.SentinelConfigFile == "" {
			config.SentinelConfigFile = "/etc/redis/sentinel.conf"
		}
	}
	if !(config.UseSentinelConfig || config.UseRedSkull) {
		config.SentinelConfigFile = "/etc/redis/sentinel.conf"
		config.UseSentinelConfig = true
	}

}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func beforePodCommand(c *cli.Context) (err error) {
	//podname := c.String("name")
	args := c.Args()
	if len(args) == 0 {
		log.Fatal("Need a podname as first argument")
	}
	podname = args[0]
	pod, _ = sentinel.GetPod(podname)
	return nil
}
func main() {
	app = cli.NewApp()
	app.Name = "pod-connector"
	app.Usage = "Interact with a Sentinel using configuration data"
	app.Version = "0.5.1"
	app.EnableBashCompletion = true
	author := cli.Author{Name: "Bill Anderson", Email: "therealbill@me.com"}
	app.Authors = append(app.Authors, author)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "sentinelconfig, s",
			Value: "/etc/redis/sentinel.conf",
			Usage: "Location of the sentinel config file",
		},
	}
	app.Before = SetConfigFile

	app.Commands = []cli.Command{
		{
			Name:   "cli",
			Usage:  "connect redis-cli to the pod",
			Action: PodCli,
			Before: beforePodCommand,
		},
		{
			Name:   "info",
			Usage:  "Show Sentinel information on the pod",
			Action: ShowInfo,
			Before: beforePodCommand,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json, j",
					Usage: "Output only JSON",
				},
			},
		},
	}
	enc = json.NewEncoder(os.Stdout)
	app.Run(os.Args)
}

func ShowInfo(c *cli.Context) {
	if c.Bool("json") {
		if err := enc.Encode(&pod); err != nil {
			log.Println(err)
		}
	} else {
		t := template.Must(template.New("podinfo").Parse(PodInfoTemplate))
		err := t.Execute(os.Stdout, pod)
		if err != nil {
			log.Println("executing template:", err)
		}
		fmt.Printf("cli string: redis-cli -h %s -p %s -a %s\n", pod.MasterIP, pod.MasterPort, pod.Authpass)
	}
	return
}

func PodCli(c *cli.Context) {
	log.Printf("Connecting to pod '%s' via CLI", podname)
	cmd := exec.Command("redis-cli", "-h", pod.MasterIP, "-p", pod.MasterPort, "-a", pod.Authpass)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
