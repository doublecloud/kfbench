package measure

import (
	"fmt"
	"io"
	"time"
)

type Printer struct {
	g        *Gauge // interval gauge
	o        *Gauge // overall gauge
	interval time.Duration
	start    time.Time
}

func NewPrinter(g *Gauge, o *Gauge, interval time.Duration) *Printer {
	return &Printer{g: g, o: o, interval: interval, start: time.Now()}
}

func (p *Printer) Print(out io.Writer, d Dump, i time.Duration) {
	message := fmt.Sprintf(
		"[%d]\t%s\t(Max: %s,\tStdDev: %s);\t%0.2f MiB/s;\t%0.2fk records/s\t%0.3f errors",
		len(d.Samples),
		d.MeanLatency().Round(time.Nanosecond),
		d.Max.Round(time.Microsecond),
		d.StdDev(),
		d.BytesRate(i)/(1024*1024),
		d.RecordsRate(i)/1000,
		d.ErrorsRate(i),
	)

	_, _ = fmt.Fprintln(out, message)
}

func (p *Printer) Report(out io.Writer, d time.Duration) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for range ticker.C {
		dump := p.g.Dump()
		p.o.Sum(dump)
		p.Print(out, dump, p.interval)
		if time.Since(p.start) > d {
			break
		}
	}

	dump := p.o.Dump()
	fmt.Fprintln(out, "=======overall stats=======")
	p.Print(out, dump, time.Since(p.start))
}
