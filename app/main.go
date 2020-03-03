package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/r2k1/sql-alert-manager/app/config"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("FATAL: CONFIG_PATH environment variable is not specified")
	}
	alerts, err := config.LoadAlerts(configPath)
	if err != nil {
		log.Fatalf("FATAL: error during loading configuration file %s: %s", configPath, err)
	}
	if len(alerts) == 0 {
		log.Fatal("FATAL: alerts are not defined, exiting")
	}
	log.Printf("INFO: prepared %v alerts", len(alerts))
	for i := range alerts {
		go alerts[i].Worker()
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	// TODO: close DB connections
	log.Print("INFO: Received exit signal, exiting")
}
