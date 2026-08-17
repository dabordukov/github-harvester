package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"repo-stat/collector/internal/dto"
	"repo-stat/collector/internal/usecase"
	platformkafka "repo-stat/platform/kafka"

	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *confluentkafka.Producer
	topic    string
	log      *slog.Logger
}

type KafkaConsumer struct {
	consumer *confluentkafka.Consumer
	topic    string
	log      *slog.Logger
}

func NewKafkaProducer(brokers []string, topic string, log *slog.Logger) (*KafkaProducer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("brokers is required")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic is required")
	}
	producer, err := confluentkafka.NewProducer(&confluentkafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
		"client.id":         "collector-producer",
	})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer, topic: topic, log: log}, nil
}

func NewKafkaConsumer(brokers []string, topic, groupID string, log *slog.Logger) (*KafkaConsumer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("brokers is required")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic is required")
	}
	if groupID == "" {
		return nil, fmt.Errorf("groupID is required")
	}
	consumer, err := confluentkafka.NewConsumer(&confluentkafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	if err := consumer.Subscribe(topic, nil); err != nil {
		_ = consumer.Close()
		return nil, err
	}
	return &KafkaConsumer{consumer: consumer, topic: topic, log: log}, nil
}

func (p *KafkaProducer) PublishCollectRequest(ctx context.Context, owner, repo string) error {
	if p == nil || p.producer == nil {
		return fmt.Errorf("kafka producer is not initialized")
	}

	payload, err := json.Marshal(platformkafka.CollectRequest{
		RequestID: fmt.Sprintf("%s/%s/%d", owner, repo, time.Now().UnixNano()),
		Owner:     owner,
		RepoName:  repo,
	})
	if err != nil {
		return fmt.Errorf("marshal collect request: %w", err)
	}

	deliveryChan := make(chan confluentkafka.Event, 1)
	if err := p.producer.Produce(&confluentkafka.Message{
		TopicPartition: confluentkafka.TopicPartition{Topic: &p.topic, Partition: confluentkafka.PartitionAny},
		Value:          payload,
	}, deliveryChan); err != nil {
		return fmt.Errorf("produce collect request: %w", err)
	}

	select {
	case <-deliveryChan:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for kafka delivery")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *KafkaProducer) PublishCollectResult(ctx context.Context, req usecase.CollectRequest, repo *dto.RepositoryDTO) error {
	if p == nil || p.producer == nil {
		return fmt.Errorf("kafka producer is not initialized")
	}
	if repo == nil {
		return fmt.Errorf("repository data is empty")
	}

	payload, err := json.Marshal(platformkafka.CollectResult{
		RequestID:    req.RequestID,
		Owner:        req.Owner,
		RepoName:     req.RepoName,
		Description:  repo.Description,
		Forks:        repo.Forks,
		Stars:        repo.Stars,
		CreatedAt:    repo.CreatedAt,
		CommitsCount: repo.CommitsCount,
	})
	if err != nil {
		return fmt.Errorf("marshal collect result: %w", err)
	}

	deliveryChan := make(chan confluentkafka.Event, 1)
	if err := p.producer.Produce(&confluentkafka.Message{
		TopicPartition: confluentkafka.TopicPartition{Topic: &p.topic, Partition: confluentkafka.PartitionAny},
		Value:          payload,
	}, deliveryChan); err != nil {
		return fmt.Errorf("produce collect result: %w", err)
	}

	select {
	case <-deliveryChan:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for kafka delivery")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *KafkaProducer) Close() error {
	if p == nil || p.producer == nil {
		return nil
	}
	p.producer.Flush(1000)
	p.producer.Close()
	return nil
}

func (c *KafkaConsumer) StartCollectLoop(ctx context.Context, handler func(context.Context, usecase.CollectRequest) error) {
	if c == nil || c.consumer == nil {
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msg, err := c.consumer.ReadMessage(1 * time.Second)
			if err != nil {
				if kerr, ok := err.(confluentkafka.Error); !ok || kerr.Code() != confluentkafka.ErrTimedOut {
					if c.log != nil {
						c.log.Error("read kafka request failed", "error", err)
					}
				}
				continue
			}

			var req platformkafka.CollectRequest
			if err := json.Unmarshal(msg.Value, &req); err != nil {
				if c.log != nil {
					c.log.Error("decode kafka request failed", "error", err)
				}
				continue
			}

			if err := handler(ctx, usecase.CollectRequest{RequestID: req.RequestID, Owner: req.Owner, RepoName: req.RepoName}); err != nil {
				if c.log != nil {
					c.log.Error("handle kafka request failed", "error", err, "owner", req.Owner, "repo", req.RepoName)
				}
			}
		}
	}()
}

func (c *KafkaConsumer) Close() error {
	if c == nil || c.consumer == nil {
		return nil
	}
	return c.consumer.Close()
}
