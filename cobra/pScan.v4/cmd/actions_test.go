package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"cobra/pScan.v4/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Create temp file
	tf, err := os.CreateTemp("", "pScan")

	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	// Initialize list id needed
	if initList {
		hl := &scan.HostsList{}

		for _, h := range hosts {
			hl.Add(h)
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	// Return temp file name and cleanup function
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostAction(t *testing.T) {
	// Define hosts for actions test
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// Test cases for Action test
	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},

		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},

		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Deleted host: host1\nDeleted host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//setup Action test
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			// Define var to capture Action output
			var out bytes.Buffer

			// Execute Action and capture output
			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("expected no error, got %q\n", err)
			}

			// Test Action output
			if out.String() != tc.expectedOut {
				t.Errorf("expected output %q, got %q\n", tc.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	// Define hosts for integration test
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// Setup integration test
	tf, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"

	hostsEnd := []string{
		"host1",
		"host2",
	}

	// Define var to capture output
	var out bytes.Buffer

	// Define expected output of all actions
	expectedOut := ""

	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", v)
	}

	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	// Add hosts to the list
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// list host
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// Delet  host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	//List hosts after delete
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// Test integration output
	if out.String() != expectedOut {
		t.Errorf("expected output %q, got %q\n", expectedOut, out.String())
	}
}
