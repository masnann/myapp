package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_DeleteUser(t *testing.T) {
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error creating mock database:", err)
	}
	defer db.Close()
	repo := NewUserRepository(db)

	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users WHERE id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteUser(1)
		assert.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		// Kembalikan kesalahan dari database
		mock.ExpectExec("DELETE FROM users WHERE id = ?").WithArgs(3).WillReturnError(errors.New("database error"))

		err := repo.DeleteUser(3)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})

	// Pastikan semua ekspektasi telah terpenuhi
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
