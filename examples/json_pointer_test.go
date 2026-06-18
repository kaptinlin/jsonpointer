package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/kaptinlin/jsonpointer"
)

var errWriteFailed = errors.New("write failed")

const wantOutput = "name: Alice\nescaped value: ready\nescaped key: tilde~key\nemail: alice@example.com\ntokens: [users 0 name]\npointer: /users/0/name\n"

func TestDefaultExampleData(t *testing.T) {
	t.Parallel()

	data := defaultExampleData()
	if data.namePath != "/users/0/name" {
		t.Fatalf("defaultExampleData() namePath = %q", data.namePath)
	}
	if data.escapePath != "/foo~1bar/tilde~0key" {
		t.Fatalf("defaultExampleData() escapePath = %q", data.escapePath)
	}
}

func TestWriteExample(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	err := writeExample(&out, defaultExampleData())
	if err != nil {
		t.Fatal(err)
	}
	if out.String() != wantOutput {
		t.Errorf("writeExample() output = %q, want %q", out.String(), wantOutput)
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	err := run(&out)
	if err != nil {
		t.Fatal(err)
	}
	if out.String() != wantOutput {
		t.Errorf("run() output = %q, want %q", out.String(), wantOutput)
	}
}

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) {
	return 0, errWriteFailed
}

func TestWriteExampleReturnsNamePointerError(t *testing.T) {
	t.Parallel()

	data := defaultExampleData()
	data.namePath = "/users/9/name"

	err := writeExample(io.Discard, data)
	if !errors.Is(err, jsonpointer.ErrIndexOutOfBounds) {
		t.Fatalf("writeExample() error = %v, want %v", err, jsonpointer.ErrIndexOutOfBounds)
	}
}

func TestWriteExampleReturnsReferenceError(t *testing.T) {
	t.Parallel()

	data := defaultExampleData()
	data.escapePath = "/missing"

	err := writeExample(io.Discard, data)
	if !errors.Is(err, jsonpointer.ErrKeyNotFound) {
		t.Fatalf("writeExample() error = %v, want %v", err, jsonpointer.ErrKeyNotFound)
	}
}

func TestWriteExampleReturnsWriterError(t *testing.T) {
	t.Parallel()

	err := writeExample(failWriter{}, defaultExampleData())
	if !errors.Is(err, errWriteFailed) {
		t.Fatalf("writeExample() error = %v, want %v", err, errWriteFailed)
	}
}

func TestMain(t *testing.T) {
	// main writes to process-wide stdout.
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	originalStdout := os.Stdout
	os.Stdout = writePipe
	defer func() {
		os.Stdout = originalStdout
	}()

	main()

	if err := writePipe.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(readPipe)
	if err != nil {
		t.Fatal(err)
	}
	if err := readPipe.Close(); err != nil {
		t.Fatal(err)
	}
	if string(output) != wantOutput {
		t.Errorf("main() output = %q, want %q", output, wantOutput)
	}
}
