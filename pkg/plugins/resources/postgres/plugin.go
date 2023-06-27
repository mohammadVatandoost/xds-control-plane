package postgres

import (
	"errors"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/plugins/resources/postgres"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core"
	core_plugins "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/plugins"
	core_store "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/store"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/runtime/component"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/events"
	postgres_events "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/resources/postgres/events"
)

var _ core_plugins.ResourceStorePlugin = &plugin{}

type plugin struct{}

func init() {
	core_plugins.Register(core_plugins.Postgres, &plugin{})
}

func (p *plugin) NewResourceStore(pc core_plugins.PluginContext, config core_plugins.PluginConfig) (core_store.ResourceStore, error) {
	cfg, ok := config.(*postgres.PostgresStoreConfig)
	if !ok {
		return nil, errors.New("invalid type of the config. Passed config should be a PostgresStoreConfig")
	}
	migrated, err := isDbMigrated(*cfg)
	if err != nil {
		return nil, err
	}
	if !migrated {
		return nil, errors.New(`database is not migrated. Run "kuma-cp migrate up" to update database to the newest schema`)
	}
	switch cfg.DriverName {
	case postgres.DriverNamePgx:
		return NewPgxStore(pc.Metrics(), *cfg, pc.PgxConfigCustomizationFn())
	case postgres.DriverNamePq:
		return NewPqStore(pc.Metrics(), *cfg)
	default:
		return nil, errors.New("unknown driver name " + cfg.DriverName)
	}
}

func (p *plugin) Migrate(pc core_plugins.PluginContext, config core_plugins.PluginConfig) (core_plugins.DbVersion, error) {
	cfg, ok := config.(*postgres.PostgresStoreConfig)
	if !ok {
		return 0, errors.New("invalid type of the config. Passed config should be a PostgresStoreConfig")
	}
	return MigrateDb(*cfg)
}

func (p *plugin) EventListener(pc core_plugins.PluginContext, out events.Emitter) error {
	postgresListener := postgres_events.NewListener(*pc.Config().Store.Postgres, out)
	return pc.ComponentManager().Add(component.NewResilientComponent(core.Log.WithName("postgres-event-listener-component"), postgresListener))
}
