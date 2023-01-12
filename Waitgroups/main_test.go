package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_main(t *testing.T) {
	output := os.Stdout

	var r, w, _ = os.Pipe()

	msg = "message"

	os.Stdout = w
	main()

	_ = w.Close()

	printed, _ := io.ReadAll(r)
	outputted := string(printed)

	os.Stdout = output

	msgs := []string{
		"Hello, universe!",
		"Hello, cosmos!",
		"Hello, world!",
	}

	for _, m := range msgs {
		if !strings.Contains(outputted, m) {
			t.Errorf("Printed content does not contain content: " + m)
		}
	}
}

func Test_updateMessage(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	updateMessage("test", &wg)
	wg.Wait()

	result := strings.Compare(msg, "test")
	if result != 0 {
		t.Errorf("Expected string not found in msg")
	}
}

func Test_printMessage(t *testing.T) {
	output := os.Stdout

	var r, w, _ = os.Pipe()

	msg = "message"

	os.Stdout = w
	printMessage()

	_ = w.Close()

	printed, _ := io.ReadAll(r)

	os.Stdout = output

	if !strings.Contains(string(printed), "message") {
		t.Errorf("Printed message does not contain test content")
	}
}
