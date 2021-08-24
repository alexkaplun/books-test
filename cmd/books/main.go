package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexkaplun/books-test/config"
	"github.com/alexkaplun/books-test/service/server"
	"github.com/alexkaplun/books-test/storage"
	"github.com/spf13/pflag"
)

var (
	configPath string
)

func init() {
	pflag.StringVar(&configPath, "config", "", "config.toml")
}

func main() {
	pflag.Parse()

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load service config: %v", err)
	}

	storage, err := storage.NewPostgres(
		storage.Params{
			ConnString: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName),
		},
	)
	if err != nil {
		log.Fatalf("failed to initiate storage: %v", err)
	}

	handler := server.NewHandler(server.HandlerParams{Storage: storage})

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: server.NewRouter(
			server.RouterParams{
				Handler: handler,
			},
		),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("starting http server on port %s...\n", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen and serve error: %s\n", err)
			quit <- os.Kill
			return
		}
		fmt.Println("closing http server...")
	}()

	<-quit
	fmt.Println("shutting down...\n")
}

func loadConfig(configFilePath string) (*config.Config, error) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, err
	}
	return config.ParseConfig(configFilePath)
}
