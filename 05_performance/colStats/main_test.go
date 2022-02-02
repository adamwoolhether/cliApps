package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
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
		{name: "RunAvg1File", col: 3, op: "avg", exp: "227.6\n", files: []string{"./testdata/example.csv"}, expErr: nil},
		{name: "RunAvgMultiFiles", col: 3, op: "avg", exp: "233.84\n", files: []string{"./testdata/example.csv", "./testdata/example2.csv"}, expErr: nil},
		{name: "RunFailRead", col: 2, op: "avg", exp: "", files: []string{"./testdata/example.csv", "./testdata/fakefile.csv"}, expErr: os.ErrNotExist},
		{name: "RunFailColumn", col: 0, op: "avg", exp: "", files: []string{"./testdata/example.csv"}, expErr: ErrInvalidColumn},
		{name: "RunFailNoFiles", col: 2, op: "avg", exp: "", files: []string{}, expErr: ErrNoFiles},
		{name: "RunFailOperation", col: 2, op: "invalid", exp: "", files: []string{"./testdata/example.csv"}, expErr: ErrInvalidOperation},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bytes.Buffer
			err := run(tc.files, tc.op, tc.col, &res)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error, go nil")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %q, got %q", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %q", err)
			}
			if res.String() != tc.exp {
				t.Errorf("expected %q, got %q", tc.exp, &res)
			}
		})
	}
}

/*
// To run benchmarks:
gotest -bench . -run ^$
gotest -bench . -benchtime=10x -run ^$
gotest -bench . -benchtime=10x -run ^$ | tee benchresults00.txt
gotest -bench . -benchtime=10x -run ^$ -benchmem | tee benchresults00m.txt

// For profiling:
gotest -bench . -benchtime=10x -run ^$ -cpuprofile cpu00.pprof
go tool pprof cpu00.pprof
top
top -cum
list csv2Float
web

// Mem profiling:
gotest -bench . -benchtime=10x -run ^$ -memprofile mem00.pprof
go tool pprof --alloc_space mem00.pprof

// To compare results, use benchcmp:
go get -u -v golang.org/x/tools/cmd/benchcmp
benchcmp benchresults00m.txt benchresults01m.txt
*/
func BenchmarkRun(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err = run(filenames, "avg", 2, ioutil.Discard); err != nil {
			b.Error(err)
		}
	}
}
