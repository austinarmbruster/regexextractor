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

package extractor

import (
	"errors"
	"regexp"
)

var (
	// ErrMissingName is returned when a user attempts to remove a named pattern
	// that does not exist.
	ErrMissingName = errors.New("missing named pattern")
)

// RegexExtractor is a simple text processor for finding subsets of text that
// match a collection of regular expressions.
type RegexExtractor struct {
	labelPatterns map[string]*regexp.Regexp
}

// NewRegex creates a valid RegexExtractor.
func NewRegex() *RegexExtractor {
	return &RegexExtractor{
		labelPatterns: make(map[string]*regexp.Regexp),
	}
}

// Add will include another regular expression for extraction of values.  The
// name is the label that is provided in extraction output.  If the pattern is
// invalid, then an error is returned indicating the issue with the pattern.
func (r *RegexExtractor) Add(name, pattern string) error {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	r.labelPatterns[name] = reg
	return nil

}

// Remove will either return ErrMissingName when the name does not exist or nil
// on success.
func (r *RegexExtractor) Remove(name string) error {
	if _, ok := r.labelPatterns[name]; !ok {
		return ErrMissingName
	}

	delete(r.labelPatterns, name)
	return nil
}

// Extract finds matches in the provided text for each of the configured
// patterns.  The label provided with the pattern is used as the label for the
// slice of matching strings found in the text.
func (r *RegexExtractor) Extract(text string) (map[string][]string, error) {
	rtnVal := make(map[string][]string)

	for l, re := range r.labelPatterns {
		hits := re.FindAllString(text, -1)
		if hits == nil {
			continue
		}

		rtnVal[l] = hits
	}
	return rtnVal, nil
}
