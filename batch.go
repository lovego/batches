package batches

import (
	"io"
	"math"
	"sync/atomic"
	"time"
)

type Batches struct {
	From        int64
	To          int64
	BatchSize   uint32 // BatchSize defaults to 1 if is 0.
	Concurrency uint16
	Work        func(int64, int64) error
	Output      io.Writer
}

func (b Batches) Run() (err error) {
	itr := b.iterator()

	pgs := b.progress(itr.isSingleStep())
	pgs.begin(b.From, b.To, itr.totalCount)

	if b.Concurrency <= 1 {
		err = b.runSerially(itr, pgs)
	} else {
		err = b.runConcurrently(itr, pgs)
	}

	if err == nil {
		pgs.end(b.From, b.To, itr.totalCount)
	}
	return
}

func (b Batches) runSerially(itr iterator, pgs *progress) error {
	for {
		from, to, ok := itr.next()
		if !ok {
			return nil
		}
		startTime := time.Now()
		pgs.start(startTime, from, to)
		if err := b.Work(from, to); err != nil {
			return err
		}
		pgs.finish(startTime)
	}
}

func (b Batches) runConcurrently(itr iterator, pgs *progress) (err error) {
	var dataChan = make(chan [2]int64)
	var resultChan = make(chan error, 1)

	b.startWorkers(itr.totalCount, dataChan, resultChan, pgs)
	for {
		from, to, ok := itr.next()
		if !ok {
			goto stop
		}

		// because select is random, check err first.
		select {
		case err = <-resultChan:
			if err != nil {
				goto stop
			}
		default:
		}

		select {
		case err = <-resultChan:
			if err != nil {
				goto stop
			}
		case dataChan <- [2]int64{from, to}:
		}
	}

stop:
	close(dataChan)
	if err != nil {
		return err
	}
	return <-resultChan
}

func (b Batches) startWorkers(
	batchCount uint64, dataChan chan [2]int64, resultChan chan error, pgs *progress,
) {
	concurrency := b.Concurrency
	if batchCount < math.MaxUint16 && concurrency > uint16(batchCount) {
		concurrency = uint16(batchCount)
	}
	var running = int32(concurrency)
	for i := uint16(0); i < concurrency; i++ {
		go func() {
			for {
				data, ok := <-dataChan
				if !ok {
					if atomic.AddInt32(&running, -1) == 0 {
						resultChan <- nil
					}
					return
				}
				from, to := data[0], data[1]

				var startTime = time.Now()
				pgs.startLine(startTime, from, to)
				if err := b.Work(from, to); err != nil {
					resultChan <- err
					return
				}
				pgs.finishLine(startTime, from, to)
			}
		}()
	}
}
