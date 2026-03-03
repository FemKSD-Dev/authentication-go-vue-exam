package postgres

import (
	"context"
	"testing"

	"authentication-project-exam/internal/core/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	// Create SQLite-compatible schema (UserRecord has Postgres-specific types)
	err = db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME,
		updated_at DATETIME,
		delete_at DATETIME
	)`).Error
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func TestNewUsersRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUsersRepository(db)
	if repo == nil {
		t.Fatal("NewUsersRepository returned nil")
	}
}

func TestUsersRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUsersRepository(db)
	ctx := context.Background()

	user := &model.User{
		ID:       "user-1",
		Username: "testuser",
		Password: "hashed",
	}

	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
}

func TestUsersRepository_FindByUsername_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUsersRepository(db)
	ctx := context.Background()

	user := &model.User{
		ID:       "user-1",
		Username: "findme",
		Password: "hashed",
	}
	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	found, err := repo.FindByUsername(ctx, "findme")
	if err != nil {
		t.Fatalf("FindByUsername failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Username != "findme" {
		t.Errorf("expected username findme, got %s", found.Username)
	}
	if found.Password != "hashed" {
		t.Errorf("expected password hashed, got %s", found.Password)
	}
}

func TestUsersRepository_FindByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUsersRepository(db)
	ctx := context.Background()

	found, err := repo.FindByUsername(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("FindByUsername failed: %v", err)
	}
	if found != nil {
		t.Errorf("expected nil for nonexistent user, got %v", found)
	}
}

func TestUsersRepository_FindByID_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUsersRepository(db)
	ctx := context.Background()

	user := &model.User{
		ID:       "user-123",
		Username: "byid",
		Password: "hashed",
	}
	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	found, err := repo.FindByID(ctx, "user-123")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.ID != "user-123" {
		t.Errorf("expected ID user-123, got %s", found.ID)
	}
}

func TestUsersRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUsersRepository(db)
	ctx := context.Background()

	found, err := repo.FindByID(ctx, "nonexistent-id")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found != nil {
		t.Errorf("expected nil for nonexistent id, got %v", found)
	}
}

func TestUserRecord_TableName(t *testing.T) {
	r := UserRecord{}
	if name := r.TableName(); name != "users" {
		t.Errorf("expected TableName users, got %s", name)
	}
}
