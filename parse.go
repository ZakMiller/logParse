package main

import (
	"fmt"
	"os"
	"regexp"
	"io/ioutil"
	"strings"
	"encoding/csv"
)

func filterLines(path string) ([][]byte, bool){
	text, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Printf("Got an error while reading the file %v\n", err)
		return nil, false
	}

	lineCount := len(text)

	r, _ := regexp.Compile(`\w+\s+\w+.\w+\s+milliseconds`)

	relevantLines := r.FindAll(text, lineCount)

	return relevantLines, true
}

func getTimes(lines [][]byte) map[string][]string {
	m := make(map[string][]string)

	for _, line := range lines {
		s := string(line)
		fields := strings.Fields(s)
		if len(fields) == 3 {
			function := fields[0]
			time := fields[1]
			m[function] = append(m[function], time)
		}
	}
	return m
}

func getKeys(m map[string][]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func writeValues(m map[string][]string, w *csv.Writer, keys []string) {
	functionCount := len(keys)
	callCount := len(m[keys[0]])

	for i := 0; i < callCount; i++ {
		
		record := make([]string, functionCount)
		for j, key := range keys {
			if i < len(m[key]) {
				record[j] = m[key][i]
			}
		}
		(*w).Write(record)
	}
}

func createCSV(m map[string][]string) {
	keys := getKeys(m)

	w := csv.NewWriter(os.Stdout)
	w.Write(keys)
	writeValues(m, w, keys)
	w.Flush()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify the file you want to parse")
		return
	}

	path := os.Args[1]	
	lines, ok := filterLines(path)

	if !ok {
		return
	}

	m := getTimes(lines)
	createCSV(m)
}
