package logs

import cmap "github.com/orcaman/concurrent-map/v2"

type concurrentRepo struct {
	counts cmap.ConcurrentMap[string, int]
}

func NewConcurrentRepo() *concurrentRepo {
	return &concurrentRepo{
		counts: cmap.New[int](),
	}
}

func (r *concurrentRepo) IncrementAttribute(attr string) {
	val, _ := r.counts.Get(attr)
	r.counts.Set(attr, val+1)
}

func (r *concurrentRepo) Reset() map[string]int {
	oldMap := r.counts
	r.counts = cmap.New[int]()
	return oldMap.Items()
}
