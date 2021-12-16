package batches

func (b Batches) iterator() iterator {
	var step, step2 = b.step()
	totalCount := (b.To-b.From)/step + 1 // step() ensures totalCount must be positive
	return iterator{
		asc:        b.From <= b.To,
		step:       step,
		step2:      step2,
		from:       b.From,
		to:         b.To,
		totalCount: uint64(totalCount),
	}
}

// from + step = nextFrom
// from + step2 = to
func (b Batches) step() (step, step2 int64) {
	step = int64(b.BatchSize)
	if step == 0 {
		step = 1
	}
	step2 = step - 1
	if b.From > b.To {
		step = -step
		step2 = -step2
	}
	return step, step2
}

type iterator struct {
	asc         bool
	from, to    int64
	step, step2 int64
	totalCount  uint64
}

func (itr *iterator) next() (from, to int64, ok bool) {
	from = itr.from
	if itr.asc && from > itr.to || !itr.asc && from < itr.to {
		return -1, -1, false
	}
	to = from + itr.step2
	if itr.asc && to > itr.to || !itr.asc && to < itr.to {
		to = itr.to
	}

	itr.from += itr.step
	return from, to, true
}

func (itr *iterator) isSingleStep() bool {
	return itr.step2 == 0
}
