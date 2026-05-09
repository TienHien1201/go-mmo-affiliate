package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TienHien1201/go-mmo-affiliate/internal/config"
)

var (
	Name    = "MMO AFFILIATE APP"
	Version = "1.0.0"
	Env     = "dev"
)

func init() {
	flag.StringVar(&Name, "name", Name, "Name")
	flag.StringVar(&Version, "version", Version, "Version")
	flag.StringVar(&Env, "env", Env, "Enviroment")
}
func main() {
	flag.Parse()

	// Load configuration for the specified environment
	cfg, err := config.LoadConfig(config.Enviroment(Env), "./config")
	if err != nil {
		log.Panicf("Failed to load configuration: %v", err)
	}

	// Override application name, version and environment
	cfg.App.Name = Name
	cfg.App.Version = Version
	cfg.App.Env = Env

	log.Printf("\033[1;36mStarting\033[0m \033[1;33m%s\033[0m \033[1;32mv%s\033[0m (\033[1;35m%s\033[0m)", Name, Version, Env)

	// Initialize application config
	app, cleanup, err := NewApp(cfg)
	if err != nil {
		log.Panicf("Failed ti initialize application: %v", err)
	}
	defer cleanup()

	// Start the application
	if err := app.Start(); err != nil {
		log.Panicf("Application err: %v", err)
	}

	// Wait for interrupt sinal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop the application
	if err := app.Stop(ctx); err != nil {
		log.Panicf("Failed to stop application: %v", err)
	}

}
