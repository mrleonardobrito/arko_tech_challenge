package pipelines

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arko_tech_challenge/go_modules/data_extractor/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type State struct {
	ID      int    `json:"id"`
	Name    string `json:"nome"`
	Acronym string `json:"sigla"`
}

func (s State) GetID() int {
	return s.ID
}

type StateEncoder struct {
	pool *pgxpool.Pool
}

func NewStateEncoder(pool *pgxpool.Pool) internal.DBEncoder[State] {
	return &StateEncoder{pool: pool}
}

func (e *StateEncoder) Encode(ctx context.Context, v State) ([]any, error) {
	return []any{v.ID, v.Name, v.Acronym}, nil
}

func RunStatesPipeline(ctx context.Context, pool *pgxpool.Pool, apiUrl string, batchSize int) error {
	src := internal.NewAPISource(apiUrl, func(data []byte) ([]State, bool, error) {
		var states []State
		if err := json.Unmarshal(data, &states); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal states: %w", err)
		}
		return states, false, nil
	})

	batcher := internal.NewFixedSizeBatcher[State](batchSize, src.ItemCount())
	encoder := NewStateEncoder(pool)
	tableSpec := internal.TableSpec{
		Name:           "state",
		Columns:        []string{"id", "name", "acronym"},
		ConflictMode:   internal.ConflictModeUpdate,
		ConflictColumn: "id",
		UpdateColumns:  []string{"name", "acronym"},
	}
	db := internal.NewPostgresRepository(pool, tableSpec, encoder, internal.PGOptions{
		TxTimeout: 10 * time.Second,
	})

	return RunPipeline(ctx, src, batcher, db)
}

type City struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	StateID int    `json:"state_id"`
}

func (c City) GetID() int {
	return c.ID
}

type CityEncoder struct {
	pool *pgxpool.Pool
}

func NewCityEncoder(pool *pgxpool.Pool) internal.DBEncoder[City] {
	return &CityEncoder{pool: pool}
}

func (e *CityEncoder) Encode(ctx context.Context, v City) ([]any, error) {
	return []any{v.ID, v.Name, v.StateID}, nil
}

type CityResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"nome"`
	MicroRegion struct {
		ID         int    `json:"id"`
		Name       string `json:"nome"`
		Mesoregion struct {
			ID   int    `json:"id"`
			Name string `json:"nome"`
			UF   struct {
				ID      int    `json:"id"`
				Name    string `json:"nome"`
				Acronym string `json:"sigla"`
			} `json:"UF"`
		} `json:"mesorregiao"`
	} `json:"microrregiao"`
}

func RunCitiesPipeline(ctx context.Context, pool *pgxpool.Pool, apiUrl string, batchSize int) error {
	rows, err := pool.Query(ctx, "SELECT id FROM state")
	if err != nil {
		return fmt.Errorf("failed to query states: %w", err)
	}
	defer rows.Close()

	validStateIDs := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan state id: %w", err)
		}
		validStateIDs[id] = true
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating state rows: %w", err)
	}

	src := internal.NewAPISource(apiUrl, func(data []byte) ([]City, bool, error) {
		var citiesResponse []CityResponse
		if err := json.Unmarshal(data, &citiesResponse); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal cities: %w", err)
		}

		cities := make([]City, 0, len(citiesResponse))
		for _, city := range citiesResponse {
			stateID := city.MicroRegion.Mesoregion.UF.ID
			if validStateIDs[stateID] {
				cities = append(cities, City{
					ID:      city.ID,
					Name:    city.Name,
					StateID: stateID,
				})
			}
		}
		return cities, false, nil
	})

	batcher := internal.NewFixedSizeBatcher[City](batchSize, src.ItemCount())
	encoder := NewCityEncoder(pool)
	tableSpec := internal.TableSpec{
		Name:           "city",
		Columns:        []string{"id", "name", "state_id"},
		UpdateColumns:  []string{"name", "state_id"},
		ConflictMode:   internal.ConflictModeUpdate,
		ConflictColumn: "id",
	}
	db := internal.NewPostgresRepository(pool, tableSpec, encoder, internal.PGOptions{
		TxTimeout: 10 * time.Second,
	})

	return RunPipeline(ctx, src, batcher, db)
}

type District struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	CityID int    `json:"city_id"`
}

func (d District) GetID() int {
	return d.ID
}

type DistrictEncoder struct {
	pool *pgxpool.Pool
}

func NewDistrictEncoder(pool *pgxpool.Pool) internal.DBEncoder[District] {
	return &DistrictEncoder{pool: pool}
}

func (e *DistrictEncoder) Encode(ctx context.Context, v District) ([]any, error) {
	return []any{v.ID, v.Name, v.CityID}, nil
}

type DistrictResponse struct {
	ID   int    `json:"id"`
	Name string `json:"nome"`
	City struct {
		ID int `json:"id"`
	} `json:"municipio"`
}

func RunDistrictsPipeline(ctx context.Context, pool *pgxpool.Pool, apiUrl string, batchSize int) error {
	rows, err := pool.Query(ctx, "SELECT id FROM city")
	if err != nil {
		return fmt.Errorf("failed to query cities: %w", err)
	}
	defer rows.Close()

	validCityIDs := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan city id: %w", err)
		}
		validCityIDs[id] = true
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating city rows: %w", err)
	}

	src := internal.NewAPISource(apiUrl, func(data []byte) ([]District, bool, error) {
		var districtsResponse []DistrictResponse
		if err := json.Unmarshal(data, &districtsResponse); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal districts: %w", err)
		}

		districts := make([]District, 0, len(districtsResponse))
		for _, district := range districtsResponse {
			if validCityIDs[district.City.ID] {
				districts = append(districts, District{
					ID:     district.ID,
					Name:   district.Name,
					CityID: district.City.ID,
				})
			}
		}
		return districts, false, nil
	})

	batcher := internal.NewFixedSizeBatcher[District](batchSize, src.ItemCount())
	encoder := NewDistrictEncoder(pool)
	tableSpec := internal.TableSpec{
		Name:           "district",
		Columns:        []string{"id", "name", "city_id"},
		UpdateColumns:  []string{"name", "city_id"},
		ConflictMode:   internal.ConflictModeUpdate,
		ConflictColumn: "id",
	}
	db := internal.NewPostgresRepository(pool, tableSpec, encoder, internal.PGOptions{
		TxTimeout: 10 * time.Second,
	})

	return RunPipeline(ctx, src, batcher, db)
}
