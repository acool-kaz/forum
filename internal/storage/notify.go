package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
)

type Notification interface {
	AddNewNotification(notification models.Notification) error
	GetAllNotificationForUser(username string) ([]models.Notification, error)
	DeleteNotification(notificationId int) error
}

type NotificationStorage struct {
	db *sql.DB
}

func newNotificationStorage(db *sql.DB) *NotificationStorage {
	return &NotificationStorage{
		db: db,
	}
}

func (s *NotificationStorage) AddNewNotification(notification models.Notification) error {
	query := `INSERT INTO notification(fromUser, toUser, description, postId) values($1, $2, $3, $4);`
	_, err := s.db.Exec(query, notification.From, notification.To, notification.Description, notification.PostId)
	if err != nil {
		return fmt.Errorf("storage: add new notification: %w", err)
	}
	return nil
}

func (s *NotificationStorage) GetAllNotificationForUser(username string) ([]models.Notification, error) {
	var notifications []models.Notification
	query := `SELECT * FROM notification WHERE toUser = $1;`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("storage: get all notification for user: %w", err)
	}
	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.Id, &notification.From, &notification.To, &notification.Description, &notification.PostId); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("storage: get all notification for user: %w", err)
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (s *NotificationStorage) DeleteNotification(NotificationId int) error {
	query := `DELETE FROM notification WHERE id = $1;`
	_, err := s.db.Exec(query, NotificationId)
	if err != nil {
		return fmt.Errorf("storage: delete notification: %w", err)
	}
	return nil
}
