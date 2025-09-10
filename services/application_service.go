package services

import (
	"shopflow/application/models"
	"shopflow/application/publisher"
	"shopflow/application/repository"
)

type ApplicationService struct {
	repo      *repository.ApplicationRepository
	publisher *publisher.ApplicationPublisher
}

func NewApplicationService(repo *repository.ApplicationRepository, publisher *publisher.ApplicationPublisher) *ApplicationService {
	return &ApplicationService{repo: repo, publisher: publisher}
}

func (s *ApplicationService) CreateApplication(userID uint, text, fileURL string, Status string) (models.Application, error) {
	app := models.Application{
		UserID:  userID,
		Text:    text,
		FileURL: fileURL,
		Status: func() string {
			if Status == "" {
				return "new"
			}
			return Status
		}(),
	}

	if err := s.repo.Create(&app); err != nil {
		return models.Application{}, err
	}

	// событие в шину
	if s.publisher != nil {
		msg := publisher.ApplicationCreatedMessage{
			ID:     app.ID,
			UserID: app.UserID,
			Text:   app.Text,
			File:   app.FileURL,
		}
		_ = s.publisher.PublishApplicationCreated(msg)
	}

	return app, nil
}

func (s *ApplicationService) GetAll(UserID uint) ([]models.Application, error) {
	return s.repo.GetAll(UserID)
}

func (s *ApplicationService) GetApplicationById(id uint) (*models.Application, error) {
	return s.repo.GetApplicationById(id)
}

func (s *ApplicationService) DeleteApplication(id uint) error {
	return s.repo.DeleteApplicationById(id)
}

func (s *ApplicationService) UpdateApplication(req models.UpdateApplicationRequest, Id uint) (*models.Application, error) {
	app := models.Application{
		Text:    req.Text,
		FileURL: req.FileURL,
		Status:  req.Status,
	}
	return s.repo.UpdateApplication(app, Id)
}
