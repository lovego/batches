package batches

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

func (b *Batches) progress(single bool) *progress {
	if b.Output == nil {
		return nil
	}
	p := &progress{
		output: b.Output,
		single: single,
	}

	width := len(strconv.FormatInt(b.From, 10))
	if n := len(strconv.FormatInt(b.To, 10)); n > width {
		width = n
	}
	w := strconv.Itoa(width)
	if p.single {
		p.batch = "[%" + w + "d]"
	} else {
		p.batch = "[%" + w + "d ~ %" + w + "d]"
	}
	return p
}

type progress struct {
	output    io.Writer
	single    bool
	batch     string
	beginTime time.Time
}

func (p *progress) begin(from, to int64, batchCount uint64) {
	if p == nil {
		return
	}
	p.beginTime = time.Now()
	p.print(p.beginTime, fmt.Sprintf(
		"[%d ~ %d] total %d batches\n", from, to, batchCount,
	))
}

func (p *progress) end(from, to int64, batchCount uint64) {
	if p == nil {
		return
	}
	p.print(time.Now(), fmt.Sprintf(
		"[%d ~ %d] total %d batches, finished in %s\n", from, to, batchCount, ElapsedTime(p.beginTime),
	))
}

// print a start message, on half newline.
func (p *progress) start(startTime time.Time, from, to int64) {
	if p == nil {
		return
	}
	p.print(startTime, p.batchString(from, to)+" ... ")
}

// print a finish message, on half line.
func (p *progress) finish(startTime time.Time) {
	if p == nil {
		return
	}
	fmt.Fprintln(p.output, " finished in "+ElapsedTime(startTime))
}

// print a start message, on full line.
func (p *progress) startLine(startTime time.Time, from, to int64) {
	if p == nil {
		return
	}
	p.print(startTime, p.batchString(from, to)+" ...\n")
}

// print a finish message, on full line.
func (p *progress) finishLine(startTime time.Time, from, to int64) {
	if p == nil {
		return
	}
	p.print(time.Now(), p.batchString(from, to)+" finished in "+ElapsedTime(startTime)+"\n")
}

func (p *progress) print(now time.Time, msg string) {
	if p == nil {
		return
	}
	fmt.Fprint(p.output, now.Format("15:04:05 ")+msg)
}

func (p *progress) batchString(from, to int64) string {
	if p.single {
		return fmt.Sprintf(p.batch, from)
	}
	return fmt.Sprintf(p.batch, from, to)
}

func ElapsedTime(start time.Time) string {
	return PrettyDuration(time.Since(start))
}

func PrettyDuration(d time.Duration) string {
	switch {
	case d > time.Minute:
		d = d.Round(time.Second)
	case d > time.Second:
		d = d.Round(time.Second / 10)
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / 10)
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / 10)
	}
	return fmt.Sprintf("%7s", d.String())
}
