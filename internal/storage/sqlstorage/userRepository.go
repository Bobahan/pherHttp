package sqlstorage

import (
	"database/sql"
	"fmt"
	"gopher/internal/model"
	"gopher/internal/storage"
)

type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	res, err := r.storage.db.Exec(`INSERT INTO users (email, pass) VALUES (?, ?)`, u.Email, u.Password)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	fmt.Println("id", id)

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}

	row := r.storage.db.QueryRow(`SELECT id, email, pass FROM users WHERE email = ?`, email)

	if err := row.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}
