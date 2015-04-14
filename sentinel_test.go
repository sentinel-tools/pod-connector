package main

import (
	"strings"
	"testing"

	"github.com/sentinel-tools/sconf-parser"
)

func TestValidPodFromSentinelConfig(t *testing.T) {
	conf, err := parser.ParseSentinelConfig("sentinel.conf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	_, err = conf.GetPod("pod1")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestINValidPodFromSentinelConfig(t *testing.T) {
	conf, err := parser.ParseSentinelConfig("sentinel.conf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	_, err = conf.GetPod("pod2")
	if err == nil {
		t.Error(err)
		t.Fail()
	} else {
		if !strings.Contains(err.Error(), "not found") {
			t.Error("Somehow found a pod which doesn't exist!")
			t.Fail()
		}
	}
}
