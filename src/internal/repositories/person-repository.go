package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/database"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/entities"
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

func (c *PersonRepository) CreatePerson(ctx context.Context, personRequest entities.PersonRequestDTO) (entities.PersonResponseDTO, error) {
	return entities.PersonResponseDTO{
		Id:        "fef06178-3685-4e9d-bcc1-4c04ad8132fb",
		Surname:   personRequest.Surname,
		Name:      personRequest.Name,
		Birthdate: personRequest.Birthdate,
		Stack:     personRequest.Stack,
	}, nil
}
