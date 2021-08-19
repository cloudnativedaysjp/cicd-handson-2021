package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"runtime"
	"testing"
)

func init() {
	stars = make(map[string]int64)
	stars["https://github.com/argoproj/argo-workflows"] = 5000
}

func Test_parseFile(t *testing.T) {
	inputyaml, _ := ioutil.ReadFile("../test/data/input.yaml")
	outputjson, _ := ioutil.ReadFile("../test/data/output.json")

	projects, err := findCicdProjects(inputyaml)
	if err != nil {
		t.Errorf("Operation failed: %s", err)
	}
	v, _ := json.MarshalIndent(projects, "", "  ")

	if runtime.GOOS == "windows" {
		outputjson = bytes.ReplaceAll(outputjson, []byte("\r"), []byte(""))
	}

	if bytes.Compare(v, outputjson) != 0 {
		t.Errorf("Invalid output.\noutput:\n%s\n\nexpected:\n%s", string(v), string(outputjson))
	}
}
