package pipelines

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/arko_tech_challenge/go_modules/data_extractor/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Company struct {
	CNPJ                     string  `csv:"cnpj"`
	SocialName               string  `csv:"social_name"`
	JuridicalNature          string  `csv:"juridical_nature"`
	ResponsibleQualification string  `csv:"responsible_qualification"`
	SocialCapital            float64 `csv:"social_capital"`
	CompanySize              string  `csv:"company_size"`
	FederativeEntity         string  `csv:"federative_entity"`
}

var _ internal.DBEncoder[Company] = &CompanyEncoder{}

type CompanyEncoder struct {
	pool *pgxpool.Pool
}

func NewCompanyEncoder(pool *pgxpool.Pool) internal.DBEncoder[Company] {
	return &CompanyEncoder{pool: pool}
}

func (e *CompanyEncoder) Encode(ctx context.Context, v Company) ([]any, error) {
	return []any{
		v.CNPJ,
		v.SocialName,
		v.JuridicalNature,
		v.ResponsibleQualification,
		v.SocialCapital,
		v.CompanySize,
		v.FederativeEntity,
	}, nil
}

func RunCompaniesPipeline(ctx context.Context, pool *pgxpool.Pool, downloadUrl, downloadPath, extractPath string, batchSize int) error {
	downloader := internal.NewHTTPDownloader()

	log.Println("Downloading file...")
	if err := downloader.Download(ctx, downloadUrl, downloadPath); err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}

	log.Println("Extracting file...")
	if err := downloader.Extract(ctx, downloadPath, extractPath); err != nil {
		log.Fatalf("Failed to extract file: %v", err)
	}

	files, err := os.ReadDir(extractPath)
	if err != nil {
		log.Fatalf("Failed to read extract directory: %v", err)
	}

	extractedFilePath := filepath.Join(extractPath, files[0].Name())
	csvPath := filepath.Join(extractPath, "companies.csv")
	companyColumns := []string{"cnpj", "social_name", "juridical_nature", "responsible_qualification", "social_capital", "company_size", "federative_entity"}

	itemCount, err := convertToCSV(extractedFilePath, csvPath, companyColumns, 10000000000)
	if err != nil {
		log.Fatalf("Failed to convert CSV: %v", err)
	}
	log.Printf("Converted %d records to CSV", itemCount)

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	log.Printf("CSV file opened: %v", csvPath)
	defer file.Close()

	src := internal.NewCSVSource(file, ';', func(cols []string) (Company, error) {
		rawCapital := strings.ReplaceAll(cols[4], ".", "")
		rawCapital = strings.ReplaceAll(rawCapital, ",", ".")
		capital, err := strconv.ParseFloat(rawCapital, 64)
		if err != nil {
			return Company{}, fmt.Errorf("failed to parse social capital: %w", err)
		}
		return Company{
			CNPJ:                     cols[0],
			SocialName:               cols[1],
			JuridicalNature:          cols[2],
			ResponsibleQualification: cols[3],
			SocialCapital:            capital,
			CompanySize:              cols[5],
			FederativeEntity:         cols[6],
		}, nil
	}, itemCount)
	defer src.Close()

	batcher := internal.NewFixedSizeBatcher[Company](batchSize, src.ItemCount())
	encoder := NewCompanyEncoder(pool)
	tableSpec := internal.TableSpec{
		Name:           "company",
		Columns:        []string{"cnpj", "social_name", "juridical_nature", "responsible_qualification", "social_capital", "company_size", "federative_entity"},
		ConflictMode:   internal.ConflictModeUpdate,
		ConflictColumn: "cnpj",
		UpdateColumns:  []string{"social_name", "juridical_nature", "responsible_qualification", "social_capital", "company_size", "federative_entity"},
	}
	db := internal.NewPostgresRepository(pool, tableSpec, encoder, internal.PGOptions{
		TxTimeout: 10 * time.Second,
	})
	return RunPipeline(ctx, src, batcher, db)
}

func convertToCSV(inputPath, outputPath string, headers []string, limit int) (int, error) {
	in, err := os.Open(inputPath)
	if err != nil {
		return 0, fmt.Errorf("error opening input file: %w", err)
	}
	defer in.Close()

	bufferedReader := bufio.NewReaderSize(in, 1024*1024)

	r := csv.NewReader(bufferedReader)
	r.Comma = ';'
	r.ReuseRecord = true
	r.FieldsPerRecord = len(headers)
	r.LazyQuotes = true
	r.TrimLeadingSpace = true

	out, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	bufferedWriter := bufio.NewWriterSize(out, 1024*1024)
	defer bufferedWriter.Flush()

	w := csv.NewWriter(bufferedWriter)
	w.Comma = ';'
	defer w.Flush()

	if len(headers) > 0 {
		if err := w.Write(headers); err != nil {
			return 0, fmt.Errorf("error writing headers: %w", err)
		}
	}

	recordCount := 0
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("error reading CSV: %w", err)
		}
		record := make([]string, len(headers))
		for i, val := range rec {
			val = strings.TrimSpace(val)
			if !utf8.ValidString(val) {
				val = strings.Map(func(r rune) rune {
					if r == utf8.RuneError {
						return -1
					}
					return r
				}, val)
			}
			record[i] = val
		}
		if err := w.Write(record); err != nil {
			return 0, fmt.Errorf("error writing CSV: %w", err)
		}

		recordCount++
		if limit > 0 && recordCount >= limit {
			break
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return 0, fmt.Errorf("error finalizing write: %w", err)
	}

	return recordCount, nil
}
