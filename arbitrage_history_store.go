package main

import (
	"math"
	"sync"
	"time"

	"github.com/sawdustofmind/currency-arbitrage/arbalgo"
)

// ArbitrageHistorySize is number of entries history stored
const ArbitrageHistorySize = 100
const priceScale = 1e-8

type ArbitrageHistoryEntry struct {
	Time   time.Time `json:"time,omitempty"`
	Cycle  string    `json:"cycle,omitempty"`
	Path   []int     `json:"-"`
	Report string    `json:"report,omitempty"`
	Profit float64   `json:"profit,omitempty"`
}

type ArbitrageHistoryStore struct {
	mu      sync.RWMutex
	entries []ArbitrageHistoryEntry
}

func (s *ArbitrageHistoryStore) Add(e *ArbitrageHistoryEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var nearestSameCycle *ArbitrageHistoryEntry
	for i := len(s.entries) - 1; i >= 0; i-- {
		if arbalgo.SameCycles(s.entries[i].Path, e.Path) {
			nearestSameCycle = &s.entries[i]
		}
	}
	if nearestSameCycle == nil || math.Abs(nearestSameCycle.Profit-e.Profit) > priceScale {
		s.put(e)
	}
}

// TODO: think of array deque here
func (s *ArbitrageHistoryStore) put(e *ArbitrageHistoryEntry) {
	if len(s.entries) == ArbitrageHistorySize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, *e)
}

func (s *ArbitrageHistoryStore) Get() []ArbitrageHistoryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.entries[:]
}
