package main

import (
	"github.com/TienHien1201/go-mmo-affiliate/internal/config"
	"github.com/TienHien1201/go-mmo-affiliate/internal/di"
	xhttp "github.com/TienHien1201/go-mmo-affiliate/pkg/http"
	xlogger "github.com/TienHien1201/go-mmo-affiliate/pkg/logger"
	xserver "github.com/TienHien1201/go-mmo-affiliate/pkg/server"
)

type App struct {
	logger *xlogger.Logger
	server []xserver.Server
}

func NewApp(cfg *config.Config) (*App, func(), error){
	// Initialize logger
	logger, err := initLogger(cfg)
	if err != nil {
		return nil, nil, err
	}

	// Initialize dependencies 
	container, cleanup , err := di.NewAppContainer(cfg, logger)
	if err != nil {
		return nil, nil , err
	} 

	// Initialize HTTP server
	httpServer := xhttp.NewHTTPServer(logger, cfg.Server.HTTP.Host, cfg.Server.HTTP.Port, container.HTTPHandler)	

	return &App{
		logger: logger,
		server: []xserver.Server{httpServer},
	}, cleanup, nil
}


func initLogger(cfg *config.Config) (*xlogger.Logger, error) {
	logCfg := &xlogger.Config{
		Level:      cfg.Logger.Level,
		Format:     cfg.Logger.Format,
		Output:     cfg.Logger.Output,
		TimeFormat: cfg.Logger.TimeFormat,
}

logger, err := xlogger.New(logCfg)
if err != nil {
	return nil, err
}
return logger, nil
}