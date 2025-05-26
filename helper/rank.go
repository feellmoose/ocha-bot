package helper

import (
	"container/heap"
	"sort"
	"sync"
)

type RankItem struct {
	Index int
	Score float64
	Item  any
}

type Ranker interface {
	Items() []RankItem
	At(index int) (RankItem, bool)
	Add(item any) RankItem
}

type heapItem struct {
	id    string
	score float64
	item  any
}

type scoreHeap []*heapItem

func (h *scoreHeap) Len() int {
	return len(*h)
}

func (h *scoreHeap) Less(i, j int) bool {
	return (*h)[i].score < (*h)[j].score
}

func (h *scoreHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *scoreHeap) Push(x any) {
	*h = append(*h, x.(*heapItem))
}

func (h *scoreHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

type QueueRank struct {
	mu       sync.Mutex
	capacity int
	heap     scoreHeap
	repo     Repo
	id       GenID
	Score    func(any) float64
}

func NewQueueRank(repo Repo, capacity int, score func(any) float64) *QueueRank {
	h := make(scoreHeap, 0, capacity)
	repo.Range(func(key, value any) bool {
		h = append(h, &heapItem{
			id:    key.(string),
			score: score(value),
			item:  value,
		})
		return true
	})
	heap.Init(&h)
	return &QueueRank{
		repo:     repo,
		id:       NewGenRandomRepoShortID(4, 16, 5, repo),
		capacity: capacity,
		heap:     h,
		Score:    score,
	}
}

func (q *QueueRank) Add(item any) RankItem {
	q.mu.Lock()
	defer q.mu.Unlock()

	id, _ := q.id.NextID()
	score := q.Score(item)
	hi := &heapItem{
		id:    id,
		score: score,
		item:  item,
	}
	heap.Push(&q.heap, hi)
	q.repo.Put(id, item)

	if q.heap.Len() > q.capacity {
		p := heap.Pop(&q.heap).(*heapItem)
		q.repo.Del(p.id)
	}

	snapshot := make([]*heapItem, q.heap.Len())
	copy(snapshot, q.heap)

	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].score > snapshot[j].score
	})

	var resultRank RankItem
	items := make([]RankItem, len(snapshot))
	for idx, node := range snapshot {
		items[idx] = RankItem{
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

func (q *QueueRank) Items() []RankItem {
	q.mu.Lock()
	defer q.mu.Unlock()

	snapshot := make([]*heapItem, q.heap.Len())
	copy(snapshot, q.heap)
	sort.Slice(snapshot, func(i, j int) bool {
		return snapshot[i].score > snapshot[j].score
	})

	items := make([]RankItem, len(snapshot))
	for idx, node := range snapshot {
		items[idx] = RankItem{
			Index: idx,
			Score: node.score,
			Item:  node.item,
		}
	}
	return items
}

func (q *QueueRank) At(index int) (RankItem, bool) {
	items := q.Items()
	if index < 0 || index >= len(items) {
		return RankItem{}, false
	}
	return items[index], true
}
