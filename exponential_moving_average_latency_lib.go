// Package exponential_moving_average_latency_lib computes an exponential moving
// average of request durations from a stream of individual measurements.
//
// The EMA is updated as:
//
//     ema = alpha*sample + (1-alpha)*ema
//
// where alpha is the smoothing factor. A smaller alpha gives more weight to
// older samples and produces a smoother, slower-moving average.
//
// The zero value is not usable; use NewEMA to construct a value. The zero value
// has an undefined alpha and count, so its behaviour is unspecified.
package exponential_moving_average_latency_lib

import (
	"errors"
	"math"
)

// EMA holds the running state of an exponential moving average.
type EMA struct {
	alpha   float64
	current float64
	count   uint64
	seeded  bool
}

// NewEMA returns an EMA with the given smoothing factor alpha.
//
// Alpha must be in the open interval (0, 1). Values outside that range are
// rejected because they would either freeze the average (alpha <= 0) or make it
// track the most recent sample exactly (alpha >= 1), neither of which is an
// exponential moving average.
func NewEMA(alpha float64) (*EMA, error) {
	if math.IsNaN(alpha) || math.IsInf(alpha, 0) {
		return nil, errors.New("alpha must be a finite number")
	}
	if alpha <= 0 || alpha >= 1 {
		return nil, errors.New("alpha must be in the open interval (0, 1)")
	}
	return &EMA{alpha: alpha}, nil
}

// Add incorporates a new latency sample into the average.
//
// The first sample seeds the average directly. This is a deliberate choice:
// starting from zero would make the average useless for the first several
// samples, and starting from an arbitrary value would bias the result. Seeding
// with the first observation means the average is always a real latency and
// converges quickly.
//
// Negative samples are rejected because a negative duration is not a duration.
func (e *EMA) Add(sample float64) error {
	if math.IsNaN(sample) || math.IsInf(sample, 0) {
		return errors.New("sample must be a finite number")
	}
	if sample < 0 {
		return errors.New("sample must be non-negative")
	}
	if !e.seeded {
		e.current = sample
		e.seeded = true
	} else {
		e.current = e.alpha*sample + (1-e.alpha)*e.current
	}
	e.count++
	return nil
}

// Value returns the current exponential moving average.
//
// If no samples have been added, Value returns 0 and false. The second return
// value distinguishes "no data yet" from "the average is genuinely zero"
// (which can happen when every sample is zero).
func (e *EMA) Value() (float64, bool) {
	if !e.seeded {
		return 0, false
	}
	return e.current, true
}

// Count returns the number of samples that have been added.
func (e *EMA) Count() uint64 {
	return e.count
}

// Reset clears the accumulated state, leaving the configured alpha in place.
// The next Add will seed the average from scratch.
func (e *EMA) Reset() {
	e.current = 0
	e.count = 0
	e.seeded = false
}
