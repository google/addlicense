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
	"os"
	"testing"
	"text/template"
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

func TestExecuteTemplate(t *testing.T) {
	tests := []struct {
		name          string
		template      string
		data          licenseData
		top, mid, bot string
		want          string
	}{
		{
			"empty template",
			"",
			licenseData{},
			"", "", "",
			"\n",
		},
		{
			"no extra",
			"{{.Holder}}{{.Year}}{{.SPDXID}}",
			licenseData{Holder: "H", Year: "Y", SPDXID: "S"},
			"", "", "",
			"HYS\n\n",
		},
		{
			"only mid",
			"{{.Holder}}{{.Year}}{{.SPDXID}}",
			licenseData{Holder: "H", Year: "Y", SPDXID: "S"},
			"", "// ", "",
			"// HYS\n\n",
		},
		{
			"top, mid, bot",
			"{{.Holder}}{{.Year}}{{.SPDXID}}",
			licenseData{Holder: "H", Year: "Y", SPDXID: "S"},
			"/*", " * ", "*/",
			"/*\n * HYS\n*/\n\n",
		},

		// ensure we don't escape HTML characters by using the wrong template package
		{
			"html chars",
			"{{.Holder}}",
			licenseData{Holder: "A&Z"},
			"", "", "",
			"A&Z\n\n",
		},

		// empty near should not add a space
		{
			"no year, apache",
			tmplApache,
			licenseData{Holder: "Holder"},
			"", "", "",
			`Copyright Holder

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

`,
		},
		{
			"no year, BSD",
			tmplBSD,
			licenseData{Holder: "Holder"},
			"", "", "",
			`Copyright (c) Holder All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.

`,
		},
		{
			"no year, MIT",
			tmplMIT,
			licenseData{Holder: "Holder"},
			"", "", "",
			`Copyright (c) Holder

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

`,
		},
		{
			"no year, SPDX",
			tmplSPDX,
			licenseData{Holder: "Holder", SPDXID: "Spdx"},
			"", "", "",
			`Copyright Holder
SPDX-License-Identifier: Spdx

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpl, err := template.New("").Parse(tt.template)
			if err != nil {
				t.Errorf("error parsing template: %v", err)
			}
			got, err := executeTemplate(tpl, tt.data, tt.top, tt.mid, tt.bot)
			if err != nil {
				t.Errorf("returned error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("returned \n%q\n, want: \n%q", string(got), tt.want)
			}
		})
	}
}
