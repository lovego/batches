package batches

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"testing"
	"time"
)

func ExampleBatches_Run() {
	err := Batches{
		From:      1,
		To:        307,
		BatchSize: 100,
		Work: func(from, to int64) error {
			fmt.Println(from, to)
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 1 100
	// 101 200
	// 201 300
	// 301 307
}

func ExampleBatches_Run_desc() {
	err := Batches{
		From:      307,
		To:        7,
		BatchSize: 100,
		Work: func(from, to int64) error {
			fmt.Println(from, to)
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 307 208
	// 207 108
	// 107 8
	// 7 7
}

func ExampleBatches_Run_withError() {
	err := Batches{
		From:      1,
		To:        307,
		BatchSize: 100,
		Work: func(from, to int64) error {
			fmt.Println(from, to)
			if from == 201 {
				return errors.New("error")
			}
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 1 100
	// 101 200
	// 201 300
	// error
}

func ExampleBatches_Run_concurrently() {
	var result = [][2]int64{}
	var lock sync.Mutex
	err := Batches{
		From:        307,
		To:          7,
		BatchSize:   100,
		Concurrency: 5,
		Work: func(from, to int64) error {
			time.Sleep(time.Millisecond)
			lock.Lock()
			defer lock.Unlock()
			result = append(result, [2]int64{from, to})
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i][0] > result[j][0]
	})
	fmt.Println(result)

	// Output:
	// [[307 208] [207 108] [107 8] [7 7]]
}

func ExampleBatches_Run_concurrentlyWithError() {
	err := Batches{
		From:        1,
		To:          7,
		Concurrency: 2,
		Work: func(from, to int64) error {
			if from == 4 {
				return errors.New("error")
			}
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Output: error
}

func ExampleBatches_Run_concurrentlyWithError2() {
	err := Batches{
		From:        1,
		To:          7,
		Concurrency: 2,
		Work: func(from, to int64) error {
			if from == 4 {
				return errors.New("error")
			}
			time.Sleep(10 * time.Millisecond)
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Output: error
}

func ExampleBatches_Run_singleStep() {
	err := Batches{
		From: 3,
		To:   1,
		Work: func(from, to int64) error {
			fmt.Println(from, to)
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 3 3
	// 2 2
	// 1 1
}

func ExampleBatches_Run_singleBatch() {
	err := Batches{
		From:      3,
		To:        1,
		BatchSize: 100,
		Work: func(from, to int64) error {
			time.Sleep(time.Second)
			fmt.Println(from, to)
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 3 1
}

func ExampleBatches_Run_oneNumber() {
	err := Batches{
		From: 3,
		To:   3,
		Work: func(from, to int64) error {
			time.Sleep(time.Millisecond)
			fmt.Println(from, to)
			return nil
		},
		Output: os.Stderr,
	}.Run()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 3 3
}

func ExampleBatches_Run_oneNumber2() {
	err := Batches{
		From:      3,
		To:        3,
		BatchSize: 3,
		Work: func(from, to int64) error {
			fmt.Println(from, to)
			return nil
		},
	}.Run()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 3 3
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func ExamplePrettyDuration() {
	fmt.Printf("%#v\n", PrettyDuration(67*time.Minute))
	fmt.Printf("%#v\n", PrettyDuration(67*time.Second))

	// Output:
	// " 1h7m0s"
	// "   1m7s"
}

func ExampleProgress() {
	var pgs *progress
	var t time.Time
	pgs.startLine(t, 0, 0)
	pgs.finishLine(t, 0, 0)
	pgs.print(t, "")

	// Output:
}
