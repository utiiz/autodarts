package models

type Player struct {
	Name    string
	Score   int
	Darts   [3]Dart
	History []Dart
}
