package models

import "strconv"

type Dart struct {
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Score Score   `json:"score"`
}

type Score struct {
	Bed     string `json:"bed"`
	Segment int    `json:"segment"`
	Score   int    `json:"score"`
}

func (s Score) String() string {
	return s.Bed + strconv.Itoa(s.Segment)
}
