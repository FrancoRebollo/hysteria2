package services

import (
	"context"
	"plantilla_api/internal/version/domains"
)

func (s *HysteriaService) AltaBossAPI(ctx context.Context, AltaBoss domains.RequestAltaBoss) (*domains.AltaBossResponse, error) {

	message, err := s.HysteriaRepository.AltaBoss(ctx, AltaBoss)

	if err != nil {
		return nil, err
	}

	altaBResponse := domains.AltaBossResponse{
		//IdBosses: AltaBoss.IdBosses,
		Nombre:  AltaBoss.Nombre,
		Message: message,
	}

	return &altaBResponse, nil
}

func (s *HysteriaService) AltaAnuncioAPI(ctx context.Context, AltaAnuncio domains.RequestAltaAnuncio) (*domains.AltaAnuncioResponse, error) {

	message, err := s.HysteriaRepository.AltaAnuncio(ctx, AltaAnuncio)

	if err != nil {
		return nil, err
	}

	returnAltaAnuncio := domains.AltaAnuncioResponse{
		//de donde saco el ID?
		Texto: AltaAnuncio.Texto,
		Fecha: AltaAnuncio.Fecha,
		Error: message,
	}

	return &returnAltaAnuncio, nil
}
