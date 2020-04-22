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
