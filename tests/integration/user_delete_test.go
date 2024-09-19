package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/username/myapp/handler"
	"github.com/username/myapp/repository"
	"github.com/username/myapp/service"

	"github.com/pressly/goose/v3"
)

var TestDB *sql.DB

func TestMain(m *testing.M) {
	// Konfigurasi DSN untuk database testing
	dsn := "postgres://postgres:mkpmobile2024@localhost:5432/myapp_test?sslmode=disable"

	var err error
	TestDB, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}
	defer TestDB.Close()

	// Jalankan migrasi
	if err := runMigrations(TestDB); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	// Jalankan tes
	code := m.Run()

	// Bersihkan database setelah tes
	cleanup(TestDB)
	cleanupDatabase(TestDB)

	os.Exit(code)
}

func runMigrations(db *sql.DB) error {
	migrationDir := filepath.Join("..", "..", "migrations")

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, migrationDir)
	if err != nil {
		return err
	}

	return nil
}

func cleanup(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS users;")
	if err != nil {
		fmt.Printf("Failed to clean up test database: %v\n", err)
	}
}

func cleanupDatabase(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS goose_db_version")
	if err != nil {
		log.Fatalf("Failed to clean up goose_db_version: %v", err)
	}
}

func insertTestUser(db *sql.DB, username, email, password string) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		username, email, password,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func setupEcho(db *sql.DB) *echo.Echo {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()
	e.DELETE("/users/:id", userHandler.DeleteUser)
	return e
}

func TestDeleteUser_Success(t *testing.T) {
	e := setupEcho(TestDB)

	// Insert test user
	userID, err := insertTestUser(TestDB, "johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	// Buat HTTP request DELETE /users/:id
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", userID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(userID))

	// Panggil handler
	handler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(TestDB)))
	err = handler.DeleteUser(c)
	assert.NoError(t, err)

	// Cek respons
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verifikasi bahwa user telah dihapus dari database
	var count int
	err = TestDB.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", userID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDeleteUser_NotFound(t *testing.T) {
	e := setupEcho(TestDB)

	// ID yang tidak ada di database
	nonExistentID := 9999

	// Buat HTTP request DELETE /users/:id
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", nonExistentID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(nonExistentID))

	// Panggil handler
	handler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(TestDB)))
	err := handler.DeleteUser(c)
	assert.NoError(t, err)

	// Cek respons
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Cek pesan error
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "user not found", resp["error"])
}

func TestDeleteUser_InvalidID(t *testing.T) {
	e := setupEcho(TestDB)

	// ID tidak valid
	invalidID := "abc"

	// Buat HTTP request DELETE /users/:id
	req := httptest.NewRequest(http.MethodDelete, "/users/"+invalidID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	// Panggil handler
	handler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(TestDB)))
	err := handler.DeleteUser(c)
	assert.NoError(t, err)

	// Cek respons
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Cek pesan error
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid user id", resp["error"])
}
