package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"sap_segmentation/internal/config"
	"sap_segmentation/internal/logger"
	"sap_segmentation/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.InitLogger("log/segmentation_import.log"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.LogFile.Close()

	logger.CleanupOldLogs("log", cfg.LogCleanupMaxAge)

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	importService := service.NewImportService(cfg, db)
	if err := importService.ImportData(); err != nil {
		log.Fatalf("Failed to import data: %v", err)
	}
}
