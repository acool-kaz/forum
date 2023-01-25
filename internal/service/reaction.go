package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/storage"
	"strconv"
)

type ReactionService struct {
	reactionRepo storage.Reaction
}

func newReactionService(reactionRepo storage.Reaction) *ReactionService {
	return &ReactionService{
		reactionRepo: reactionRepo,
	}
}

func (s *ReactionService) Set(ctx context.Context, postId, commentId string, react int, userId uint) error {
	reaction, err := newReaction(postId, commentId, react, userId)
	if err != nil {
		return fmt.Errorf("reaction service: set: %w", err)
	}

	curReact, err := s.reactionRepo.Get(ctx, reaction.PostId, reaction.CommentId, reaction.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if reaction.CommentId == 0 {
				if err = s.reactionRepo.CreateForPost(ctx, reaction); err != nil {
					return fmt.Errorf("reaction service: set: %w", err)
				}
			} else {
				if err = s.reactionRepo.CreateForComment(ctx, reaction); err != nil {
					return fmt.Errorf("reaction service: set: %w", err)
				}
			}
			return nil
		}
		return fmt.Errorf("reaction service: set: %w", err)
	}

	if curReact.React == reaction.React {
		if err = s.reactionRepo.Delete(ctx, curReact.Id); err != nil {
			return fmt.Errorf("reaction service: set: %w", err)
		}
	} else {
		if err = s.reactionRepo.Change(ctx, curReact.Id, reaction.React); err != nil {
			return fmt.Errorf("reaction service: set: %w", err)
		}
	}
	return nil
}

func newReaction(postId, commentId string, react int, userId uint) (models.Reaction, error) {
	postIdInt, err := strconv.ParseUint(postId, 10, 32)
	if err != nil {
		return models.Reaction{}, fmt.Errorf("reaction service: new reaction: %w: invalid post id %v", models.ErrInvalidReaction, postId)
	}

	commentIdInt, err := strconv.ParseUint(commentId, 10, 32)
	if err != nil {
		return models.Reaction{}, fmt.Errorf("reaction service: new reaction: %w: invalid comment id %v", models.ErrInvalidReaction, commentId)
	}

	if react != 1 && react != -1 {
		return models.Reaction{}, fmt.Errorf("reaction service: new reaction: %w: invalid reaction %v", models.ErrInvalidReaction, react)
	}

	return models.Reaction{
		PostId:    uint(postIdInt),
		CommentId: uint(commentIdInt),
		UserId:    userId,
		React:     react,
	}, nil
}
