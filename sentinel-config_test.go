package main

import (
	"strings"
	"testing"
)

func TestValidPodFromSentinelConfig(t *testing.T) {
	_, err := getPodInfoFromConfigFile("sentinel.conf", "pod1")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestINValidPodFromSentinelConfig(t *testing.T) {
	_, err := getPodInfoFromConfigFile("sentinel.conf", "pod2")
	if err == nil {
		t.Error(err)
		t.Fail()
	} else {
		if !strings.Contains(err.Error(), "not found in config file") {
			t.Error("Somehow found a pod which doesn't exist!")
			t.Fail()
		}
	}
}
