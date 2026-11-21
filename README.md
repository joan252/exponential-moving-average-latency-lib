# exponential-moving-average-latency-lib

Computes an exponential moving average of request durations from a stream of individual measurements.

```go
package main

import (
	"fmt"

	ema "github.com/example/exponential-moving-average-latency-lib"
)

func main() {
	avg, err := ema.NewEMA(0.2)
	if err != nil {
		panic(err)
	}

	for _, d := range []float64{12.5, 13.1, 11.9, 14.0} {
		if err := avg.Add(d); err != nil {
			panic(err)
		}
	}

	v, ok := avg.Value()
	fmt.Println(v, ok, avg.Count())
}
```

## Why this exists

Latency is noisy. A single slow request can make an average jump, and a plain
mean over a long window hides recent changes. An exponential moving average
weights recent samples more heavily while still smoothing spikes, which is what
you want when feeding a latency signal into a circuit breaker or an autoscaler.

The trade-off is that EMA has a smoothing factor `alpha` with no universally
correct value. Small `alpha` is stable but slow to react; large `alpha` reacts
quickly but is jumpy. This library does not pick for you: you pass `alpha` to
`NewEMA` and it must be in the open interval `(0, 1)`.

## The awkward edge

The first sample seeds the average directly instead of blending against an
initial zero. That means `Value()` after one sample is exactly that sample, not
a fraction of it. It also means `Value()` returns a second boolean: before any
samples have been added there is no average, and the boolean is `false`. Do not
ignore it if a zero latency is a meaningful value in your system.

## API

- `NewEMA(alpha float64) (*EMA, error)` — construct an average. `alpha` must be
  finite and strictly between 0 and 1.
- `(*EMA) Add(sample float64) error` — add one latency measurement. Samples must
  be finite and non-negative.
- `(*EMA) Value() (float64, bool)` — current average and whether any samples have
  been added.
- `(*EMA) Count() uint64` — number of accepted samples.
- `(*EMA) Reset()` — clear all samples while keeping `alpha`.
