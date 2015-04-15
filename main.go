package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"log"
	"os"
	"os/exec"
	"text/template"

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
	config LaunchConfig
	enc    *json.Encoder
)

// cli flags
var (
	podname            string
	cli                bool
	info               bool
	jsonout            bool
	validatemasterrole bool
)

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

func main() {
	flag.StringVar(&podname, "podname", "", "Name of the pod")
	flag.BoolVar(&cli, "cli", false, "launch redis-cli to connect to the pod")
	flag.BoolVar(&jsonout, "jsonout", false, "output info in JSON format (requires -info)")
	flag.BoolVar(&info, "info", false, "output info")
	flag.Parse()

	var pod parser.PodConfig
	var err error
	if podname == "" {
		log.Print("Need a podname. Try using '-podname <podname>'")
		flag.PrintDefaults()
		return
	}
	if jsonout {
		enc = json.NewEncoder(os.Stdout)
	}

	if config.UseSentinelConfig {
		log.Print("using sentinel config file")
		sentinel, err := parser.ParseSentinelConfig(config.SentinelConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		pod, err = sentinel.GetPod(podname)
		if err != nil {
			log.Fatal(err)
		}
	}
	if config.UseRedSkull {
		log.Print("Using RedSkull connection")
		pod, err = getPodInfoFromRedSkull(podname)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	}
	if cli {
		//args := flag.Args()
		log.Printf("Connecting to pod '%s' via CLI", podname)
		cmd := exec.Command("redis-cli", "-h", pod.MasterIP, "-p", pod.MasterPort, "-a", pod.Authpass)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	if info {
		if jsonout {
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
}
