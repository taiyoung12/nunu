package service

import (
	"nunu/pkg/log"
)

type Service struct {
	logger *log.Logger
}

func NewService(
	logger *log.Logger,
) *Service {
	return &Service{
		logger: logger,
	}
}
