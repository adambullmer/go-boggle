package gameboard

import (
	"strings"
)

// Letters is an array of the Letters making up the word in progress.
type Letters []string

// Words is a unique collection of words found in the gameboard
type Words map[string]struct{}

// Groups is a collection of words grouped by their length
type Groups map[int][]string

func (l Letters) String() string {
	return strings.Join(l, "")
}

func (w Words) String() string {
	keys := ""
	for k := range w {
		keys += k
	}

	return keys
}
