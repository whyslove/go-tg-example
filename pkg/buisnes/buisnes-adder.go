package buisnes

import (
	"context"

	"github.com/rs/zerolog"
)

type Service struct {
	log zerolog.Logger
}

func NewService(log zerolog.Logger) (svc *Service) {
	return &Service{log: log}
}

func (s Service) Add(ctx context.Context, a, b int) (result int, err error) {
	return a + b, nil
}
