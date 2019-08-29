package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/exercism/go-test-runner/exreport"
)

const (
	statPass = "pass"
	statFail = "fail"
	statSkip = "skip"
	statErr  = "error"
)

func main() {
	lines, err := readFile("test.out")
	if err != nil {
		log.Panic(err)
	}

	output := getStructure(lines)
	bts, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(bts))
}

func readFile(file string) ([][]byte, error) {
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(bts) < 5 {
		return nil, errors.New("file is empty")
	}
	return bytes.Split(bts, []byte{'\n'}), nil
}

func getStructure(lines [][]byte) *exreport.Report {
	var tests = map[string]*exreport.Test{}
	for _, lineBytes := range lines {
		var line testLine
		if err := json.Unmarshal(lineBytes, &line); err != nil {
			log.Println(err)
			continue
		}

		if line.Test == "" {
			continue
		}

		switch line.Action {
		case "run":
			tests[line.Test] = &exreport.Test{
				Name:   line.Test,
				Status: statSkip,
			}
		case "output":
			tests[line.Test].Message += "\n" + line.Output
		case statFail:
			tests[line.Test].Status = statFail
		case statPass:
			tests[line.Test].Status = statPass
		}
	}

	report := &exreport.Report{
		Status: statPass,
		Tests:  nil,
	}

	for _, test := range tests {
		if test == nil {
			// just to be sure we dont get a nil pointer exception
			continue
		}
		if test.Status == statSkip {
			report.Status = statErr
		}
		if report.Status == statPass && test.Status == statFail {
			report.Status = statFail
		}

		report.Tests = append(report.Tests, *test)
	}

	return report
}

type testLine struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Elapsed float64
	Output  string
}