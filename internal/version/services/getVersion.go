package services

import (
	"context"
	"plantilla_api/cmd/config"
	"plantilla_api/internal/version/domains"
	"plantilla_api/internal/version/ports"
)

type SecurityService struct {
	ports.SecurityRepository
	config.App
}

//Constructor de la estructura VersionService

func NewSecurityService(repo ports.SecurityRepository, app config.App) *SecurityService {
	return &SecurityService{
		SecurityRepository: repo,
		App:                app,
	}
}

type HysteriaService struct {
	ports.HysteriaRepository
	config.App
}

//Constructor de la estructura VersionService

func NewHysteriaService(repo ports.HysteriaRepository, app config.App) *HysteriaService {
	return &HysteriaService{
		HysteriaRepository: repo,
		App:                app,
	}
}

func (s *SecurityService) GetVersionAPI(ctx context.Context) (*domains.Version, error) {

	version_api, err := s.SecurityRepository.GetVersion(ctx)
	if err != nil {
		return nil, err
	}

	newVersion := domains.Version{
		NombreApi:     s.App.Name,
		Cliente:       s.App.Client,
		Version:       s.App.Version,
		FechaStartUp:  s.App.FechaStartUp,
		VersionModelo: version_api,
	}
	return &newVersion, nil

}
