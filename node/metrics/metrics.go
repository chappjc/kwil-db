package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	// RPC metrics
	requests    metric.Int64Counter
	latencyHist metric.Float64Histogram

	// DB metrics
	dbConnsActive      metric.Int64UpDownCounter
	dbQueryLatencyHist metric.Float64Histogram
	dbQueryErrorCount  metric.Int64Counter

	// Engine metrics
	engineNumNamespaces metric.Int64Gauge

	// Accounts metrics
	// accountsNum metric.Int64ObservableGauge // callback should get account count?
)

// These are the global instances of the metrics grouped by the meter name.
// Until and unless Start is called, these are no-op meters.
var (
	RPC RPCMetrics = rpcMetrics{}
	DB  DBMetrics  = dbMetrics{}
)

const (
	// DBMeterName is the name of the meter for DB metrics. Use the global DB
	// instance to access use these metrics.
	DBMeterName = "kwil-db/postgres"

	// RPCMeterName is the name of the meter for RPC metrics. Use the global RPC
	// instance to access use these metrics.
	RPCMeterName = "kwil-db/rpc"

	EngineMeterName = "kwil-db/engine"

	ConsensusMeterName = "kwil-db/consensus"

	BlockProcessorMeterName = "kwil-db/blocks"

	NodeMeterName = "kwil-db/node"

	MempoolMeterName = "kwil-db/mempool"

	BlockStoreMeterName = "kwil-db/blockstore"

	AccountsMeterName = "kwil-db/accounts"
)

// init sets up all meters and instruments. Initially, the no-op meter
// provider is used until and unless the actual OTEL providers and exporters are
// configured and started with Start.
func init() {
	// DB metrics
	dbMeter := otel.Meter(DBMeterName)
	// active connections from the DB connection pool
	dbConnsActive, _ = dbMeter.Int64UpDownCounter("connections.active")
	dbQueryLatencyHist, _ = dbMeter.Float64Histogram("query.latency")
	dbQueryErrorCount, _ = dbMeter.Int64Counter("query.errors")

	// RPC metrics
	rpcMeter := otel.Meter(RPCMeterName)
	requests, _ = rpcMeter.Int64Counter("requests.total")
	latencyHist, _ = rpcMeter.Float64Histogram("requests.duration")
}

type DBMetrics interface {
	AcquiredConnections(ctx context.Context, dbName string)
	ReleasedConnection(ctx context.Context)
	RecordQuery(ctx context.Context, crudType string, duration time.Duration)
	RecordQueryFailure(ctx context.Context, crudType string, err error)
}

type dbMetrics struct{}

// AcquiredConnections logs a new connection to the database
func (dbMetrics) AcquiredConnections(ctx context.Context, dbName string) {
	// include attribute for the db name
	dbConnsActive.Add(ctx, 1, metric.WithAttributes(attribute.String("db_name", dbName)))
}

// ReleasedConnection logs a connection to the database being released
func (dbMetrics) ReleasedConnection(ctx context.Context) {
	dbConnsActive.Add(ctx, -1)
}

func (dbMetrics) RecordQuery(ctx context.Context, crudType string, duration time.Duration) {
	dbQueryLatencyHist.Record(ctx, 1000*duration.Seconds(),
		metric.WithAttributes(
			attribute.String("type", crudType),
		),
	)
}

func (dbMetrics) RecordQueryFailure(ctx context.Context, crudType string, err error) {
	dbQueryErrorCount.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("type", crudType),
			attribute.String("error", err.Error()),
		),
	)
}

type RPCMetrics interface {
	RecordRequest(ctx context.Context, method string, status int)
	RecordLatency(ctx context.Context, method string, latency time.Duration)
}

type rpcMetrics struct{}

// RecordRequest logs an HTTP request count
func (rpcMetrics) RecordRequest(ctx context.Context, method string, status int) {
	requests.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.Int("status", status),
		),
	)
}

// RecordLatency logs an HTTP request latency
func (rpcMetrics) RecordLatency(ctx context.Context, method string, latency time.Duration) {
	latencyHist.Record(ctx, 1000*latency.Seconds(),
		metric.WithAttributes(attribute.String("method", method)),
	)
}
