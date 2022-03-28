package scan_test

import (
	"net"
	"strconv"
	"testing"
	
	"github.com/adamwoolhether/cliApps/cobraCli/pScan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}
	
	if ps.Open.String() != "closed" {
		t.Errorf("Exp %q, got %q\n", "closed", ps.Open.String())
	}
	
	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("exp %q, got %q\n", "open", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name          string
		expectedState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}
	
	host := "localhost"
	hl := &scan.HostsList{}
	
	hl.Add(host)
	
	ports := []int{}
	
	// Init 1 open and 1 closed port.
	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		
		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}
		
		ports = append(ports, port)
		
		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}
	
	res := scan.Run(hl, ports)
	
	// Verify results for HostFound test.
	if len(res) != 1 {
		t.Fatalf("exp 1 reult, got %d\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Exp host %q, got %q\n", host, res[0].Host)
	}
	if res[0].NotFound {
		t.Errorf("Exp host %q to be found", host)
	}
	if len(res[0].PortStates) != 2 {
		t.Fatalf("Exp 2 ports, go %d\n", len(res[0].PortStates))
	}
	
	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Exp port %d, got %d\n", ports[0], res[0].PortStates[i].Port)
		}
		if res[0].PortStates[i].Open.String() != tc.expectedState {
			t.Errorf("Exp port %d to be %s\n", ports[i], tc.expectedState)
		}
	}
}

func TestRunHostNotFound(t *testing.T) {
	host := "389.389.389.389"
	hl := &scan.HostsList{}
	
	hl.Add(host)
	
	res := scan.Run(hl, []int{})
	
	// Verify results for HostNotFound test
	if len(res) != 1 {
		t.Fatalf("Exp 1 result, got %d\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Exp host %q, got %q", host, res[0].Host)
	}
	if !res[0].NotFound {
		t.Errorf("Exp host %q NOT to be found\n", host)
	}
	if len(res[0].PortStates) != 0 {
		t.Errorf("Exp 0 ports states, got %d instead\n", len(res[0].PortStates))
	}
}
