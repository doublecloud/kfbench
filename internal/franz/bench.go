package franz

import (
	"context"
	"fmt"

	"github.com/doublecloud/kfbench/internal/measure"

	// "github.com/doublecloud/kafka-vs-msk/bench/internal/must"

	"os"

	"github.com/doublecloud/kfbench/internal/value"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ValueGenerator interface {
	Generate(int64) []byte
}

func Produce(c Config) error {
	var vg ValueGenerator
	if c.StaticRecord {
		vg = value.NewStaticBytes(c.RecordBytes)
	} else {
		vg = value.NewBytesGenerator(c.RecordBytes)
	}

	g, o := &measure.Gauge{}, &measure.Gauge{}

	for i := 0; i < int(c.Concurrency); i++ {
		go func() {
			cn, err := kgo.NewClient(c.opts...)
			if err != nil {
				fmt.Fprintln(os.Stderr, "client error: %v", err)
			}
			data := vg.Generate(1)

			for {
				err := g.Take(func() (v []byte, err error) {
					rec := kgo.SliceRecord(data)
					if !c.AsyncMode {
						res := cn.ProduceSync(context.Background(), rec)
						return rec.Value, res.FirstErr()
					} else {
						cn.Produce(context.Background(), rec, func(r *kgo.Record, err error) {
							if err != nil {
								fmt.Printf("record had a produce error: %v\n", err)
								// How do I return an error here and exit sendAll function?
							}
						})
					}
					return rec.Value, nil
				})
				if err != nil {
					fmt.Fprintln(os.Stderr, "produce error: %v", err)
				}
			}
		}()
	}
	measure.NewPrinter(g, o, c.ReportingInterval).Report(os.Stdout, c.Duration)
	return nil
}
