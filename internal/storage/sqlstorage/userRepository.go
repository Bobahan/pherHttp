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

	u.Id = int(id)

	return nil
}

func (r *UserRepository) FindById(id string) (*model.User, error) {
	user, err := r.getUserBy("id", id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user, err := r.getUserBy("email", email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) getUserBy(param string, value string) (*model.User, error) {
	user := &model.User{}

	query := fmt.Sprintf("SELECT id, email, pass FROM users WHERE %s = ?", param)

	row := r.storage.db.QueryRow(query, value)

	if err := row.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}
