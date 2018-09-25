// The MIT License (MIT)
//
// Copyright © 2018 Austin Armbruster <austin@iits.me>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/austinarmbruster/regexextractor/pkg/web"
)

// Simple HTTP Server that parses text and provides the label extracts of that
// text based upon the named patterns in a CSV file.
func main() {
	csvFileName := flag.String("file", "", "CSV File Name")
	addr := flag.String("addr", ":8080", "server address")

	flag.Parse()

	if *csvFileName == "" {
		log.Fatal("Missing the CSV file name")
	}

	patterns, err := readCSV(*csvFileName)
	if err != nil {
		log.Fatalf("Failed to read the CSV file:  %v", err)
	}

	h, err := web.RegexHandler(patterns)
	http.ListenAndServe(*addr, h)
}

// readCSV provides the collection of name patterns.
func readCSV(fileName string) (map[string]string, error) {
	rtnVal := make(map[string]string)

	cf, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the CSV file:  file name: %v, err: %v", fileName, err)
	}
	defer cf.Close()

	r := csv.NewReader(cf)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 2 {
			return nil, fmt.Errorf("Missing the name / pattern: have: %v", record)
		}

		rtnVal[record[0]] = record[1]
	}

	return rtnVal, nil
}
