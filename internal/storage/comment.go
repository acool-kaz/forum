package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/models"
)

type CommentStorage struct {
	db *sql.DB
}

func newCommentStorage(db *sql.DB) *CommentStorage {
	return &CommentStorage{
		db: db,
	}
}

func (s *CommentStorage) Create(ctx context.Context, comment models.Comment) error {
	query := fmt.Sprintf("INSERT INTO %s(post_id, user_id, text) VALUES($1, $2, $3);", commentTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("comment storage: create: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, comment.PostId, comment.UserId, comment.Text); err != nil {
		return fmt.Errorf("comment storage: create: %w", err)
	}

	return nil
}

func (s *CommentStorage) GetAll(ctx context.Context, postId uint) ([]models.FullComment, error) {
	query := fmt.Sprintf(`
	SELECT 
		u.username,
		c.text,
		count(SELECT * FROM %s WHERE comment_id = c.id and reaction=1),
		count(SELECT * FROM %s WHERE comment_id = c.id and reaction=-1),
	FROM %s c
	INNER JOIN %s u ON u.id = c.user_id
	WHERE c.id = $1
	GROUP BY c.id;
	`, reactionTable, reactionTable, commentTable, userTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("comment storage: get all: %w", err)
	}
	defer prep.Close()

	var (
		allComments []models.FullComment
		oneComment  models.FullComment
	)

	rows, err := prep.QueryContext(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("comment storage: get all: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&oneComment.Username, &oneComment.Text, &oneComment.Likes, &oneComment.Dislikes); err != nil {
			return nil, fmt.Errorf("comment storage: get all: %w", err)
		}
		allComments = append(allComments, oneComment)
	}

	return allComments, nil
}
