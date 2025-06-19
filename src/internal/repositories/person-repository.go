package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/database"
)

type Person struct {
	apelido    string
	nome       string
	nascimento time.Time
}

type PersonRepository struct {
	pool *database.Pool
}

func NewUserRepository(pool *database.Pool) *PersonRepository {
	return &PersonRepository{pool: pool}
}

func (p *PersonRepository) GetPeopleCount(ctx context.Context) (int, error) {
	var count int
	err := p.pool.QueryRow(ctx, "SELECT COUNT(person_id) FROM PEOPLE").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query row failed: %v", err)
	}
	return count, nil
}

func (c *PersonRepository) CreatePerson(ctx context.Context) {
	// d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
	// year, month, day := d.Date()
}
