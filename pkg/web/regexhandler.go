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

package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/austinarmbruster/regexextractor/pkg/extractor"
)

// RegexHandler configures an HTTP handler to provide extracts of received text
// with the labels given.
func RegexHandler(labeledPatterns map[string]string) (http.Handler, error) {
	// create the "service"
	regEx := extractor.NewRegex()
	for l, p := range labeledPatterns {
		err := regEx.Add(l, p)
		if err != nil {
			return nil, err
		}
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read content", http.StatusInternalServerError)
			return
		}

		extracts, err := regEx.Extract(string(content))
		if err != nil {
			http.Error(w, "failed to read content", http.StatusBadRequest)
			return
		}

		jBytes, err := json.Marshal(extracts)
		if err != nil {
			http.Error(w, "failed to produce the proper output", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jBytes)
	})
	return h, nil
}
