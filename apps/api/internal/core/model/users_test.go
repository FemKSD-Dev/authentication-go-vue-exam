package model

import (
	"testing"
	"time"
)

func TestUser_Struct(t *testing.T) {
	u := User{
		ID:        "id-1",
		Username:  "john",
		Password:  "hashed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if u.ID != "id-1" {
		t.Errorf("expected ID id-1, got %s", u.ID)
	}
	if u.Username != "john" {
		t.Errorf("expected Username john, got %s", u.Username)
	}
}
