//go:build wireinject
// +build wireinject

package wire

import (
	"nunu/internal/handler"
	"nunu/internal/repository"
	"nunu/internal/server"
	"nunu/internal/service"
	"nunu/pkg/app"
	"nunu/pkg/embedding"
	"nunu/pkg/log"
	"nunu/pkg/server/http"
	"nunu/pkg/slackclient"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewQueryDB,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewMemoryRepository,
	repository.NewKnowledgeRepository,
	repository.NewConversationRepository,
	repository.NewPostgresQueryEngine,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewMemoryService,
	service.NewKnowledgeService,
	provideQueryService,
	provideCSVService,
	service.NewAgentService,
	service.NewSlackService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewSlackHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewSlackServer,
)

var pkgSet = wire.NewSet(
	embedding.NewClient,
	slackclient.NewClient,
)

func provideQueryService(engine repository.QueryEngine, logger *log.Logger, conf *viper.Viper) service.QueryService {
	return service.NewQueryService(engine, logger, conf.GetInt("agent.max_query_rows"))
}

func provideCSVService(engine repository.QueryEngine, logger *log.Logger, conf *viper.Viper) service.CSVService {
	return service.NewCSVService(
		engine,
		logger,
		conf.GetString("csv.storage_path"),
		conf.GetString("csv.base_url"),
	)
}

// build App
func newApp(
	httpServer *http.Server,
	slackServer *server.SlackServer,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, slackServer),
		app.WithName("nunu-agent"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		pkgSet,
		newApp,
	))
}
