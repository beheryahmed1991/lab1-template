package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bmstu-rsoi/lab1-template.git/internal/model"
)

type PersonRepository struct {
	db *pgxpool.Pool
}

func NewPersonRepository(db *pgxpool.Pool) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, person model.PersonRequest) (int64, error) {
	var id int64

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO persons (name, age, address, work)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		person.Name,
		person.Age,
		person.Address,
		person.Work,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PersonRepository) GetByID(ctx context.Context, id int64) (*model.Person, error) {
	var person model.Person

	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, age, address, work
		 FROM persons
		 WHERE id = $1`,
		id,
	).Scan(
		&person.ID,
		&person.Name,
		&person.Age,
		&person.Address,
		&person.Work,
	)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (r *PersonRepository) GetAll(ctx context.Context) ([]model.Person, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, name, age, address, work
		 FROM persons
		 ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := make([]model.Person, 0)
	for rows.Next() {
		var person model.Person

		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Age,
			&person.Address,
			&person.Work,
		)
		if err != nil {
			return nil, err
		}

		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (r *PersonRepository) Update(ctx context.Context, id int64, person model.PersonRequest) (*model.Person, error) {
	var updated model.Person
	err := r.db.QueryRow(
		ctx,
		`UPDATE persons
		 SET name = $1, age = $2, address = $3, work = $4
		 WHERE id = $5
		 RETURNING id, name, age, address, work`,
		person.Name,
		person.Age,
		person.Address,
		person.Work,
		id,
	).Scan(&updated.ID, &updated.Name, &updated.Age, &updated.Address, &updated.Work)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *PersonRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM persons
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
