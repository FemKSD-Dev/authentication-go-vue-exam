package postgres

import (
	"authentication-project-exam/internal/core/model"

	"gorm.io/gorm"
)

func toModel(record *UserRecord) *model.User {
	if record == nil {
		return nil
	}
	return &model.User{
		ID:        record.ID,
		Username:  record.Username,
		Password:  record.Password,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		DeleteAt:  record.DeleteAt.Time,
	}
}

func toRecord(model *model.User) *UserRecord {
	if model == nil {
		return nil
	}
	return &UserRecord{
		ID:        model.ID,
		Username:  model.Username,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeleteAt:  gorm.DeletedAt{Time: model.DeleteAt, Valid: !model.DeleteAt.IsZero()},
	}
}
