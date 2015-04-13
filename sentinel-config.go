package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// getPodInfoFromConfig returns a PodConfig struct with information about
// the found pod, error and an empty struct otherwise
func getPodInfoFromConfig(podname string) (PodConfig, error) {
	return getPodInfoFromConfigFile(config.SentinelConfigFile, podname)
}

// getPodInfoFromConfigFile actually parses the config file and is intended to
// be called by getPodInfoFromConfig
func getPodInfoFromConfigFile(configfile string, podname string) (PodConfig, error) {
	var pod PodConfig
	file, err := os.Open(configfile)
	defer file.Close()
	if err != nil {
		log.Print(err)
		return pod, err
	}
	bf := bufio.NewReader(file)
	foundpod := false
	for {
		rawline, err := bf.ReadString('\n')
		if err == nil || err == io.EOF {
			line := strings.TrimSpace(rawline)
			// ignore comments
			if strings.Contains(line, "#") {
				continue
			}
			if rawline == "" {
				break
			}
			if strings.Contains(rawline, podname) {
				entries := strings.Split(line, " ")
				if len(entries) == 0 {
					continue
				}
				if entries[2] == podname {
					foundpod = true
					//Most values are key/value pairs
					switch entries[1] {
					case "monitor": // Have a sentinel directive, initialze struct
						//pod.KnownSentinels = make([]string, 5)
						//pod.KnownSlaves = make([]string, 12)
						pod.Settings = make(map[string]string)
						pod.Name = podname
						pod.MasterIP = entries[3]
						pod.MasterPort = entries[4]
						continue
					case "known-slave":
						pod.KnownSlaves = append(pod.KnownSlaves, entries[3])
						continue
					case "known-sentinel":
						pod.KnownSentinels = append(pod.KnownSentinels, entries[3])
						continue
					case "auth-pass":
						pod.Authpass = entries[3]
						continue
					case "config-epoch", "leader-epoch", "current-epoch", "down-after-milliseconds", "maxclients", "parallel-syncs", "can-failover", "failover-timeout":
						pod.Settings[entries[1]] = entries[3]
						continue
					default:
						fmt.Printf("Unhandled directive line: %+v\n", rawline)
						log.Printf("entries: %+v", entries)
						continue
					}
				} else {
					continue
				}
			} else {
				// sentinel ensures all entries for a pod are together. So if
				// we've parsed entries for a pod, but now we read a line which
				// doesn't have the pod name in it, we're doen with it. No need
				// to continue reading the file.
				if foundpod {
					break
				}
			}
		} else {
			break
		}
	}
	if !foundpod {
		return pod, fmt.Errorf("Pod '%s' not found in config file", podname)
	}
	return pod, nil
}
