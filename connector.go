package dcpsql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Trendyol/go-dcp"
	"github.com/Trendyol/go-dcp-sql/config"
	"github.com/Trendyol/go-dcp-sql/couchbase"
	"github.com/Trendyol/go-dcp/logger"
	"github.com/Trendyol/go-dcp/models"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Connector interface {
	Start()
	Close()
}

type connector struct {
	dcp    dcp.Dcp
	mapper Mapper
	config *config.Connector
	db     *sql.DB
}

func (c *connector) Start() {
	go func() {
		<-c.dcp.WaitUntilReady()
	}()
	c.dcp.Start()
}

func (c *connector) Close() {
	c.dcp.Close()
}

func (c *connector) listener(ctx *models.ListenerContext) {
	var e couchbase.Event
	switch event := ctx.Event.(type) {
	case models.DcpMutation:
		e = couchbase.NewMutateEvent(event.Key, event.Value, event.CollectionName, event.EventTime, event.Cas, event.VbID)
	case models.DcpExpiration:
		e = couchbase.NewExpireEvent(event.Key, nil, event.CollectionName, event.EventTime, event.Cas, event.VbID)
	case models.DcpDeletion:
		e = couchbase.NewDeleteEvent(event.Key, nil, event.CollectionName, event.EventTime, event.Cas, event.VbID)
	default:
		return
	}

	actions := c.mapper(e)

	if len(actions) == 0 {
		ctx.Ack()
		return
	}

	for _, action := range actions {
		var query = action.ConvertSql()
		result, err := c.db.Exec(query)
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
}

type ConnectorBuilder struct {
	mapper Mapper
	config any
}

func newConnectorConfigFromPath(path string) (*config.Connector, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c config.Connector
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func newConfig(cf any) (*config.Connector, error) {
	switch v := cf.(type) {
	case *config.Connector:
		return v, nil
	case config.Connector:
		return &v, nil
	case string:
		return newConnectorConfigFromPath(v)
	default:
		return nil, errors.New("invalid config")
	}
}

func newConnector(cf any, mapper Mapper) (Connector, error) {
	cfg, err := newConfig(cf)
	if err != nil {
		return nil, err
	}
	cfg.ApplyDefaults()

	connector := &connector{
		mapper: mapper,
		config: cfg,
	}

	dcp, err := dcp.NewDcp(&cfg.Dcp, connector.listener)
	if err != nil {
		logger.Log.Error("Dcp error: %v", err)
		return nil, err
	}

	dcpConfig := dcp.GetConfig()
	dcpConfig.Checkpoint.Type = "manual"

	connector.dcp = dcp

	var driverExist = false
	for _, driver := range sql.Drivers() {
		if driver == cfg.Sql.DriverName {
			driverExist = true
			break
		}
	}

	if !driverExist {
		panic(fmt.Errorf("driver: %s not found", cfg.Sql.DriverName))
	}

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Sql.Host, cfg.Sql.Port, cfg.Sql.User, cfg.Sql.Password, cfg.Sql.DbName, cfg.Sql.SslMode,
	)
	connector.db, err = sql.Open(cfg.Sql.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return connector, nil
}

func NewConnectorBuilder(config any) ConnectorBuilder {
	return ConnectorBuilder{
		config: config,
		mapper: DefaultMapper,
	}
}

func (c ConnectorBuilder) SetMapper(mapper Mapper) ConnectorBuilder {
	c.mapper = mapper
	return c
}

func (c ConnectorBuilder) Build() (Connector, error) {
	return newConnector(c.config, c.mapper)
}

func (c ConnectorBuilder) SetLogger(l slog.Logger) ConnectorBuilder {
	//TODO slog should be the logger
	return c
}
