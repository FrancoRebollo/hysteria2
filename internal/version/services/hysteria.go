package services

import (
	"context"
	"plantilla_api/internal/version/domains"
)

func (s *HysteriaService) AltaBossAPI(ctx context.Context, AltaBoss domains.RequestAltaBoss) error {

	_, err := s.HysteriaRepository.AltaBoss(ctx, AltaBoss)

	if err != nil {
		return err
	}

	return nil
}
