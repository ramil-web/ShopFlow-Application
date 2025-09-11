package services

import (
	"context"
	"fmt"
	"shopflow/application/models"
	"shopflow/application/publisher"
	"shopflow/application/repository"
)

type ApplicationService struct {
	repo      *repository.ApplicationRepository
	publisher *publisher.ApplicationPublisher
	auth      AuthClient
}

func NewApplicationService(repo *repository.ApplicationRepository, publisher *publisher.ApplicationPublisher, auth AuthClient) *ApplicationService {
	return &ApplicationService{
		repo:      repo,
		publisher: publisher,
		auth:      auth,
	}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, userID uint, token, text, fileURL, status string) (models.Application, error) {
	if s.auth != nil {
		valid, _, err := s.auth.VerifyToken(uint32(userID), token)
		if err != nil || !valid {
			return models.Application{}, fmt.Errorf("invalid token")
		}
	}

	app := models.Application{
		UserID:  userID,
		Text:    text,
		FileURL: fileURL,
		Status:  status,
	}
	if app.Status == "" {
		app.Status = "new"
	}

	if err := s.repo.Create(&app); err != nil {
		return models.Application{}, err
	}

	if s.publisher != nil {
		_ = s.publisher.PublishApplicationCreated(publisher.ApplicationCreatedMessage{
			ID:     app.ID,
			UserID: app.UserID,
			Text:   app.Text,
			File:   app.FileURL,
		})
	}

	return app, nil
}

func (s *ApplicationService) GetAll(userID uint) ([]models.Application, error) {
	return s.repo.GetAll(userID)
}

func (s *ApplicationService) GetApplicationById(id uint) (*models.Application, error) {
	return s.repo.GetApplicationById(id)
}

func (s *ApplicationService) DeleteApplication(id uint) error {
	return s.repo.DeleteApplicationById(id)
}

func (s *ApplicationService) UpdateApplication(req models.UpdateApplicationRequest, id uint) (*models.Application, error) {
	app := models.Application{
		Text:    req.Text,
		FileURL: req.FileURL,
		Status:  req.Status,
	}
	return s.repo.UpdateApplication(app, id)
}
