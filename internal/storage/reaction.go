package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/models"
	"strings"
)

type ReactionStorage struct {
	db *sql.DB
}

func newReactionStorage(db *sql.DB) *ReactionStorage {
	return &ReactionStorage{
		db: db,
	}
}

func (s *ReactionStorage) CreateForPost(ctx context.Context, reaction models.Reaction) error {
	query := fmt.Sprintf("INSERT INTO %s(post_id, user_id, react) VALUES($1, $2, $3);", reactionTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("post reaction storage: create for post: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, reaction.PostId, reaction.UserId, reaction.React); err != nil {
		return fmt.Errorf("post reaction storage: create for post: %w", err)
	}

	return nil
}

func (s *ReactionStorage) CreateForComment(ctx context.Context, reaction models.Reaction) error {
	query := fmt.Sprintf("INSERT INTO %s(post_id, comment_id, user_id, react) VALUES($1, $2, $3, $4);", reactionTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("post reaction storage: create for comment: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, reaction.PostId, reaction.CommentId, reaction.UserId, reaction.React); err != nil {
		return fmt.Errorf("post reaction storage: create for comment: %w", err)
	}

	return nil
}

func (s *ReactionStorage) Get(ctx context.Context, postId, commentId, userId uint) (models.Reaction, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE post_id = $1 AND user_id = $2 AND comment_id IS NULL;", reactionTable)

	args := []interface{}{postId, userId}

	if commentId != 0 {
		query = strings.ReplaceAll(query, "IS NULL", "= $3")
		args = append(args, commentId)
	}

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.Reaction{}, fmt.Errorf("post reaction storage: get: %w", err)
	}
	defer prep.Close()

	var (
		reaction models.Reaction
		cId      sql.NullInt64
	)
	if err = prep.QueryRowContext(ctx, args...).Scan(&reaction.Id, &reaction.PostId, &cId, &reaction.UserId, &reaction.React); err != nil {
		return models.Reaction{}, fmt.Errorf("post reaction storage: get: %w", err)
	}
	reaction.CommentId = uint(cId.Int64)

	return reaction, nil
}

func (s *ReactionStorage) Change(ctx context.Context, id uint, newReact int) error {
	query := fmt.Sprintf("UPDATE %s SET react = $1 WHERE id = $2;", reactionTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("post reaction storage: change: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, newReact, id); err != nil {
		return fmt.Errorf("post reaction storage: change: %w", err)
	}

	return nil
}

func (s *ReactionStorage) Delete(ctx context.Context, id uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", reactionTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("post reaction storage: delete: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("post reaction storage: delete: %w", err)
	}

	return nil
}
