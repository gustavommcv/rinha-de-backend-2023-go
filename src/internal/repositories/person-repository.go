package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
		Id:        personID.String(),
		Surname:   personRequest.Surname,
		Name:      personRequest.Name,
		Birthdate: personRequest.Birthdate,
		Stack:     personRequest.Stack,
	}, nil
}

func (p *PersonRepository) FindById(ctx context.Context, id uuid.UUID) (*entities.PersonResponseDTO, error) {
	rows, err := p.pool.Query(
		ctx, `
			SELECT people.*, languages.name AS language_name FROM PEOPLE 
			LEFT JOIN stack ON people.person_id = stack.person_id
			LEFT JOIN languages ON stack.language_id = languages.language_id
			WHERE people.person_id = $1 
		`, id)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("query row failed: %w", err)
	}
	defer rows.Close()

	type person_row struct {
		person_id string
		name      string
		surname   string
		birthdate time.Time
		languages []string
	}

	pr := person_row{}
	for rows.Next() {
		var l string
		rows.Scan(&pr.person_id, &pr.name, &pr.surname, &pr.birthdate, &l)
		if l != "" {
			pr.languages = append(pr.languages, l)
		}
	}

	return &entities.PersonResponseDTO{
		Id:        pr.person_id,
		Surname:   pr.surname,
		Name:      pr.name,
		Birthdate: pr.birthdate.Format("2006-01-02"),
		Stack:     pr.languages,
	}, nil
}

func (p *PersonRepository) Search(ctx context.Context, searchTerm string) ([]entities.PersonResponseDTO, error) {
	var sb strings.Builder

	sb.WriteString("%")
	sb.WriteString(searchTerm)
	sb.WriteString("%")

	rows, err := p.pool.Query(
		ctx, `
		SELECT p.person_id FROM people p
LEFT JOIN stack s ON p.person_id = s.person_id
LEFT JOIN languages l ON l.language_id = s.language_id
WHERE LOWER(p.name) LIKE LOWER($1)
OR LOWER(p.surname) LIKE LOWER($1)
OR LOWER(l.name) LIKE LOWER($1)
LIMIT 50;
		`, sb.String())

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("query row failed: %w", err)
	}
	defer rows.Close()

	type person_row struct {
		person_id string
		name      string
		surname   string
		birthdate time.Time
		languages []string
	}

	people := []entities.PersonResponseDTO{}
	ids := []string{}
	for rows.Next() {
		var currentId string
		rows.Scan(&currentId)
		ids = append(ids, currentId)
	}

	for _, id := range ids {
		person, err := p.FindById(ctx, uuid.MustParse(id))
		if err != nil {
			return nil, err
		}

		people = append(people, *person)
	}

	return people, nil
}
