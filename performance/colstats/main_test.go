package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{name: "RunAvgFile", col: 3, op: "avg", exp: "227.6\n", files: []string{"./testdata/example.csv"}, expErr: nil},
		{name: "RunAvgFile", col: 3, op: "avg", exp: "223.86\n", files: []string{"./testdata/example2.csv", "./testdata/example.csv"}, expErr: nil},
		{name: "RunFailRead", col: 2, op: "avg", exp: "", files: []string{"./testdata/example.csv", "./testdata/fakefile.csv"}, expErr: os.ErrNotExist},
		{name: "RunFailColumn", col: 0, op: "avg", exp: "", files: []string{"./testdata/example2.csv"}, expErr: ErrInvalidColumn},
		{name: "RunFailNOFiles", col: 2, op: "avg", exp: "", files: []string{}, expErr: ErrNoFiles},
		{name: "RunFailOperation", col: 2, op: "invalid", exp: "", files: []string{"./testdata/example.csv"}, expErr: ErrInvalidOperation},
	}

	// Run tests execution
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bytes.Buffer
			err := run(tc.files, tc.op, tc.col, &res)

			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error. Got nil insted")
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %q, got %q instead", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %q", err)
			}

			if res.String() != tc.exp {
				t.Errorf("expected %q, got %q instead", tc.exp, &res)
			}
		})
	}
}
