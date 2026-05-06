package di

import (
	"github.com/TienHien1201/go-mmo-affiliate/internal/config"
	"github.com/TienHien1201/go-mmo-affiliate/internal/server/http/handler"

	xhttp "github.com/TienHien1201/go-mmo-affiliate/pkg/http"
	xlogger "github.com/TienHien1201/go-mmo-affiliate/pkg/logger"
)

type AppContainer struct {
	HTTPHandler xhttp.Handler
}

func NewAppContainer(cfg *config.Config, logger *xlogger.Logger) (*AppContainer, func(), error) {
	// Create dependencies
	postpresClient, err := cfg.InitPostgresDB()
	if err != nil {
		return nil, nil, err
	}

		httpHandler := handler.NewHTTPHandler(logger)

	cleanup := func() {
		if err := postpresClient.Close(); err != nil {
			logger.Error("Close Postgres Client faild", xlogger.Error(err))
		}
	}

	return  &AppContainer{
		HTTPHandler: httpHandler,

	}, cleanup, nil
}