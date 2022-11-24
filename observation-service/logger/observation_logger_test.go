package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/caraml-dev/observation-service/observation-service/config"
)

func TestNoopLogConsumer(t *testing.T) {
	logConsumer, err := NewNoopLogConsumer()
	expected := &NoopLogConsumer{}

	assert.NoError(t, nil, err)
	assert.Equal(t, expected, logConsumer)
}

func TestNoopLogProducer(t *testing.T) {
	logProducer, err := NewNoopLogProducer()
	expected := &NoopLogProducer{}

	assert.NoError(t, nil, err)
	assert.Equal(t, expected, logProducer)
}

func TestObservationLogger(t *testing.T) {
	// Configs
	consumerConfig := config.LogConsumerConfig{}
	producerConfig := config.LogProducerConfig{}

	observationLogger, err := NewObservationLogger(
		consumerConfig,
		producerConfig,
	)
	assert.NoError(t, nil, err)

	// Expected
	logConsumer, err := NewNoopLogConsumer()
	assert.NoError(t, nil, err)
	logProducer, err := NewNoopLogProducer()
	assert.NoError(t, nil, err)
	expected := &ObservationLogger{
		logsChannel:   observationLogger.logsChannel,
		consumer:      logConsumer,
		producer:      logProducer,
		flushInterval: time.Duration(producerConfig.FlushIntervalSeconds),
	}

	assert.NoError(t, nil, err)
	assert.Equal(t, expected, observationLogger)
}
