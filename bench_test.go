// Copyright 2016 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

// Only build with go1.7 or above because b.Run did not exist before.
// +build go1.7

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/mtail/mtail"
	"github.com/google/mtail/watcher"
	"github.com/spf13/afero"
)

var (
	recordBenchmark = flag.Bool("record_benchmark", false, "Record the benchmark results to 'benchmark_results.csv'.")
)

func BenchmarkProgram(b *testing.B) {
	// exampleProgramTests live in ex_test.go
	for _, bm := range exampleProgramTests {
		b.Run(fmt.Sprintf("%s on %s", bm.programfile, bm.logfile), func(b *testing.B) {
			b.ReportAllocs()
			w := watcher.NewFakeWatcher()
			fs := afero.NewOsFs()
			log, err := fs.Create("/tmp/test")
			if err != nil {
				b.Fatalf("failed to create test file descriptor")
			}
			logs = []string{log.Name()}
			o := mtail.Options{Progs: bm.programfile, LogPathPatterns: logs, W: w, FS: fs}
			mtail, err := mtail.New(o)
			if err != nil {
				b.Fatalf("Failed to create mtail: %s", err)
			}
			mtail.StartTailing()

			var total int64
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				l, err := os.Open(bm.logfile)
				if err != nil {
					b.Fatalf("Couldn't open logfile: %s", err)
				}
				count, err := io.Copy(log, l)

				if err != nil {
					b.Fatalf("Write of test data failed to test file: %s", err)
				}
				total += count
			}
			mtail.Close()
			b.StopTimer()
			b.SetBytes(total)
		})
	}
}
