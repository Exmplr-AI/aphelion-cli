package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	return s
}