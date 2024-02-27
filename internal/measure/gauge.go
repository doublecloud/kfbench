package measure

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type (
	Gauge struct {
		Samples []time.Duration
		Bytes   int64
		Errors  int64

		min, max time.Duration
		mx       sync.Mutex
	}

	Dump struct {
		Samples []time.Duration
		Bytes   int64
		Errors  int64

		Min, Max time.Duration
	}
)

func (g *Gauge) Dump() Dump {
	g.mx.Lock()
	defer func() {
		g.Samples = nil
		g.Bytes = 0
		g.Errors = 0
		g.min = 0
		g.max = 0
		g.mx.Unlock()
	}()

	return Dump{
		Samples: g.Samples,
		Bytes:   g.Bytes,
		Errors:  g.Errors,
		Min:     g.min,
		Max:     g.max,
	}
}

func (g *Gauge) Sum(d Dump) {
	g.mx.Lock()
	defer g.mx.Unlock()

	g.Samples = append(g.Samples, d.Samples...)
	g.Bytes += d.Bytes
	g.Errors += d.Errors
	if d.Max > g.max {
		g.max = d.Max
	}
	if d.Min < g.min {
		g.min = d.Min
	}
}

func (g *Gauge) Take(f func() ([]byte, error)) error {
	t := time.Now()
	b, err := f()
	s := time.Since(t)

	g.mx.Lock()
	defer g.mx.Unlock()
	if err != nil {
		// Do not count error samples as usual ones
		g.Errors++
		return err
	}
	g.Samples = append(g.Samples, s)
	g.Bytes += int64(len(b))
	if s > g.max {
		g.max = s
	} else if g.min == 0 || s < g.min {
		g.min = s
	}
	return nil
}

func (d Dump) MeanLatency() time.Duration {
	total := len(d.Samples)
	if total == 0 {
		return 0
	}

	sum := time.Duration(0)
	for _, s := range d.Samples {
		sum += s
	}

	return sum / time.Duration(total)
}

func (d Dump) StdDev() time.Duration {
	total := len(d.Samples)
	if total == 0 {
		return 0
	}

	var sd float64
	mean := d.MeanLatency().Microseconds()
	for _, s := range d.Samples {
		sd += math.Pow(float64(s.Microseconds()-mean), 2)
	}
	sd = math.Sqrt(sd / float64(total))
	x, _ := time.ParseDuration(fmt.Sprintf("%fµs", sd))
	return x
}

// RecordsRate per second records rate
func (d Dump) RecordsRate(interval time.Duration) float64 {
	total := len(d.Samples)
	if total == 0 {
		return 0
	}

	return float64(total) / interval.Seconds()
}

func (d Dump) ErrorsRate(interval time.Duration) float64 {
	samples := int64(len(d.Samples))
	errors := d.Errors

	if samples+errors == 0 {
		return 0.0
	}

	return float64(errors) * 100.0 / float64(samples+errors)
}

// BytesRate per second bytes rate
func (d Dump) BytesRate(interval time.Duration) float64 {
	if d.Bytes == 0 {
		return 0
	}

	return float64(d.Bytes) / interval.Seconds()
}
