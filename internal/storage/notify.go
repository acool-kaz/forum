package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
)

type Notify interface {
	AddNewNotify(notify models.Notify) error
	GetAllNotifyForUser(username string) ([]models.Notify, error)
	DeleteNotify(notifyId int) error
}

type NotifyStorage struct {
	db *sql.DB
}

func newNotifyStorage(db *sql.DB) *NotifyStorage {
	return &NotifyStorage{
		db: db,
	}
}

func (s *NotifyStorage) AddNewNotify(notify models.Notify) error {
	query := `INSERT INTO notify(fromUser, toUser, description, postId) values($1, $2, $3, $4);`
	_, err := s.db.Exec(query, notify.From, notify.To, notify.Description, notify.PostId)
	if err != nil {
		return fmt.Errorf("storage: add new notify: %w", err)
	}
	return nil
}

func (s *NotifyStorage) GetAllNotifyForUser(username string) ([]models.Notify, error) {
	var notifies []models.Notify
	query := `SELECT id, fromUser, toUser, description, postId FROM notify WHERE toUser = $1;`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("storage: get all notify for user: %w", err)
	}
	for rows.Next() {
		var notify models.Notify
		if err := rows.Scan(&notify.Id, &notify.From, &notify.To, &notify.Description, &notify.PostId); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("storage: get all notify for user: %w", err)
		}
		notifies = append(notifies, notify)
	}
	return notifies, nil
}

func (s *NotifyStorage) DeleteNotify(notifyId int) error {
	query := `DELETE FROM notify WHERE id = $1;`
	_, err := s.db.Exec(query, notifyId)
	if err != nil {
		return fmt.Errorf("storage: delete notify: %w", err)
	}
	return nil
}
