package bulk

import (
	"database/sql"
	"github.com/Trendyol/go-dcp-sql/config"
	"github.com/Trendyol/go-dcp-sql/model"
	"github.com/Trendyol/go-dcp-sql/sql/client"
	"github.com/Trendyol/go-dcp/logger"
	"github.com/Trendyol/go-dcp/models"
	"sync"
)

type Bulk struct {
	sqlClient           *sql.DB
	batch               []model.DcpSqlItem
	flushLock           sync.Mutex
	isDcpRebalancing    bool
	dcpCheckpointCommit func()
	batchSizeLimit      int
}

func NewBulk(
	cfg *config.Connector,
	dcpCheckpointCommit func(),
) (*Bulk, error) {
	c, err := client.NewSqlClient(cfg.Sql)
	if err != nil {
		return nil, err
	}

	b := Bulk{
		sqlClient:           c,
		batch:               make([]model.DcpSqlItem, 0, cfg.Sql.BatchSizeLimit),
		dcpCheckpointCommit: dcpCheckpointCommit,
		batchSizeLimit:      cfg.Sql.BatchSizeLimit,
	}
	return &b, nil
}

func (b *Bulk) AddActions(ctx *models.ListenerContext, actions []model.DcpSqlItem) {
	b.flushLock.Lock()
	if b.isDcpRebalancing {
		logger.Log.Warn("could not add new message to batch while rebalancing")
		b.flushLock.Unlock()
		return
	}

	for _, action := range actions {
		b.batch = append(b.batch, action)
	}
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
	for _, model := range b.batch {
		result, err := b.sqlClient.Exec(model.ConvertSql())
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
	b.batch = b.batch[:0]
	b.dcpCheckpointCommit()
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
