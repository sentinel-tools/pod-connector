package main

import "fmt"
import "github.com/sentinel-tools/sconf-parser"

func getPodInfoFromRedSkull(podname string) (parser.PodConfig, error) {
	var pod parser.PodConfig
	err := fmt.Errorf("Redskull support not yet implemented")
	return pod, err
}
