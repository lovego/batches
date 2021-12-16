// +build example

package batches_test

import (
	"fmt"
	"os"

	"github.com/lovego/batches"
)

// This is just for doc in README.md, it will not pass.
// Run with command:
//   go test -run _example -tags example
func ExampleBatches_example() {
	err := batches.Batches{
		From:        1,
		To:          307,
		BatchSize:   100,
		Concurrency: 1,
		Work: func(from, to int64) error {
			fmt.Printf("\t*do work for [%3d ~ %3d]*\t", from, to)
			return nil
		},
		Output: os.Stdout,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Output:
}
