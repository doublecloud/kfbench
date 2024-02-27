package franz

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

func CreateTopic(c Config, name string) error {
	cl, err := kgo.NewClient(c.opts...)
	if err != nil {
		return fmt.Errorf("connection error: %v", err)
	}
	defer cl.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	admc := kadm.NewClient(cl)
	topics, err := admc.ListTopics(ctx)
	if err != nil {
		return err
	}

	if topics.Has(name) {
		fmt.Printf("[Debug] Topic [%s] already exist, skipping creation\n", name)
		return nil
	}

	bb, err := admc.ListBrokers(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve brokers list: %w", err)
	}

	fmt.Printf("[Debug] Creating topic [%s]\n", name)
	vPtr := func(v string) *string {
		return &v
	}
	_, err = admc.CreateTopic(ctx, int32(c.Partitions), int16(len(bb)), map[string]*string{
		// "retention.ms":    vPtr("21600000"),
		// "retention.bytes": vPtr("5368709120"), //5 GB
		"retention.ms":    vPtr("5000"),
		"retention.bytes": vPtr("104857600"), // 1 GB
		// "file.delete.delay.ms": vPtr("0"),
	}, name)
	return err
}
