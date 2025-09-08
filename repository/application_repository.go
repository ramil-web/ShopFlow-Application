package repository

import (
	"database/sql"
	"shopflow/application/models"
)

type ApplicationRepository struct {
	DB *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{DB: db}
}

// Create — создать новую заявку
func (r *ApplicationRepository) Create(app *models.Application) error {
	query := `
		INSERT INTO user_applications (user_id, text, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.DB.QueryRow(
		query,
		app.UserID,
		app.Text,
		app.Status,
	).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
}

// GetAll — получить одну последнюю заявку (если нужен только *Application)
func (r *ApplicationRepository) GetAll(userID uint) ([]models.Application, error) {
	query := `
		SELECT id, user_id, text, status, created_at, updated_at
		FROM user_applications
	`
	var rows *sql.Rows
	var err error

	if userID != 0 {
		query += " WHERE user_id = $1 ORDER BY created_at DESC"
		rows, err = r.DB.Query(query, userID)
	} else {
		query += " ORDER BY created_at DESC"
		rows, err = r.DB.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []models.Application{}
	for rows.Next() {
		var app models.Application
		if err := rows.Scan(&app.ID, &app.UserID, &app.Text, &app.Status, &app.CreatedAt, &app.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}
