package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/arko_tech_challenge/go_modules/data_extractor/internal/pipelines"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue != "" {
			return defaultValue
		}
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func newPGPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	cfg.MinConns = 0
	cfg.MaxConns = int32(max(4, runtime.NumCPU()*4))
	cfg.MaxConnLifetime = 10 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second

	cfg.ConnConfig.RuntimeParams = map[string]string{
		"client_encoding":             "UTF8",
		"standard_conforming_strings": "on",
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping pool: %w", err)
	}
	return pool, nil
}

func main() {
	ctx := context.Background()
	log.Println("Starting data extraction")

	companyZipUrl := getEnv("COMPANY_ZIP_URL", "https://arquivos.receitafederal.gov.br/dados/cnpj/dados_abertos_cnpj/2025-05/Empresas0.zip")
	companyStoragePath := getEnv("COMPANY_STORAGE_PATH", "data")
	companyZipPath := filepath.Join(companyStoragePath, "companies.zip")
	extractPath := filepath.Join(companyStoragePath, "extracted")
	databaseUrl := getEnv("DATABASE_URL", "postgresql://admin:1234567@localhost:5437/arko_tech_challenge")
	locationUrl := getEnv("LOCATION_API_URL", "https://servicodados.ibge.gov.br/api/v1/localidades")

	stateLocationUrl := fmt.Sprintf("%s/estados", locationUrl)
	cityLocationUrl := fmt.Sprintf("%s/municipios", locationUrl)
	districtLocationUrl := fmt.Sprintf("%s/distritos", locationUrl)

	pool, err := newPGPool(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	if err := pipelines.RunStatesPipeline(ctx, pool, stateLocationUrl, 300); err != nil {
		log.Fatalf("Failed to run states pipeline: %v", err)
	}
	log.Println("States pipeline completed successfully")

	if err := pipelines.RunCitiesPipeline(ctx, pool, cityLocationUrl, 10); err != nil {
		log.Fatalf("Failed to run cities pipeline: %v", err)
	}
	log.Println("Cities pipeline completed successfully")

	if err := pipelines.RunDistrictsPipeline(ctx, pool, districtLocationUrl, 30); err != nil {
		log.Fatalf("Failed to run districts pipeline: %v", err)
	}
	log.Println("Districts pipeline completed successfully")
	pipelines.RunCompaniesPipeline(ctx, pool, companyZipUrl, companyZipPath, extractPath, 5000)
}
