package data_processor

import (
	"testing"
)

// Two functions require fuzzying tests, which are parser and DataProcessor. parser is a part of DataProcessor, hence, DP needs to be fuzzy tested.
func FuzzDataProcessor(f *testing.F) {
	tcSeed := []string{"test-1\n+\n1\n1","test-2\n-\n2\n1","test-3\n*\n3\n1","test-4\n/\n4\n2","invalid-test", ""}
	for _, tc := range tcSeed {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, inp string) {
		inpChan := make(chan []byte, 1)
		inpChan<-[]byte(inp)
		close(inpChan)
		outChan := make(chan Result, 1)
		DataProcessor(inpChan, outChan)
		_ = <-outChan
		close(outChan)
	})
}
// The task itself asks to test parser function, so here it is.
func FuzzParser(f *testing.F) {
	tcSeed := []string{"test-1\n+\n1\n1","test-2\n-\n2\n1","test-3\n*\n3\n1","test-4\n/\n4\n2","invalid-test", ""}
	for _, tc := range tcSeed {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, inp string) {
		inpVal := []byte(inp)
		_, _ = parser(inpVal)
	}
}
