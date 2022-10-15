package service

import (
	"fmt"
	"forum/internal/storage"
	"forum/models"
)

type Notification interface {
	AddNewNotification(notification models.Notification) error
	GetAllNotificationForUser(user models.User) ([]models.Notification, error)
	DeleteNotification(notificationId int) error
}

type NotificationService struct {
	storage storage.Notification
}

func newNotificationService(storage storage.Notification) *NotificationService {
	return &NotificationService{
		storage: storage,
	}
}

func (s *NotificationService) AddNewNotification(notification models.Notification) error {
	err := s.storage.AddNewNotification(notification)
	if err != nil {
		return fmt.Errorf("service: add new notification: %w", err)
	}
	return nil
}

func (s *NotificationService) GetAllNotificationForUser(user models.User) ([]models.Notification, error) {
	if user == (models.User{}) {
		return nil, nil
	}
	notifies, err := s.storage.GetAllNotificationForUser(user.Username)
	if err != nil {
		return nil, fmt.Errorf("service: get all notification for user: %w", err)
	}
	return notifies, nil
}

func (s *NotificationService) DeleteNotification(notificationId int) error {
	err := s.storage.DeleteNotification(notificationId)
	if err != nil {
		return fmt.Errorf("service: delete notification: %w", err)
	}
	return nil
}
