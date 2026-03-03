package outbound

import (
	"authentication-project-exam/internal/core/model"
	"context"
)

type UsersRepository interface {
	Save(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}
