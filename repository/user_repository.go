package repository

import (
    "database/sql"
    "fmt"
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

type userRepositoryImpl struct {
    DB *sql.DB
}

// NewUserRepository adalah konstruktor untuk UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepositoryImpl{
        DB: db,
    }
}

// DeleteUser menghapus pengguna berdasarkan ID
func (r *userRepositoryImpl) DeleteUser(id int) error {
    result, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}
