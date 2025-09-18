package logs

import "sync"

type Repo interface{
	IncrementAttribute(string)
	Reset() map[string]int
}

type repo struct {
	counts map[string]int
	mu     sync.Mutex
}

func NewRepo() *repo {
	return &repo{
		counts: make(map[string]int),
	}
}

func (r *repo) IncrementAttribute(attr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counts[attr]++
}

func (r *repo) Reset() map[string]int {
	r.mu.Lock()
	defer r.mu.Unlock()
	oldMap := r.counts
	r.counts = map[string]int{}
	return oldMap
}
