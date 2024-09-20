package repository

import (
	"database/sql"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	DeleteUser(id int) error
	// Metode lain seperti GetUserByID, CreateUser, dsb.
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

// NewUserRepository adalah konstruktor untuk UserRepository
func NewUserRepository(db *sql.DB) UserRepositoryImpl {
	return UserRepositoryImpl{
		DB: db,
	}
}

// DeleteUser menghapus pengguna berdasarkan ID
func (r UserRepositoryImpl) DeleteUser(id int) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
