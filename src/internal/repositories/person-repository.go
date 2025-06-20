package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/database"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/entities"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
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
		return 0, fmt.Errorf("query row failed: %w", err)
	}
	return count, nil
}

func (p *PersonRepository) CreatePerson(ctx context.Context, personRequest entities.PersonRequestDTO) (*entities.PersonResponseDTO, error) {
	// Insert person
	var personID uuid.UUID
	err := p.pool.QueryRow(
		ctx, `
			INSERT INTO people (surname, name, birthdate) 
			VALUES ($1, $2, $3) 
			RETURNING person_id`,
		personRequest.Surname, personRequest.Name, personRequest.Birthdate).Scan(&personID)

	if err != nil {
		return nil, fmt.Errorf("failed to insert person: %w", err)
	}

	// Insert languages
	for _, language := range personRequest.Stack {
		var languageID uuid.UUID
		err := p.pool.QueryRow(
			ctx, `
				INSERT INTO languages (name)
				VALUES ($1) 
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
				RETURNING language_id`,
			language).Scan(&languageID)

		if err != nil {
			return nil, fmt.Errorf("failed to get/insert language %s: %w", language, err)
		}

		// Insert stack
		_, err = p.pool.Exec(
			ctx, `
				INSERT INTO stack (person_id, language_id)
				VALUES ($1, $2)
				ON CONFLICT (person_id, language_id) DO NOTHING`,
			personID,
			languageID,
		)

		if err != nil {
			return nil, fmt.Errorf("failed insert stack %s: %w", language, err)
		}
	}

	return &entities.PersonResponseDTO{
		Id:        personID.UUID.String(),
		Surname:   personRequest.Surname,
		Name:      personRequest.Name,
		Birthdate: personRequest.Birthdate,
		Stack:     personRequest.Stack,
	}, nil
}
