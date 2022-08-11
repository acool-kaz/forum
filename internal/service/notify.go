package service

import (
	"fmt"
	"forum/internal/storage"
	"forum/models"
)

type Notify interface {
	AddNewNotify(notify models.Notify) error
	GetAllNotifyForUser(username string) ([]models.Notify, error)
	DeleteNotify(notifyId int) error
}

type NotifyService struct {
	storage storage.Notify
}

func newNotifyService(storage storage.Notify) *NotifyService {
	return &NotifyService{
		storage: storage,
	}
}

func (s *NotifyService) AddNewNotify(notify models.Notify) error {
	err := s.storage.AddNewNotify(notify)
	if err != nil {
		return fmt.Errorf("service: add new notify: %w", err)
	}
	return nil
}

func (s *NotifyService) GetAllNotifyForUser(username string) ([]models.Notify, error) {
	notifies, err := s.storage.GetAllNotifyForUser(username)
	if err != nil {
		return nil, fmt.Errorf("service: get all notify for user: %w", err)
	}
	return notifies, nil
}

func (s *NotifyService) DeleteNotify(notifyId int) error {
	err := s.storage.DeleteNotify(notifyId)
	if err != nil {
		return fmt.Errorf("service: delete notify: %w", err)
	}
	return nil
}
