package bulk

import (
	rawSql "database/sql"
	"sync"
	"time"

	"github.com/Trendyol/go-dcp-sql/config"
	"github.com/Trendyol/go-dcp-sql/sql"
	"github.com/Trendyol/go-dcp-sql/sql/client"
	"github.com/Trendyol/go-dcp/logger"
	"github.com/Trendyol/go-dcp/models"
)

type Bulk struct {
	sqlClient           *rawSql.DB
	dcpCheckpointCommit func()
	batchTicker         *time.Ticker
	metric              *Metric
	batchTickerDuration time.Duration
	flushLock           sync.Mutex
	isDcpRebalancing    bool
}

func NewBulk(
	cfg *config.Connector,
	dcpCheckpointCommit func(),
) (*Bulk, error) {
	c, err := client.NewSQLClient(cfg.SQL)
	if err != nil {
		return nil, err
	}

	b := Bulk{
		sqlClient:           c,
		dcpCheckpointCommit: dcpCheckpointCommit,
		batchTickerDuration: cfg.SQL.BatchTickerDuration,
		batchTicker:         time.NewTicker(cfg.SQL.BatchTickerDuration),
		metric:              &Metric{},
	}
	return &b, nil
}

type Metric struct {
	ProcessLatencyMs            int64
	BulkRequestProcessLatencyMs int64
}

func (b *Bulk) GetMetric() *Metric {
	return b.metric
}

func (b *Bulk) StartBulk() {
	for range b.batchTicker.C {
		b.dcpCheckpointCommit()
	}
}

func (b *Bulk) Close() {
	b.batchTicker.Stop()
}

func (b *Bulk) AddActions(ctx *models.ListenerContext, eventTime time.Time, actions []sql.Model) {
	b.flushLock.Lock()
	if b.isDcpRebalancing {
		logger.Log.Warn("could not add new message to batch while rebalancing")
		b.flushLock.Unlock()
		return
	}
	b.flushLock.Unlock()
	b.metric.ProcessLatencyMs = time.Since(eventTime).Milliseconds()
	b.flush(ctx, actions)
}

func (b *Bulk) flush(ctx *models.ListenerContext, models []sql.Model) {
	b.flushLock.Lock()
	defer b.flushLock.Unlock()
	if b.isDcpRebalancing {
		return
	}

	startedTime := time.Now()
	for _, model := range models {
		query := model.Convert()
		result, err := b.sqlClient.Exec(query.Query, query.Args...)
		if err != nil {
			logger.Log.Error("error while sql exec, err: %v", err)
			panic(err)
		} else {
			_, err = result.RowsAffected()
			if err != nil {
				logger.Log.Error("error while rows affected, err: %v", err)
				panic(err)
			}
		}
	}
	b.metric.BulkRequestProcessLatencyMs = time.Since(startedTime).Milliseconds()
	ctx.Ack()
}

func (b *Bulk) PrepareStartRebalancing() {
	b.flushLock.Lock()
	defer b.flushLock.Unlock()

	b.isDcpRebalancing = true
}

func (b *Bulk) PrepareEndRebalancing() {
	b.flushLock.Lock()
	defer b.flushLock.Unlock()

	b.isDcpRebalancing = false
}
