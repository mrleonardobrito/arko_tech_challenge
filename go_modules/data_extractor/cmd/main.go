package main

import (
	"log"
	"os"
	"path/filepath"
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

func init() {
	companyZipUrl := getEnv("COMPANY_ZIP_URL", "https://arquivos.receitafederal.gov.br/dados/cnpj/dados_abertos_cnpj/2025-05/Empresas0.zip")
	companyStoragePath := getEnv("COMPANY_STORAGE_PATH", "go_modules/data_extractor/data")
	companyZipPath := filepath.Join(companyStoragePath, "companies.zip")
}

func main() {
	log.Println("Starting data extraction")

	if err := os.MkdirAll(companyStoragePath, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	zipFile, err := os.Create(filepath.Join(companyStoragePath, "companies.zip"))
	if err != nil {
		log.Fatalf("Failed to create zip file: %v", err)
	}
}
