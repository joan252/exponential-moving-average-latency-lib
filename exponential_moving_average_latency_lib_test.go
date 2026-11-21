package exponential_moving_average_latency_lib

import (
	"math"
	"testing"
)

const eps = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < eps
}

func TestNewEMARejectsInvalidAlpha(t *testing.T) {
	cases := []float64{0, 1, -0.1, 1.5, math.NaN(), math.Inf(1), math.Inf(-1)}
	for _, alpha := range cases {
		if _, err := NewEMA(alpha); err == nil {
			t.Errorf("NewEMA(%v) = nil error, want error", alpha)
		}
	}
}

func TestNewEMARejectsNaN(t *testing.T) {
	if _, err := NewEMA(math.NaN()); err == nil {
		t.Fatal("NewEMA(NaN) = nil error, want error")
	}
}

func TestValueBeforeAnySamples(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	v, ok := e.Value()
	if ok {
		t.Fatalf("Value() = %v, true; want 0, false", v)
	}
	if v != 0 {
		t.Fatalf("Value() = %v, want 0", v)
	}
	if e.Count() != 0 {
		t.Fatalf("Count() = %d, want 0", e.Count())
	}
}

func TestFirstSampleSeedsAverage(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	if err := e.Add(100); err != nil {
		t.Fatalf("Add: %v", err)
	}
	v, ok := e.Value()
	if !ok {
		t.Fatal("Value() = _, false; want _, true")
	}
	if !almostEqual(v, 100) {
		t.Fatalf("Value() = %v, want 100", v)
	}
	if e.Count() != 1 {
		t.Fatalf("Count() = %d, want 1", e.Count())
	}
}

func TestSecondSampleAppliesAlpha(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	_ = e.Add(100)
	_ = e.Add(200)
	v, _ := e.Value()
	want := 0.5*200 + 0.5*100
	if !almostEqual(v, want) {
		t.Fatalf("Value() = %v, want %v", v, want)
	}
}

func TestAlphaOneTracksLatestSample(t *testing.T) {
	// alpha is strictly less than 1, so even a tiny alpha should not
	// produce a value equal to the latest sample after the first update.
	e, err := NewEMA(0.999)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	_ = e.Add(100)
	_ = e.Add(200)
	v, _ := e.Value()
	want := 0.999*200 + 0.001*100
	if !almostEqual(v, want) {
		t.Fatalf("Value() = %v, want %v", v, want)
	}
}

func TestAddRejectsNegativeSample(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	if err := e.Add(-1); err == nil {
		t.Fatal("Add(-1) = nil error, want error")
	}
	if e.Count() != 0 {
		t.Fatalf("Count() = %d, want 0 after rejected sample", e.Count())
	}
}

func TestAddRejectsNaN(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	if err := e.Add(math.NaN()); err == nil {
		t.Fatal("Add(NaN) = nil error, want error")
	}
}

func TestAddRejectsInf(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	if err := e.Add(math.Inf(1)); err == nil {
		t.Fatal("Add(+Inf) = nil error, want error")
	}
	if err := e.Add(math.Inf(-1)); err == nil {
		t.Fatal("Add(-Inf) = nil error, want error")
	}
}

func TestZeroSampleIsARealSample(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	_ = e.Add(0)
	v, ok := e.Value()
	if !ok {
		t.Fatal("Value() = _, false; want _, true")
	}
	if v != 0 {
		t.Fatalf("Value() = %v, want 0", v)
	}
	if e.Count() != 1 {
		t.Fatalf("Count() = %d, want 1", e.Count())
	}
}

func TestResetClearsState(t *testing.T) {
	e, err := NewEMA(0.5)
	if err != nil {
		t.Fatalf("NewEMA: %v", err)
	}
	_ = e.Add(100)
	_ = e.Add(200)
	e.Reset()
	if e.Count() != 0 {
		t.Fatalf("Count() = %d, want 0 after Reset", e.Count())
	}
	if _, ok := e.Value(); ok {
		t.Fatal("Value() = _, true after Reset; want _, false")
	}
	// After reset the next sample seeds again.
	_ = e.Add(50)
	v, ok := e.Value()
	if !ok {
		t.Fatal("Value() = _, false after seed; want _, true")
	}
	if !almostEqual(v, 50) {
		t.Fatalf("Value() = %v, want 50", v)
	}
}
