package recommendationservice

import "sync"

type RecommendationStore struct {
	mu              sync.Mutex
	recommendations map[string]map[string]int
}

func NewStore() *RecommendationStore {
	return &RecommendationStore{
		recommendations: make(map[string]map[string]int),
	}
}

func (rs *RecommendationStore) Track(userID, category string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, ok := rs.recommendations[userID]; !ok {
		rs.recommendations[userID] = make(map[string]int)
	}
	rs.recommendations[userID][category]++
}

func (rs *RecommendationStore) GetRecommendation(userID string) string {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	top := ""
	max := 0

	for cat, cnt := range rs.recommendations[userID] {
		if cnt > max {
			max = cnt
			top = cat
		}
	}
	return top
}
