package kafka

import (
	"encoding/json"
	"paychain/internal/blockchain"

	"github.com/IBM/sarama"
)

type Producer struct {
	ap    sarama.AsyncProducer
	topic string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = false
	p, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{ap: p, topic: topic}, nil
}

func (p *Producer) PublishTransaction(tx blockchain.Transaction) error {
	b, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	p.ap.Input() <- &sarama.ProducerMessage{Topic: p.topic, Value: sarama.ByteEncoder(b)}
	return nil
}

func (p *Producer) Close() error { return p.ap.Close() }
