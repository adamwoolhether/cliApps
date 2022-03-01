package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestListAction(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		expOut   string
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{name: "Results", expError: nil, expOut: "-  1  Task 1\n-  2  Task 2\n", resp: testResp["resultsMany"]},
		{name: "NoResults", expError: ErrNotFound, resp: testResp["noResults"]},
		{name: "InvalidURL", expError: ErrConnection, resp: testResp["noResults"], closeServer: true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})
			defer cleanup()
			
			if tc.closeServer {
				cleanup()
			}
			
			var out bytes.Buffer
			err := listAction(&out, url)
			
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("exp error %q, got no error", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("exp error %q, got %q", tc.expError, err)
				}
				
				return
			}
			if err != nil {
				t.Fatalf("exp no error, got %q", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("exp output %q, got %q", tc.expOut, out.String())
			}
		})
	}
}

func TestViewAcount(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		expOut   string
		resp     struct {
			Status int
			Body   string
		}
		id string
	}{
		{name: "ResultsOne",
			expError: nil,
			expOut: `Task:         Task 1
Created at:   Oct/28 @08:23
Completed:    No
`,
			resp: testResp["resultsOne"],
			id:   "1"},
		{name: "NotFound", expError: ErrNotFound, resp: testResp["notFound"], id: "1"},
		{name: "InvalidID", expError: ErrNotNumber, resp: testResp["noResults"], id: "a"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintln(w, tc.resp.Body)
				})
			defer cleanup()
			
			var out bytes.Buffer
			
			err := viewAction(&out, url, tc.id)
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Exp error %q, got no error", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("Exp error %q, got %q", tc.expError, err)
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Exp no error, got %q", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Exp output %q, got %q", tc.expOut, out.String())
			}
		})
	}
}
