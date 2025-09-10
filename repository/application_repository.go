package repository

import (
	"database/sql"
	"errors"
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
		INSERT INTO user_applications (user_id, text, file_url, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.DB.QueryRow(
		query,
		app.UserID,
		app.Text,
		app.FileURL,
		app.Status,
	).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
}

// GetAll — получить одну последнюю заявку (если нужен только *Application)
func (r *ApplicationRepository) GetAll(userID uint) ([]models.Application, error) {
	query := `
		SELECT id, user_id, text, file_url, status, created_at, updated_at
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

	var apps []models.Application
	for rows.Next() {
		var app models.Application
		if err := rows.Scan(&app.ID, &app.UserID, &app.Text, &app.FileURL, &app.Status, &app.CreatedAt, &app.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *ApplicationRepository) GetApplicationById(id uint) (*models.Application, error) {
	var app models.Application
	query := `
    SELECT id, user_id, text, file_url, status, created_at, updated_at
    FROM user_applications
    WHERE id = $1`

	err := r.DB.QueryRow(query, id).Scan(
		&app.ID,
		&app.UserID,
		&app.Text,
		&app.FileURL,
		&app.Status,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) DeleteApplicationById(id uint) error {
	result, err := r.DB.Exec(`DELETE FROM user_applications WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ApplicationRepository) UpdateApplication(app models.Application, id uint) (*models.Application, error) {
	query := `
        UPDATE user_applications
        SET text = $1, status = $2, file_url = $3, updated_at = NOW()
        WHERE id = $4
        RETURNING id, created_at, updated_at`
	err := r.DB.QueryRow(query, app.Text, app.Status, app.FileURL, id).Scan(
		&app.ID,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &app, nil
}
