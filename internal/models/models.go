package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Title         string    `db:"title" json:"title"`
	Content       string    `db:"content" json:"content"`
	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	AllowComments bool      `db:"allow_comments" json:"allow_comments"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type Comment struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	PostID    uuid.UUID  `db:"post_id" json:"post_id"`
	ParentID  *uuid.UUID `db:"parent_id" json:"parent_id,omitempty"` // nil if top-level comment
	Content   string     `db:"content" json:"content"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type StructureTree struct {
	AncestorID        uuid.UUID `db:"ancestor_id" json:"ancestor_id"`
	DescendantID      uuid.UUID `db:"descendant_id" json:"descendant_id"`
	NearestAncestorID uuid.UUID `db:"nearest_ancestor_id" json:"nearest_ancestor_id"`
	Level             int       `db:"level" json:"level"`
	SubjectID         uuid.UUID `db:"subject_id" json:"subject_id"`
}

type Storage interface {
	CreatePost(ctx context.Context, post Post) error
	GetPostByID(ctx context.Context, postID uuid.UUID) (Post, error)
	ListPosts(ctx context.Context, page, pageSize int) ([]Post, error)
	CreateComment(ctx context.Context, comment Comment) error
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]Comment, error)
	UpdatePost(ctx context.Context, post Post) error
	GetCommentByID(ctx context.Context, commentID uuid.UUID) (Comment, error)
}

var ErrPostNotFound = errors.New("post not found")
var ErrCommentNotFound = errors.New("comment not found")
