package bulk

import (
	rawSql "database/sql"
	"fmt"
	"strings"
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
	batch               []sql.Model
	batchSizeLimit      int
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
		batch:               make([]sql.Model, 0, cfg.SQL.BatchSizeLimit),
		dcpCheckpointCommit: dcpCheckpointCommit,
		batchSizeLimit:      cfg.SQL.BatchSizeLimit,
		batchTickerDuration: cfg.SQL.BatchTickerDuration,
		batchTicker:         time.NewTicker(cfg.SQL.BatchTickerDuration),
	}
	return &b, nil
}

func (b *Bulk) StartBulk() {
	for range b.batchTicker.C {
		b.flushBatch()
	}
}

func (b *Bulk) Close() {
	b.batchTicker.Stop()

	b.flushBatch()
}

func (b *Bulk) AddActions(ctx *models.ListenerContext, actions []sql.Model) {
	b.flushLock.Lock()
	if b.isDcpRebalancing {
		logger.Log.Warn("could not add new message to batch while rebalancing")
		b.flushLock.Unlock()
		return
	}

	b.batch = append(b.batch, actions...)

	ctx.Ack()

	b.flushLock.Unlock()

	if len(b.batch) >= b.batchSizeLimit {
		b.flushBatch()
	}
}

func (b *Bulk) flushBatch() {
	b.flushLock.Lock()
	defer b.flushLock.Unlock()
	if b.isDcpRebalancing {
		return
	}

	queries := prepareBulkSQLQueries(b.batch)

	for _, query := range queries {
		result, err := b.sqlClient.Exec(query)
		if err != nil {
			panic(err)
		} else {
			affected, err := result.RowsAffected()
			if err != nil {
				panic(err)
			} else {
				logger.Log.Info("affected = %v", affected)
			}
		}
	}
	b.batchTicker.Reset(b.batchTickerDuration)
	b.batch = b.batch[:0]
	b.dcpCheckpointCommit()
}

func prepareBulkSQLQueries(batch []sql.Model) []string {
	queries := make([]string, 0)

	insertQueries := make(map[string][]string)
	for _, model := range batch {
		query := model.Convert()

		if !strings.HasPrefix(query, "INSERT INTO") {
			queries = append(queries, query)
			continue
		}

		prepareInsertQueries(query, insertQueries)
	}

	for prefix, data := range insertQueries {
		queries = append(queries, fmt.Sprintf("%s %s", prefix, strings.Join(data, ",")))
	}

	return queries
}

func prepareInsertQueries(query string, insertQueries map[string][]string) {
	s := strings.Split(query, "VALUES")
	prefix := s[0] + " VALUES"
	data := s[1]

	exist, ok := insertQueries[prefix]
	if !ok {
		insertQueries[prefix] = []string{data}
	} else {
		insertQueries[prefix] = append(exist, data)
	}
}

func (b *Bulk) PrepareStartRebalancing() {
	b.flushLock.Lock()
	defer b.flushLock.Unlock()

	b.isDcpRebalancing = true
	b.batch = b.batch[:0]
}

func (b *Bulk) PrepareEndRebalancing() {
	b.flushLock.Lock()
	defer b.flushLock.Unlock()

	b.isDcpRebalancing = false
}
