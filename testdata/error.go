package testdata

import "time"

type TFrameError struct {
	error
	Timing bool
	Delay *time.Duration
}


func (psError *TFrameError) MIsTiming() bool {
	return psError.Timing
}

func (psError *TFrameError) MGetFrameLeft() *time.Duration {
	return psError.Delay
}
