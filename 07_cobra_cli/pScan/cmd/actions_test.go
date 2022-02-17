package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	
	"github.com/adamwoolhether/cliApps/07_cobra_cli/pScan/scan"
)

// setup prepares the necessary hostsFile to be used for testing.
func setup(t *testing.T, hostsList []string, initList bool) (string, func()) {
	// Create temp file
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()
	
	if initList {
		hl := &scan.HostsList{}
		
		for _, h := range hostsList {
			hl.Add(h)
		}
		if err = hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}
	
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	// Define the hosts to be used for the test.
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}
	
	// Test cases:
	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOutput: "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			args:           hosts,
			expectedOutput: "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeletedAction",
			args:           []string{"host1", "host2"},
			expectedOutput: "Deleted host: host1\nDeleted host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Action test
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()
			
			// Define var to capture the Action output
			var out bytes.Buffer
			
			// Execute action and capture output.
			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("Exp no error, got %q\n", err)
			}
			
			// Test Actions output
			if out.String() != tc.expectedOutput {
				t.Errorf("Exp output %q, got %q\n", tc.expectedOutput, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	
	tf, cleanup := setup(t, hosts, false)
	defer cleanup()
	
	// We'll delete this host
	delHost := "host2"
	hostsEnd := []string{"host1", "host3"}
	
	var out bytes.Buffer
	
	// Define expected output for all actions
	expectedOutput := ""
	for _, v := range hosts {
		expectedOutput += fmt.Sprintf("Added host: %s\n", v)
	}
	expectedOutput += strings.Join(hosts, "\n")
	expectedOutput += fmt.Sprintln()
	expectedOutput += fmt.Sprintf("Deleted host: %s\n", delHost)
	expectedOutput += strings.Join(hostsEnd, "\n")
	expectedOutput += fmt.Sprintln()
	
	// Run add -> list -> delete -> list sequence
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Exp no error, got %q\n", err)
	}
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Exp no error, got %q\n", err)
	}
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Exp no error, got %q\n", err)
	}
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Exp no error, got %q\n", err)
	}
	
	// Test integration output
	if out.String() != expectedOutput {
		t.Errorf("Exp output %q, got %q\n", expectedOutput, out.String())
	}
}
