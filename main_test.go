package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/VirtusLab/render/files"

	"github.com/VirtusLab/render/constants"
)

const (
	testBinaryName = "testrender"
	killIn         = 30 * time.Second
)

var (
	exeSuffix string // ".exe" on Windows
)

func init() {
	switch runtime.GOOS {
	case "windows":
		exeSuffix = ".exe"
	}
}

// The TestMain function creates a the binary for testing purposes
// and deletes it after the tests have been run.
func TestMain(m *testing.M) {
	// build the test binary
	args := []string{"build", "-o", testBinaryName + exeSuffix}
	out, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "building %s failed: %v\n%s", testBinaryName, err, out)
		os.Exit(2)
	}
	// remove test binary
	defer os.Remove(testBinaryName + exeSuffix)

	flag.Parse()
	merr := m.Run()
	if merr != 0 {
		fmt.Printf("Main tests failed.\n")
		os.Exit(merr)
	}

	os.Exit(0)
}

func run(args ...string) (stdout, stderr string, err error) {
	prog := "./" + testBinaryName + exeSuffix
	// always add debug flag
	newargs := append([]string{"-d"}, args...)
	ctx, _ := context.WithTimeout(context.TODO(), killIn)

	fmt.Printf("$ %s %s\n\n", prog, strings.Join(newargs, " "))
	stdout, stderr, err = sh(ctx, prog, newargs...)
	fmt.Printf("stdout:\n%s\n\n", stdout)
	fmt.Printf("stderr:\n%s\n\n", stderr)

	return
}

func sh(ctx context.Context, prog string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, prog, args...)

	// Set output to Byte Buffers
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	stdout = outb.String()
	stderr = errb.String()

	return
}

func TestHelp(t *testing.T) {
	stdout, _, err := run("-h")
	if err != nil {
		t.Fatalf("output: '%s', error: %v", string(stdout), err)
	}

	expected := fmt.Sprintf("%s - %s", constants.Name, constants.Description)
	if !strings.Contains(stdout, expected) {
		t.Fatalf("expected contains:\n%s\ngot:\n%s", expected, stdout)
	}
}

func TestRender(t *testing.T) {
	stdout, _, err := run("--config", "examples/example.config.yaml", "--in", "examples/example.yaml.tmpl")
	if err != nil {
		t.Fatalf("output: '%s', error: %v", string(stdout), err)
	}

	expectedPath := "examples/example.yaml.expected"
	expected, err := files.ReadInput(expectedPath)
	if err != nil {
		t.Fatalf("cannot read test file: '%s'", expectedPath)
	}

	if stdout != string(expected) {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, stdout)
	}
}

func TestNoArgs(t *testing.T) {
	stdout, stderr, err := run()
	if ee, ok := err.(*exec.ExitError); ok {
		if ee.String() != "exit status 1" {
			t.Fatal("expected exit status 1")
		}
	} else if err != nil {
		t.Fatalf("output: '%s', error: %v", string(stdout), err)
	}

	expectedStdout := ``
	if stdout != expectedStdout {
		t.Fatalf("expected:\n%s\ngot:\n%s", expectedStdout, stdout)
	}

	expectedStderr := `expected either stdin or --in parameter, for usage use --help`
	if !strings.Contains(stderr, expectedStderr) {
		t.Fatalf("expected contains:\n%s\ngot:\n%s", expectedStderr, stderr)
	}
}
