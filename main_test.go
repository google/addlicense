// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func run(t *testing.T, name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
}

func tempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "addlicense")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestInitial(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	run(t, "cp", "-r", "testdata/initial", tmp)

	// run at least 2 times to ensure the program is idempotent
	for i := 0; i < 2; i++ {
		t.Logf("run #%d", i)
		targs := []string{"-test.run=TestInitial"}
		cargs := []string{"-l", "apache", "-c", "Google LLC", "-y", "2018", tmp}
		c := exec.Command(os.Args[0], append(targs, cargs...)...)
		c.Env = []string{"RUNME=1"}
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}

		run(t, "diff", "-r", filepath.Join(tmp, "initial"), "testdata/expected")
	}
}

func TestMultiyear(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")
	const sampleLicensed = "testdata/multiyear_file.c"

	run(t, "cp", "testdata/initial/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestMultiyear",
		"-l", "bsd", "-c", "Google LLC",
		"-y", "2015-2017,2019", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	run(t, "diff", samplefile, sampleLicensed)
}

func TestWriteErrors(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/initial/file.c", samplefile)
	run(t, "chmod", "0444", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestWriteErrors",
		"-l", "apache", "-c", "Google LLC", "-y", "2018",
		samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	out, err := cmd.CombinedOutput()
	if err == nil {
		run(t, "chmod", "0644", samplefile)
		t.Fatalf("TestWriteErrors exited with a zero exit code.\n%s", out)
	}
	run(t, "chmod", "0644", samplefile)
}

func TestReadErrors(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/initial/file.c", samplefile)
	run(t, "chmod", "a-r", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestReadErrors",
		"-l", "apache", "-c", "Google LLC", "-y", "2018",
		samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	out, err := cmd.CombinedOutput()
	if err == nil {
		run(t, "chmod", "0644", samplefile)
		t.Fatalf("TestWriteErrors exited with a zero exit code.\n%s", out)
	}
	run(t, "chmod", "0644", samplefile)
}

func TestCheckSuccess(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/expected/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestCheckSuccess",
		"-l", "apache", "-c", "Google LLC", "-y", "2018",
		"-check", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}

func TestCheckFail(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/initial/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestCheckFail",
		"-l", "apache", "-c", "Google LLC", "-y", "2018",
		"-check", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("TestCheckFail exited with a zero exit code.\n%s", out)
	}
}

func TestMPL(t *testing.T) {
	if os.Getenv("RUNME") != "" {
		main()
		return
	}

	tmp := tempDir(t)
	t.Logf("tmp dir: %s", tmp)
	samplefile := filepath.Join(tmp, "file.c")

	run(t, "cp", "testdata/expected/file.c", samplefile)
	cmd := exec.Command(os.Args[0],
		"-test.run=TestMPL",
		"-l", "mpl", "-c", "Google LLC", "-y", "2018",
		"-check", samplefile,
	)
	cmd.Env = []string{"RUNME=1"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}

func createTempFile(contents string, pattern string) (*os.File, error) {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(f.Name(), []byte(contents), 0644); err != nil {
		return nil, err
	}

	return f, nil
}

func TestAddLicense(t *testing.T) {
	tmpl := template.Must(template.New("").Parse("{{.Holder}}{{.Year}}{{.SPDXID}}"))
	data := licenseData{Holder: "H", Year: "Y", SPDXID: "S"}

	tests := []struct {
		contents     string
		wantContents string
		wantUpdated  bool
	}{
		{"", "// HYS\n\n", true},
		{"content", "// HYS\n\ncontent", true},

		// various headers that should be left intact. Many don't make
		// sense for our temp file extension, but that doesn't matter.
		{"#!/bin/bash\ncontent", "#!/bin/bash\n// HYS\n\ncontent", true},
		{"<?xml version='1.0'?>\ncontent", "<?xml version='1.0'?>\n// HYS\n\ncontent", true},
		{"<!doctype html>\ncontent", "<!doctype html>\n// HYS\n\ncontent", true},
		{"<!DOCTYPE HTML>\ncontent", "<!DOCTYPE HTML>\n// HYS\n\ncontent", true},
		{"# encoding: UTF-8\ncontent", "# encoding: UTF-8\n// HYS\n\ncontent", true},
		{"# frozen_string_literal: true\ncontent", "# frozen_string_literal: true\n// HYS\n\ncontent", true},
		{"<?php\ncontent", "<?php\n// HYS\n\ncontent", true},
		{"# escape: `\ncontent", "# escape: `\n// HYS\n\ncontent", true},
		{"# syntax: docker/dockerfile:1.3\ncontent", "# syntax: docker/dockerfile:1.3\n// HYS\n\ncontent", true},

		// ensure files with existing license or generated files are
		// skipped. No need to test all permutations of these, since
		// there are specific tests below.
		{"// Copyright 2000 Acme\ncontent", "// Copyright 2000 Acme\ncontent", false},
		{"// Code generated by go generate; DO NOT EDIT.\ncontent", "// Code generated by go generate; DO NOT EDIT.\ncontent", false},
	}

	for _, tt := range tests {
		// create temp file with contents
		f, err := createTempFile(tt.contents, "*.go")
		if err != nil {
			t.Error(err)
		}
		fi, err := f.Stat()
		if err != nil {
			t.Error(err)
		}

		// run addlicense
		updated, err := addLicense(f.Name(), fi.Mode(), tmpl, data)
		if err != nil {
			t.Error(err)
		}

		// check results
		if updated != tt.wantUpdated {
			t.Errorf("addLicense with contents %q returned updated: %t, want %t", tt.contents, updated, tt.wantUpdated)
		}
		gotContents, err := ioutil.ReadFile(f.Name())
		if err != nil {
			t.Error(err)
		}
		if got := string(gotContents); got != tt.wantContents {
			t.Errorf("addLicense with contents %q returned contents: %q, want %q", tt.contents, got, tt.wantContents)
		}

		// if all tests passed, cleanup temp file
		if !t.Failed() {
			_ = os.Remove(f.Name())
		}
	}
}

// Test that license headers are added using the appropriate prefix for
// different filenames and extensions.
func TestLicenseHeader(t *testing.T) {
	tpl := template.Must(template.New("").Parse("{{.Holder}}{{.Year}}{{.SPDXID}}"))
	data := licenseData{Holder: "H", Year: "Y", SPDXID: "S"}

	tests := []struct {
		paths []string // paths passed to licenseHeader
		want  string   // expected result of executing template
	}{
		{
			[]string{"f.unknown"},
			"",
		},
		{
			[]string{"f.c", "f.h", "f.gv", "f.java", "f.scala", "f.kt", "f.kts"},
			"/*\n * HYS\n */\n\n",
		},
		{
			[]string{"f.js", "f.mjs", "f.cjs", "f.jsx", "f.tsx", "f.css", "f.scss", "f.sass", "f.ts"},
			"/**\n * HYS\n */\n\n",
		},
		{
			[]string{"f.cc", "f.cpp", "f.cs", "f.go", "f.hcl", "f.hh", "f.hpp", "f.m", "f.mm", "f.proto",
				"f.rs", "f.swift", "f.dart", "f.groovy", "f.v", "f.sv", "f.php"},
			"// HYS\n\n",
		},
		{
			[]string{"f.py", "f.sh", "f.yaml", "f.yml", "f.dockerfile", "dockerfile", "f.rb", "gemfile", "f.tcl", "f.tf", "f.bzl", "f.pl", "f.pp", "build"},
			"# HYS\n\n",
		},
		{
			[]string{"f.el", "f.lisp"},
			";; HYS\n\n",
		},
		{
			[]string{"f.erl"},
			"% HYS\n\n",
		},
		{
			[]string{"f.hs", "f.sql", "f.sdl"},
			"-- HYS\n\n",
		},
		{
			[]string{"f.html", "f.xml", "f.vue", "f.wxi", "f.wxl", "f.wxs"},
			"<!--\n HYS\n-->\n\n",
		},
		{
			[]string{"f.ml", "f.mli", "f.mll", "f.mly"},
			"(**\n   HYS\n*)\n\n",
		},
		{
			[]string{"cmakelists.txt", "f.cmake", "f.cmake.in"},
			"# HYS\n\n",
		},

		// ensure matches are case insenstive
		{
			[]string{"F.PY", "DoCkErFiLe"},
			"# HYS\n\n",
		},
	}

	for _, tt := range tests {
		for _, path := range tt.paths {
			header, _ := licenseHeader(path, tpl, data)
			if got := string(header); got != tt.want {
				t.Errorf("licenseHeader(%q) returned: %q, want: %q", path, got, tt.want)
			}
		}
	}
}

// Test that generated files are properly recognized.
func TestIsGenerated(t *testing.T) {
	tests := []struct {
		content string
		want    bool
	}{
		{"", false},
		{"Generated", false},
		{"// Code generated by go generate; DO NOT EDIT.", true},
		{"/*\n* Code generated by go generate; DO NOT EDIT.\n*/\n", true},
		{"DO NOT EDIT! Replaced on runs of cargo-raze", true},
	}

	for _, tt := range tests {
		b := []byte(tt.content)
		if got := isGenerated(b); got != tt.want {
			t.Errorf("isGenerated(%q) returned %v, want %v", tt.content, got, tt.want)
		}
	}
}

// Test that existing license headers are identified.
func TestHasLicense(t *testing.T) {
	tests := []struct {
		content string
		want    bool
	}{
		{"", false},
		{"This is my license", false},
		{"This code is released into the public domain.", false},
		{"SPDX: MIT", false},

		{"Copyright 2000", true},
		{"CoPyRiGhT 2000", true},
		{"Subject to the terms of the Mozilla Public License", true},
		{"SPDX-License-Identifier: MIT", true},
		{"spdx-license-identifier: MIT", true},
	}

	for _, tt := range tests {
		b := []byte(tt.content)
		if got := hasLicense(b); got != tt.want {
			t.Errorf("hasLicense(%q) returned %v, want %v", tt.content, got, tt.want)
		}
	}
}

func TestFileMatches(t *testing.T) {
	tests := []struct {
		pattern   string
		path      string
		wantMatch bool
	}{
		// basic single directory patterns
		{"", "file.c", false},
		{"*.c", "file.h", false},
		{"*.c", "file.c", true},

		// subdirectory patterns
		{"*.c", "vendor/file.c", false},
		{"**/*.c", "vendor/file.c", true},
		{"vendor/**", "vendor/file.c", true},
		{"vendor/**/*.c", "vendor/file.c", true},
		{"vendor/**/*.c", "vendor/a/b/file.c", true},

		// single character "?" match
		{"*.?", "file.c", true},
		{"*.?", "file.go", false},
		{"*.??", "file.c", false},
		{"*.??", "file.go", true},

		// character classes - sets and ranges
		{"*.[ch]", "file.c", true},
		{"*.[ch]", "file.h", true},
		{"*.[ch]", "file.ch", false},
		{"*.[a-z]", "file.c", true},
		{"*.[a-z]", "file.h", true},
		{"*.[a-z]", "file.go", false},
		{"*.[a-z]", "file.R", false},

		// character classes - negations
		{"*.[^ch]", "file.c", false},
		{"*.[^ch]", "file.h", false},
		{"*.[^ch]", "file.R", true},
		{"*.[!ch]", "file.c", false},
		{"*.[!ch]", "file.h", false},
		{"*.[!ch]", "file.R", true},

		// comma-separated alternative matches
		{"*.{c,go}", "file.c", true},
		{"*.{c,go}", "file.go", true},
		{"*.{c,go}", "file.h", false},

		// negating alternative matches
		{"*.[^{c,go}]", "file.c", false},
		{"*.[^{c,go}]", "file.go", false},
		{"*.[^{c,go}]", "file.h", true},
	}

	for _, tt := range tests {
		patterns := []string{tt.pattern}
		if got := fileMatches(tt.path, patterns); got != tt.wantMatch {
			t.Errorf("fileMatches(%q, %q) returned %v, want %v", tt.path, patterns, got, tt.wantMatch)
		}
	}
}
