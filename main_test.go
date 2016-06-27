// Copyright 2016 Google Inc.
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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
		cargs := []string{"-l", "apache", "-c", "Google Inc.", "-y", "2016", tmp}
		c := exec.Command(os.Args[0], append(targs, cargs...)...)
		c.Env = []string{"RUNME=1"}
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}

		run(t, "diff", "-r", filepath.Join(tmp, "initial"), "testdata/expected")
	}
}
