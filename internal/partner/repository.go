package partner

import (
	"adex/internal/domain"
	"encoding/json"
	"fmt"
	"os"
)

type Repository struct {
	partners []domain.Partner
}

func (r *Repository) List() []domain.Partner {
	return r.partners
}

func NewRepository(path string) (*Repository, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error read file %w", err)
	}

	var partners []domain.Partner
	if err := json.Unmarshal(data, &partners); err != nil {
		return nil, fmt.Errorf("error json: %w", err)
	}

	return &Repository{partners: partners}, nil
}
