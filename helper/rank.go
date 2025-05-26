package helper

import (
	"container/heap"
	"sort"
	"sync"
)

type RankItem[T any] struct {
	Index int
	Score float64
	Item  T
}

type Ranker[T any] interface {
	Items() []RankItem[T]
	At(index int) (RankItem[T], bool)
	Add(item T) RankItem[T]
}

type heapItem[T any] struct {
	id    string
	score float64
	item  T
}

type scoreHeap[T any] []*heapItem[T]

func (h *scoreHeap[T]) Len() int {
	return len(*h)
}

func (h *scoreHeap[T]) Less(i, j int) bool {
	return (*h)[i].score < (*h)[j].score
}

func (h *scoreHeap[T]) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *scoreHeap[T]) Push(x any) {
	*h = append(*h, x.(*heapItem[T]))
}

func (h *scoreHeap[T]) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

type QueueRank[T any] struct {
	mu       sync.Mutex
	capacity int
	heap     scoreHeap[T]
	repo     Repo[T]
	id       GenID
	Score    func(T) float64
}

func NewQueueRank[T any](repo Repo[T], capacity int, score func(T) float64) *QueueRank[T] {
	h := make(scoreHeap[T], 0, capacity)
	repo.Range(func(key string, value T) bool {
		h = append(h, &heapItem[T]{
			id:    key,
			score: score(value),
			item:  value,
		})
		return true
	})
	heap.Init(&h)
	return &QueueRank[T]{
		repo:     repo,
		id:       NewGenRandomRepoShortID(4, 16, 5, repo),
		capacity: capacity,
		heap:     h,
		Score:    score,
	}
}

func (q *QueueRank[T]) Add(item T) RankItem[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	id, _ := q.id.NextID()
	score := q.Score(item)
	hi := &heapItem[T]{
		id:    id,
		score: score,
		item:  item,
	}
	heap.Push(&q.heap, hi)
	q.repo.Put(id, item)

	if q.heap.Len() > q.capacity {
		p := heap.Pop(&q.heap).(*heapItem[T])
		q.repo.Del(p.id)
	}

	snapshot := make([]*heapItem[T], q.heap.Len())
	copy(snapshot, q.heap)

	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].score > snapshot[j].score
	})

	var resultRank RankItem[T]
	items := make([]RankItem[T], len(snapshot))
	for idx, node := range snapshot {
		items[idx] = RankItem[T]{
			Index: idx,
			Score: node.score,
			Item:  node.item,
		}
		if node == hi {
			resultRank = items[idx]
		}
	}
	return resultRank
}

func (q *QueueRank[T]) Items() []RankItem[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	snapshot := make([]*heapItem[T], q.heap.Len())
	copy(snapshot, q.heap)
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].score > snapshot[j].score
	})

	items := make([]RankItem[T], len(snapshot))
	for idx, node := range snapshot {
		items[idx] = RankItem[T]{
			Index: idx,
			Score: node.score,
			Item:  node.item,
		}
	}
	return items
}

func (q *QueueRank[T]) At(index int) (RankItem[T], bool) {
	items := q.Items()
	if index < 0 || index >= len(items) {
		return RankItem[T]{}, false
	}
	return items[index], true
}
