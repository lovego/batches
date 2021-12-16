# batches
Run a sequence of data by batch, with optional concurrency, and optional progress printing. 

[![Build Status](https://github.com/lovego/batches/actions/workflows/go.yml/badge.svg)](https://github.com/lovego/batches/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/lovego/batches/badge.svg?branch=master)](https://coveralls.io/github/lovego/batches)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/batches)](https://goreportcard.com/report/github.com/lovego/batches)
[![Documentation](https://pkg.go.dev/badge/github.com/lovego/batches)](https://pkg.go.dev/github.com/lovego/batches@v0.0.1)

## Install
`$ go get github.com/lovego/batches`

## Example
```go
func ExampleBatches_example() {
	err := batches.Batches{
		From:      1,
		To:        307,
		BatchSize: 100,
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
	// 17:49:07 [1 ~ 307] total 4 batches
	// 17:49:07 [  1 ~ 100] ... 	*do work for [  1 ~ 100]*	 finished in  31.6µs
	// 17:49:07 [101 ~ 200] ... 	*do work for [101 ~ 200]*	 finished in   8.9µs
	// 17:49:07 [201 ~ 300] ... 	*do work for [201 ~ 300]*	 finished in   7.5µs
	// 17:49:07 [301 ~ 307] ... 	*do work for [301 ~ 307]*	 finished in   7.3µs
	// 17:49:07 [1 ~ 307] total 4 batches, finished in 518.5µs
}
```
