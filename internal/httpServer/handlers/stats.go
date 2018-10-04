package handlers

import (
	"encoding/json"
	"sync"
	"time"
)

type Stats struct {
	sync.RWMutex
	timeTaken         time.Duration
	Total             int64   `json:"total"`
	TimeTakenMicroSec int64 `json:"average"`
}

func (s *Stats) Update(taken time.Duration) {
	s.Lock()
	defer s.Unlock()

	s.Total++
	s.timeTaken += taken
}

func (s *Stats) ToJSON() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	if s.Total != 0 {
		s.TimeTakenMicroSec = convertNsToMicroSeconds(s.timeTaken) / s.Total

	} else {
		s.TimeTakenMicroSec = 0
	}

	b, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

func convertNsToMicroSeconds(t time.Duration) int64 {
	return t.Nanoseconds() / 1000
}
