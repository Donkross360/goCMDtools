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
		closeSever bool
	}{
		{name: "Results", expError: nil, expOut: "- 1 Task 1\n- 2 Task 2\n", resp: testResp["resultsMany"]},
		{name: "NoResults", expError: ErrNotFound, resp: testResp["noResults"]},
		{name: "InvalidURL", expError: ErrConnection, resp: testResp["noResults"], closeSever: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})

			defer cleanup()
			if tc.closeSever {
				cleanup()
			}

			var out bytes.Buffer
			err := listAction(&out, url)

			if tc.expError != nil {
				if err == nil {
					t.Fatalf("expected error %q, got no error.", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("expected error %q, got %q.", tc.expError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %q.", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("expected out %q, got %q", tc.expOut, out.String())
			}

		})
	}
}

func TestViewAction(t *testing.T) {
	// testCases for ViewAction test
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
		{name: "ResultsOne", expError: nil, expOut: `Task: 	Task1    Created at: Feb/01 @10:44    Completed: No`, resp: testResp["resultsOne"], id: "1"},
		{name: "NotFound", expError: ErrNotFound, resp: testResp["notFound"], id: "1"},
		{name: "InvalidID", expError: ErrNotNumber, resp: testResp["noResults"], id: "a"},
	}

	// Execute ViewAction test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				fmt.Fprintln(w, tc.resp.Body)
			})
			defer cleanup()

			var out bytes.Buffer

			err := viewAction(&out, url, tc.id)

			if tc.expError != nil {
				if err == nil {
					t.Fatalf("expected error %q, got no error.", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("expected error %q, got %q.", tc.expError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %q.", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("expected output %q, got %q", tc.expOut, out.String())
			}

		})
	}
}
