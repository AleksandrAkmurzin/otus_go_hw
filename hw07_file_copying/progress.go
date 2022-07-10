package main

import (
	"errors"
	"fmt"
	"math"
)

type Progress struct {
	currentPercent float64
	total          int64
	isFinished     bool
}

func newProgress(total int64) Progress {
	return Progress{total: total}
}

func (p *Progress) setCurrent(n int64) error {
	if n > p.total {
		return errors.New("current progress cannot be more than total")
	}

	percent := math.Floor(float64(n) * 100 / float64(p.total))
	if percent < p.currentPercent {
		return errors.New("current progress cannot be less than previous")
	}

	// Skip duplicating values if progress is too long.
	if percent == p.currentPercent {
		return nil
	}

	p.currentPercent = percent
	fmt.Print(percent, "%..")

	if percent == 100 {
		p.finish()
	}

	return nil
}

func (p *Progress) finish() error {
	if p.isFinished {
		return errors.New("progress was already finished")
	}

	p.isFinished = true
	fmt.Println("Done!")
	return nil
}
