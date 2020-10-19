package node

// workpool is a stable of worker goroutines to save on gouroutine
// initialization (optimization).
type workpool struct {
	workerCount    int
	availableCount int
	chunkSize      int
	workQueue      chan []byte
}

func makeWorkpool(cap int, chunksize int) *workpool {
	return &workpool{workerCount: cap, chunkSize: chunksize, workQueue: make(chan []byte, cap)}
}

func (w *workpool) Write(p []byte) (n int, err error) {
	var last int = 0
	if len(p) > w.chunkSize {
		chunks := (len(p) / w.chunkSize) + 1
		for i := 0; i < chunks; i++ {
			w.workQueue <- p[last : last+w.chunkSize]
			last = last + w.chunkSize
		}
	}
	w.workQueue <- p
	return last, nil
}
