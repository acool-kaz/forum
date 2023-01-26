package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/models"
	"strings"
)

type PostStorage struct {
	db *sql.DB
}

func newPostStorage(db *sql.DB) *PostStorage {
	return &PostStorage{
		db: db,
	}
}

func (s *PostStorage) SaveImages(ctx context.Context, postId uint, url string) error {
	query := fmt.Sprintf("INSERT INTO %s(post_id, url) VALUES ($1, $2);", imageTable)

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("post storage: save images: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, postId, url); err != nil {
		return fmt.Errorf("post storage: save images: %w", err)
	}

	return nil
}

func (s *PostStorage) Create(ctx context.Context, post models.Post) (uint, error) {
	query := fmt.Sprintf("INSERT INTO %s(user_id, title, description) VALUES ($1, $2, $3) RETURNING id;", postTable)

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("post storage: create: %w", err)
	}
	defer prep.Close()

	var id uint
	if err = prep.QueryRowContext(ctx, post.UserId, post.Title, post.Description).Scan(&id); err != nil {
		return 0, fmt.Errorf("post storage: crete: %w", err)
	}

	return id, nil
}

func (s *PostStorage) GetAll(ctx context.Context) ([]models.FullPost, error) {
	args := []interface{}{}

	filterCondition := ""
	filter := ctx.Value(models.Filter)
	if filter != nil {
		switch filter.(string) {
		case "likes-most":
			filterCondition = " ORDER BY likes DESC"
		case "likes-least":
			filterCondition = " ORDER BY likes ASC"
		case "time-new":
			filterCondition = " ORDER BY created_at DESC"
		case "time-old":
			filterCondition = " ORDER BY created_at ASC"
		}
	}

	whereCondition := ""

	tag := ctx.Value(models.Tags)
	if tag != nil {
		whereCondition = fmt.Sprintf("WHERE p.id IN (SELECT post_id FROM %s WHERE name = $1)", tagTable)
		args = append(args, tag.(string))
	}

	userId := ctx.Value(models.UserId)

	profilePostFilter := ctx.Value(models.ProfilePostFilter)
	if profilePostFilter != nil {
		switch profilePostFilter.(string) {
		case "created":
			whereCondition = fmt.Sprintf("WHERE u.username = (SELECT username FROM %s WHERE id = $1)", userTable)
		case "liked":
			whereCondition = fmt.Sprintf("WHERE p.id IN (SELECT post_id FROM %s WHERE user_id = $1 AND react=%d)", reactionTable, models.LikeReaction)
		case "disliked":
			whereCondition = fmt.Sprintf("WHERE p.id IN (SELECT post_id FROM %s WHERE user_id = $1 AND react=%d)", reactionTable, models.DislikeReaction)
		case "commented":
			whereCondition = fmt.Sprintf("WHERE p.id IN (SELECT post_id FROM %s WHERE user_id = $1)", commentTable)
		default:
			return nil, fmt.Errorf("post storage: get all: %w", models.ErrInvalidFilter)
		}
		args = append(args, userId.(uint))
	}

	query := fmt.Sprintf(`
	SELECT 
		p.id,
		u.username,
		p.title,
		(SELECT GROUP_CONCAT(t.name, ' ') FROM %s t WHERE t.post_id = p.id) AS tags,
		p.description,
		(SELECT COUNT(*) FROM %s r WHERE r.post_id = p.id AND r.comment_id IS NULL AND react=%d) AS 'likes',
    	(SELECT COUNT(*) FROM %s r WHERE r.post_id = p.id AND r.comment_id IS NULL AND react=%d) AS 'dislikes',
		p.created_at
	FROM %s p 
	INNER JOIN %s u ON u.id = p.user_id
	%s
	GROUP BY p.id%s;
	`, tagTable, reactionTable, models.LikeReaction, reactionTable, models.DislikeReaction, postTable, userTable, whereCondition, filterCondition)

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("post storage: get all: %w", err)
	}
	defer prep.Close()

	var (
		allPosts []models.FullPost
		onePost  models.FullPost
		tags     string
	)

	rows, err := prep.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("post storage: get all: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&onePost.Id, &onePost.Username, &onePost.Title, &tags, &onePost.Description, &onePost.Likes, &onePost.Dislikes, &onePost.CreatedAt); err != nil {
			return nil, fmt.Errorf("post storage: get all: %w", err)
		}
		onePost.Tags = strings.Split(tags, " ")
		allPosts = append(allPosts, onePost)
	}

	return allPosts, nil
}

func (s *PostStorage) GetById(ctx context.Context, id uint) (models.FullPost, error) {
	query := fmt.Sprintf(`
	SELECT 
		p.id,
		u.username,
		p.title,
		(SELECT GROUP_CONCAT(t.name, ' ') FROM %s t WHERE t.post_id = p.id),
		IFNULL((SELECT GROUP_CONCAT(i.url, ' ') FROM %s i WHERE i.post_id = p.id), ''),
		p.description,
		(SELECT COUNT(*) FROM %s r WHERE r.post_id = p.id AND r.comment_id IS NULL AND react=%d) AS 'likes',
    	(SELECT COUNT(*) FROM %s r WHERE r.post_id = p.id AND r.comment_id IS NULL AND react=%d) AS 'dislikes',
		p.created_at
	FROM %s p 
	INNER JOIN %s u ON u.id = p.user_id
	WHERE p.id = $1
	GROUP BY p.id;
	`, tagTable, imageTable, reactionTable, models.LikeReaction, reactionTable, models.DislikeReaction, postTable, userTable)

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.FullPost{}, fmt.Errorf("post storage: get by id: %w", err)
	}
	defer prep.Close()

	var (
		onePost models.FullPost
		tags    string
		images  string
	)

	if err = prep.QueryRowContext(ctx, id).Scan(&onePost.Id, &onePost.Username, &onePost.Title, &tags, &images, &onePost.Description, &onePost.Likes, &onePost.Dislikes, &onePost.CreatedAt); err != nil {
		return models.FullPost{}, fmt.Errorf("post storage: get by id: %w", err)
	}
	onePost.Tags = strings.Split(tags, " ")
	onePost.Images = strings.Split(images, " ")

	return onePost, nil
}

func (s *PostStorage) Delete(ctx context.Context, id uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", postTable)

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("post storage: delete: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("post storage: delete: %w", err)
	}

	return nil
}
