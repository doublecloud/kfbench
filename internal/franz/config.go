package franz

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/aws"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"golang.org/x/exp/slices"
)

type Config struct {
	Brokers           []string
	Topic             string
	CreateTopic       bool
	RecordBytes       int64
	StaticRecord      bool
	Compression       string
	SaslMethod        string
	SaslUser          string
	SaslPassword      string
	TLS               bool
	InsecureTLS       bool
	Concurrency       int64
	Partitions        int64
	AsyncMode         bool
	LogLevel          string
	ReportingInterval time.Duration
	Duration          time.Duration

	opts []kgo.Opt
}

func (c *Config) init() error {
	c.opts = []kgo.Opt{
		kgo.SeedBrokers(c.Brokers...),
		kgo.DefaultProduceTopic(c.Topic),
		kgo.MaxBufferedRecords(250<<20/int(c.RecordBytes) + 1),
		kgo.MaxConcurrentFetches(3),
		// We have good compression, so we want to limit what we read
		// back because snappy deflation will balloon our memory usage.
		kgo.FetchMaxBytes(5 << 20),
		kgo.ProducerBatchMaxBytes(1000000), // the maximum batch size to allow per-partition (must be less than Kafka's max.message.bytes, producing)
	}

	switch strings.ToLower(c.LogLevel) {
	case "":
	case "debug":
		c.opts = append(c.opts, kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelDebug, nil)))
	case "info":
		c.opts = append(c.opts, kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, nil)))
	case "warn":
		c.opts = append(c.opts, kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelWarn, nil)))
	case "error":
		c.opts = append(c.opts, kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelError, nil)))
	default:
		return fmt.Errorf("unrecognized log level %s", c.LogLevel)
	}

	switch strings.ToLower(c.Compression) {
	case "", "none":
		c.opts = append(c.opts, kgo.ProducerBatchCompression(kgo.NoCompression()))
	case "gzip":
		c.opts = append(c.opts, kgo.ProducerBatchCompression(kgo.GzipCompression()))
	case "snappy":
		c.opts = append(c.opts, kgo.ProducerBatchCompression(kgo.SnappyCompression()))
	case "lz4":
		c.opts = append(c.opts, kgo.ProducerBatchCompression(kgo.Lz4Compression()))
	case "zstd":
		c.opts = append(c.opts, kgo.ProducerBatchCompression(kgo.ZstdCompression()))
	default:
		return fmt.Errorf("unrecognized compression %s", c.Compression)
	}

	if c.TLS {
		c.opts = append(c.opts, kgo.DialTLSConfig(&tls.Config{
			InsecureSkipVerify: c.InsecureTLS,
		}))
	}

	if (c.SaslMethod + c.SaslUser + c.SaslPassword) != "" {
		if c.SaslMethod == "" || c.SaslUser == "" || c.SaslPassword == "" {
			return errors.New("all of -sasl-method, -sasl-user, -sasl-password must be specified if any are")
		}

		method := strings.ToLower(c.SaslMethod)
		method = strings.ReplaceAll(method, "-", "")
		method = strings.ReplaceAll(method, "_", "")
		if !slices.Contains([]string{"plain", "scramsha256", "scramsha512"}, method) {
			return fmt.Errorf("unrecognized sasl method: %s", c.SaslMethod)
		}

		switch method {
		case "plain":
			c.opts = append(c.opts, kgo.SASL(plain.Auth{
				User: c.SaslUser,
				Pass: c.SaslPassword,
			}.AsMechanism()))
		case "scramsha256":
			c.opts = append(c.opts, kgo.SASL(scram.Auth{
				User: c.SaslUser,
				Pass: c.SaslPassword,
			}.AsSha256Mechanism()))
		case "scramsha512":
			c.opts = append(c.opts, kgo.SASL(scram.Auth{
				User: c.SaslUser,
				Pass: c.SaslPassword,
			}.AsSha512Mechanism()))
		case "awsmskiam":
			c.opts = append(c.opts, kgo.SASL(aws.Auth{
				AccessKey: c.SaslUser,
				SecretKey: c.SaslPassword,
			}.AsManagedStreamingIAMMechanism()))
		default:
			return fmt.Errorf("unrecognized sasl method %s", c.SaslMethod)
		}
	}

	return nil
}
