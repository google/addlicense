// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"html/template"
	"os"
	"testing"
)

func init() {
	// ensure that pre-defined templates must parse
	template.Must(template.New("").Parse(tmplApache))
	template.Must(template.New("").Parse(tmplMIT))
	template.Must(template.New("").Parse(tmplBSD))
	template.Must(template.New("").Parse(tmplMPL))
}

func TestFetchTemplate(t *testing.T) {
	tests := []struct {
		description  string   // test case description
		license      string   // license passed to fetchTemplate
		templateFile string   // templatefile passed to fetchTemplate
		spdx         spdxFlag // spdx value passed to fetchTemplate
		wantTemplate string   // expected returned template
		wantErr      error    // expected returned error
	}{
		// custom template files
		{
			"non-existent template file",
			"",
			"/does/not/exist",
			spdxOff,
			"",
			os.ErrNotExist,
		},
		{
			"custom template file",
			"",
			"testdata/custom.tpl",
			spdxOff,
			"Copyright {{.Year}} {{.Holder}}\n\nCustom License Template\n",
			nil,
		},

		{
			"unknown license",
			"unknown",
			"",
			spdxOff,
			"",
			errors.New(`unknown license: "unknown". Include the '-s' flag to request SPDX style headers using this license`),
		},

		// pre-defined license templates, no SPDX
		{
			"apache license template",
			"Apache-2.0",
			"",
			spdxOff,
			tmplApache,
			nil,
		},
		{
			"mit license template",
			"MIT",
			"",
			spdxOff,
			tmplMIT,
			nil,
		},
		{
			"bsd license template",
			"bsd",
			"",
			spdxOff,
			tmplBSD,
			nil,
		},
		{
			"mpl license template",
			"MPL-2.0",
			"",
			spdxOff,
			tmplMPL,
			nil,
		},

		// SPDX variants
		{
			"apache license template with SPDX added",
			"Apache-2.0",
			"",
			spdxOn,
			tmplApache + spdxSuffix,
			nil,
		},
		{
			"apache license template with SPDX only",
			"Apache-2.0",
			"",
			spdxOnly,
			tmplSPDX,
			nil,
		},
		{
			"unknown license with SPDX only",
			"unknown",
			"",
			spdxOnly,
			tmplSPDX,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tpl, err := fetchTemplate(tt.license, tt.templateFile, tt.spdx)
			if tt.wantErr != nil && (err == nil || (!errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error())) {
				t.Fatalf("fetchTemplate(%q, %q) returned error: %#v, want %#v", tt.license, tt.templateFile, err, tt.wantErr)
			}
			if tpl != tt.wantTemplate {
				t.Errorf("fetchTemplate(%q, %q) returned template: %q, want %q", tt.license, tt.templateFile, tpl, tt.wantTemplate)
			}
		})
	}
}
