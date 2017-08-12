package scanner

import (
	"github.com/clickyab/services/assert"
)

// Factory is needed to make a new scanner
type Factory func(string) Scanner

var mainFunc Factory

// Scanner is the interface to scan
type Scanner interface {
	// Next set another iterate
	Next() ([]string, bool)
}

// NewScanner returns a scanner
func NewScanner(pattern string) Scanner {
	assert.NotNil(mainFunc)
	return mainFunc(pattern)
}

// Register is used for mocking
func Register(f Factory) {
	mainFunc = f
}
