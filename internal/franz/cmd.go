package franz

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	var c Config
	return &cli.Command{
		Name: "franz",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "brokers",
				Required: true,
				Action: func(ctx *cli.Context, s string) error {
					c.Brokers = strings.Split(s, ",")
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "topic",
				Required:    true,
				Destination: &c.Topic,
			},
			&cli.BoolFlag{
				Name:        "create-topic",
				Value:       true,
				Destination: &c.CreateTopic,
			},
			&cli.Int64Flag{
				Name:        "record-bytes",
				Value:       1024,
				Destination: &c.RecordBytes,
			},
			&cli.BoolFlag{
				Name:        "static-record",
				Value:       true,
				Destination: &c.StaticRecord,
			},
			&cli.StringFlag{
				Name:        "compression",
				Value:       "none",
				Destination: &c.Compression,
			},
			&cli.StringFlag{
				Name:        "sasl-method",
				Destination: &c.SaslMethod,
			},
			&cli.StringFlag{
				Name:        "sasl-user",
				Destination: &c.SaslUser,
			},
			&cli.StringFlag{
				Name:        "sasl-password",
				Destination: &c.SaslPassword,
			},
			&cli.BoolFlag{
				Name:        "tls",
				Destination: &c.TLS,
			},
			&cli.BoolFlag{
				Name:        "insecure",
				Destination: &c.InsecureTLS,
			},
			&cli.Int64Flag{
				Name:        "concurrency",
				Value:       1,
				Destination: &c.Concurrency,
			},
			&cli.BoolFlag{
				Name:        "async",
				Destination: &c.AsyncMode,
			},
			&cli.Int64Flag{
				Name:        "partitions",
				Value:       1,
				Destination: &c.Partitions,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Value:       "info",
				Destination: &c.LogLevel,
			},
			&cli.DurationFlag{
				Name:        "reporting-interval",
				Value:       1 * time.Second,
				Destination: &c.ReportingInterval,
			},
			&cli.DurationFlag{
				Name:        "duration",
				Value:       100 * time.Second,
				Destination: &c.Duration,
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := c.init(); err != nil {
				return err
			}
			if c.CreateTopic {
				if err := CreateTopic(c, c.Topic); err != nil {
					return fmt.Errorf("failed to create topic: %w", err)
				}
			}
			return Produce(c)
		},
	}
}
