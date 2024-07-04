package postgres

import (
	"context"
	"database/sql"
	"ozon-test/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

type PostgresStorage struct {
	db *sqlx.DB
}

// NewPostgresStorage creates a new instance of PostgresStorage.
func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// CreatePost inserts a new post into the database.
func (s *PostgresStorage) CreatePost(ctx context.Context, post models.Post) error {
	query := `INSERT INTO posts (id, title, content, user_id, allow_comments, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, query, post.ID, post.Title, post.Content, post.UserID, post.AllowComments, post.CreatedAt)
	if err != nil {
		slog.Error("Failed to create post", "error", err, "postID", post.ID)
	}
	return err
}

// GetPostByID retrieves a post by its ID from the database.
func (s *PostgresStorage) GetPostByID(ctx context.Context, postID uuid.UUID) (models.Post, error) {
	var post models.Post
	query := `SELECT id, title, content, user_id, allow_comments, created_at FROM posts WHERE id = $1`
	err := s.db.GetContext(ctx, &post, query, postID)
	if err == sql.ErrNoRows {
		slog.Warn("Post not found", "postID", postID)
		return post, models.ErrPostNotFound
	}
	if err != nil {
		slog.Error("Failed to get post", "error", err, "postID", postID)
	}
	return post, err
}

// ListPosts retrieves a paginated list of posts from the database.
func (s *PostgresStorage) ListPosts(ctx context.Context, page, pageSize int) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT id, title, content, user_id, allow_comments, created_at 
              FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := s.db.SelectContext(ctx, &posts, query, pageSize, (page-1)*pageSize)
	if err != nil {
		slog.Error("Failed to list posts", "error", err)
	}
	return posts, err
}

// CreateComment inserts a new comment into the database and updates the structure_tree table.
func (s *PostgresStorage) CreateComment(ctx context.Context, comment models.Comment) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return err
	}

	query := `INSERT INTO comments (id, post_id, parent_id, content, user_id, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, query, comment.ID, comment.PostID, comment.ParentID, comment.Content, comment.UserID, comment.CreatedAt)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			slog.Error("Failed to rollback transaction", "error", rbErr)
			return rbErr
		}
		slog.Error("Failed to create comment", "error", err)
		return err
	}

	if comment.ParentID == nil {
		query = `INSERT INTO structure_tree (ancestor_id, descendant_id, nearest_ancestor_id, level, subject_id) 
			 VALUES ($1, $1, $1, 0, $2)`
		_, err = tx.ExecContext(ctx, query, comment.ID, comment.PostID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				slog.Error("Failed to rollback transaction", "error", rbErr)
				return rbErr
			}
			slog.Error("Failed to update structure tree", "error", err)
			return err
		}
	} else {
		query = `INSERT INTO structure_tree (ancestor_id, descendant_id, nearest_ancestor_id, level, subject_id) 
			 SELECT ancestor_id, $1, $2, level + 1, subject_id 
			 FROM structure_tree 
			 WHERE descendant_id = $2`
		_, err = tx.ExecContext(ctx, query, comment.ID, comment.ParentID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				slog.Error("Failed to rollback transaction", "error", rbErr)
				return rbErr
			}
			slog.Error("Failed to update structure tree", "error", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to commit transaction", "error", err)
	}
	return err
}

// GetCommentsByPostID retrieves a paginated list of comments for a given postID from the database.
func (s *PostgresStorage) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]models.Comment, error) {
	var comments []models.Comment
	query := `SELECT comments.id, comments.post_id, comments.parent_id, comments.content, comments.user_id, comments.created_at 
              FROM comments 
              JOIN structure_tree ON comments.id = structure_tree.descendant_id 
              WHERE structure_tree.subject_id = $1 
              AND structure_tree.nearest_ancestor_id = structure_tree.ancestor_id
              ORDER BY comments.created_at ASC 
              LIMIT $2 OFFSET $3`
	err := s.db.SelectContext(ctx, &comments, query, postID, pageSize, (page-1)*pageSize)
	if err != nil {
		slog.Error("Failed to get comments by post ID", "error", err, "postID", postID)
	}
	return comments, err
}

// UpdatePost updates the details of an existing post in the database.
func (s *PostgresStorage) UpdatePost(ctx context.Context, post models.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, allow_comments = $3 WHERE id = $4`
	_, err := s.db.ExecContext(ctx, query, post.Title, post.Content, post.AllowComments, post.ID)
	if err != nil {
		slog.Error("Failed to update post", "error", err, "postID", post.ID)
	}
	return err
}

// GetCommentByID retrieves a comment by its ID from the database.
func (s *PostgresStorage) GetCommentByID(ctx context.Context, commentID uuid.UUID) (models.Comment, error) {
	var comment models.Comment
	query := `SELECT id, post_id, parent_id, content, user_id, created_at FROM comments WHERE id = $1`
	err := s.db.GetContext(ctx, &comment, query, commentID)
	if err == sql.ErrNoRows {
		slog.Warn("Comment not found", "commentID", commentID)
		return comment, models.ErrCommentNotFound
	}
	if err != nil {
		slog.Error("Failed to get comment", "error", err, "commentID", commentID)
	}
	return comment, err
}
