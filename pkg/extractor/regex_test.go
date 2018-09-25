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

package extractor_test

import (
	"testing"

	"github.com/austinarmbruster/regexextractor/pkg/extractor"
)

func TestRemoveNoPatterns(t *testing.T) {
	r := extractor.NewRegex()

	result := r.Remove("dne")
	if result != extractor.ErrMissingName {
		t.Errorf("Should get an error when removing a non-existent name:  have:  %v, want, %v", result, extractor.ErrMissingName)
	}
}

func TestAddRemove(t *testing.T) {

	var cases = []struct {
		testName     string
		pattName     string
		patt         string
		expectAddErr bool
	}{
		{"invalid", "invalid", "*", true},
		{"star", "all", ".*", false},
		{"email", "email", ".*@.*", false},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			r := extractor.NewRegex()
			addResult := r.Add(tc.pattName, tc.patt)

			if tc.expectAddErr {
				if addResult == nil {
					t.Errorf("succeeded in adding an invalid pattern:  pattern: %v, err:  %v", tc.patt, addResult)
				}
				return
			}

			if addResult != nil {
				t.Errorf("Failed to add a valid pattern:  pattern %v, err: %v", tc.patt, addResult)
				return
			}

			rmResult := r.Remove(tc.pattName)
			if rmResult != nil {
				t.Errorf("Failed to remove the pattern:  patt:  %v, err:  %v", tc.patt, rmResult)
			}
		})
	}
}

type pair struct {
	label string
	patt  string
}

func TestAddExtract(t *testing.T) {
	var cases = []struct {
		name     string
		patterns []pair
		input    string
		expected map[string][]string
	}{
		{"all", []pair{{"all", ".*"}}, "rabbit jumped", map[string][]string{"all": []string{"rabbit jumped"}}},
		{"simpleEmail", []pair{{"email", ".*@.*"}}, "me@example.com", map[string][]string{"email": []string{"me@example.com"}}},
		{"emailPlus", []pair{{"email", "[a-z]*@[a-z]*\\.[a-z]*"}}, "My email address is me@example.com today", map[string][]string{"email": []string{"me@example.com"}}},
		{"complexEmailPatt", []pair{{"email", emailPatt}}, "My email address is me@example.com today", map[string][]string{"email": []string{"me@example.com"}}},
		{"nameEmail", []pair{{"name", "[A-Z][a-z]*"}, {"email", emailPatt}}, "Bill Gates bill@ms.com", map[string][]string{"email": []string{"bill@ms.com"}, "name": []string{"Bill", "Gates"}}},
		{"multiRuleOneHit", []pair{{"name", "[A-Z][a-z]*"}, {"email", emailPatt}}, "bill@ms.com", map[string][]string{"email": []string{"bill@ms.com"}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := extractor.NewRegex()

			for _, p := range tc.patterns {
				r.Add(p.label, p.patt)
			}

			result, err := r.Extract(tc.input)
			if err != nil {
				t.Errorf("Failed to extract a value: pattern:  %v, input: %v, err: %v", tc.patterns, tc.input, err)
				return
			}

			if result == nil && tc.expected != nil {
				t.Errorf("failed extract:  have: %v, want %v", result, tc.expected)
			}

			var expectedV []string
			var ok bool
			for k, v := range result {
				if expectedV, ok = tc.expected[k]; !ok {
					t.Errorf("Extracted an unexpected value: key:  %v, values:  %v", k, v)
					return
				}

				if len(v) != len(expectedV) {
					t.Errorf("Expected length and extracted length vary: have: %v, want: %v", v, expectedV)
				}

				compare(t, v, expectedV)
			}
		})
	}
}

// Since this examples targets just using the standard library, a simple array
// comparitor was needed.
func compare(t *testing.T, have, want []string) {
	wantMap := make(map[string]bool)
	for _, w := range want {
		wantMap[w] = true
	}

	for _, h := range have {
		if !wantMap[h] {
			t.Errorf("Extracted extra value:  have: %v, want: %v", h, nil)
		}

		delete(wantMap, h)
	}

	for w := range wantMap {
		t.Errorf("Failed to extract a value:  have: %v, want: %v", nil, w)
	}
}

const (
	// The following pattern is from https://emailregex.com/ on 24 Sept 2018.
	emailPatt = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"
)
