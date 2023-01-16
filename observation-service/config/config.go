package config

import (
	"fmt"

	"github.com/gojek/mlp/api/pkg/instrumentation/newrelic"
	"github.com/gojek/mlp/api/pkg/instrumentation/sentry"

	commonconfig "github.com/caraml-dev/timber/common/config"
)

// Config captures the config related to starting Observation Service
type Config struct {
	Port int `envconfig:"PORT" default:"9001"`

	DeploymentConfig  commonconfig.DeploymentConfig
	NewRelicConfig    newrelic.Config
	SentryConfig      sentry.Config
	LogConsumerConfig LogConsumerConfig
	LogProducerConfig LogProducerConfig
	MonitoringConfig  MonitoringConfig
}

// ObservationLoggerConsumerKind captures the consumer config for reading Observation Service logs
type ObservationLoggerConsumerKind = string

const (
	// LoggerNoopConsumer is a No-Op ObservationLog Consumer
	LoggerNoopConsumer ObservationLoggerConsumerKind = ""
	// LoggerKafkaConsumer is a Kafka ObservationLog Consumer
	LoggerKafkaConsumer ObservationLoggerConsumerKind = "kafka"
)

// LogConsumerConfig captures the config related to consuming ObservationLog via a background process
type LogConsumerConfig struct {
	// The type of Data Source for Observation logs
	Kind ObservationLoggerConsumerKind `default:""`

	// KafkaConfig captures the config related to initializing a Kafka Consumer
	KafkaConfig *KafkaConfig
}

// KafkaConfig captures all configurable parameters when configuring a Kafka Consumer and Producer
type KafkaConfig struct {
	// Kafka Brokers to connect to, comma-delimited, in the form of "<broker_host>:<broker_port>"
	Brokers string
	// Kafka Topic to produce to/consume from
	Topic string
	// Largest record batch size allowed by Kafka (after compression if compression is enabled)
	MaxMessageBytes int `default:"1048588"`
	// The compression type for all data generated by the Producer
	CompressionType string `default:"none"`
	// ConnectTimeoutMS is the maximum duration (ms) the Kafka Producer/Consumer will block for to get Metadata, before timing out
	ConnectTimeoutMS int `default:"1000"`
	// PollInterval is the maximum duration (ms) the Kafka Consumer will block for, before timing out
	PollInterval int `default:"1000"`
	// What to do when there is no initial offset in Kafka or if the current offset does not exist any more on the server
	AutoOffsetReset string `default:"latest"`
}

// FluentdConfig captures necessary configuration to flush ObservationLog to multiple sinks
type FluentdConfig struct {
	// The type of Data Sink for Observation logs flushed via Fluentd
	Kind FluentdProducerSinkKind `default:""`
	// Fluentd Host to connect to
	Host string `default:"localhost"`
	// Fluentd Port to connect to
	Port int `default:"24224"`
	// Fluentd Tag to match messages
	Tag string `default:"observation-service"`

	// BQConfig captures the config related to initializing a BQ Sink
	BQConfig *BQConfig
}

// BQConfig captures GCP BigQuery information for writing Observation Service logs to
type BQConfig struct {
	// GCP Project
	Project string
	// GCP Dataset
	Dataset string
	// GCP Table
	Table string
}

// ObservationLoggerProducerKind captures the producer config for flushing Observation Service logs
type ObservationLoggerProducerKind = string

const (
	// LoggerNoopProducer is a No-Op ObservationLog Producer
	LoggerNoopProducer ObservationLoggerProducerKind = ""
	// LoggerStdOutProducer is a Standard Output ObservationLog Producer
	LoggerStdOutProducer ObservationLoggerProducerKind = "stdout"
	// LoggerKafkaProducer is a Kafka ObservationLog Producer
	LoggerKafkaProducer ObservationLoggerProducerKind = "kafka"
	// LoggerFluentdProducer is a Fluentd ObservationLog Producer
	LoggerFluentdProducer ObservationLoggerProducerKind = "fluentd"
)

// FluentdProducerSinkKind captures the sink config for flushing Observation Service logs via Fluentd
type FluentdProducerSinkKind = string

const (
	// LoggerNoopSinkFluentdProducer is a Fluentd No-Op ObservationLog Producer
	LoggerNoopSinkFluentdProducer FluentdProducerSinkKind = ""
	// LoggerBQSinkFluentdProducer is a Fluentd ObservationLog BigQuery Producer
	LoggerBQSinkFluentdProducer FluentdProducerSinkKind = "bq"
)

// LogProducerConfig captures the config related to producing ObservationLog
type LogProducerConfig struct {
	// The type of Data Sink for Observation logs
	Kind ObservationLoggerProducerKind `default:""`
	// Maximum no. of Observation logs to be stored in-memory prior to flushing to Data sink
	QueueLength int `default:"100"`
	// Duration that specifies how often in-memory Observation logs should be flushed to Data sink
	FlushIntervalSeconds int `default:"1"`

	// KafkaConfig captures the config related to initializing a Kafka Producer
	KafkaConfig *KafkaConfig
	// FluentdConfig captures the config related to initializing a Fluentd instance
	FluentdConfig *FluentdConfig
}

// MetricSinkKind captures type of metrics sink
type MetricSinkKind = string

const (
	// PrometheusMetricSink represents the Prometheus Metric Storage
	PrometheusMetricSink MetricSinkKind = "prometheus"
	// NoopMetricSink represents no Metric Storage
	NoopMetricSink MetricSinkKind = ""
)

// MonitoringConfig captures the config for monitoring metrics
type MonitoringConfig struct {
	// The type of Metrics Sink for Observation logs
	Kind MetricSinkKind `default:""`
}

// ListenAddress returns the Observation API port
func (c *Config) ListenAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}

// Load parses multiple file configs specified via filepaths and returns a Config struct
func Load(filepaths ...string) (*Config, error) {
	var cfg Config
	err := commonconfig.ParseConfig(&cfg, filepaths)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %s", err)
	}

	return &cfg, nil
}
