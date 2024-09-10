package database

import (
	"context"
	"io"

	"github.com/h3th-IV/mackerel/internal/models"
)

type Database interface {
	io.Closer

	CaptureData(ctx context.Context, user *models.User) (bool, error)
}
