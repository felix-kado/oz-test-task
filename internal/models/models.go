package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	UserID        uuid.UUID `json:"user_id"`
	AllowComments bool      `json:"allow_comments"`
	CreatedAt     time.Time `json:"created_at"`
}

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	PostID    uuid.UUID  `json:"post_id"`
	ParentID  *uuid.UUID `json:"parent_id"` // nil if top-level comment
	Content   string     `json:"content"`
	UserID    uuid.UUID  `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
}

type StructureTree struct {
	AncestorID        uuid.UUID `json:"ancestor_id"`
	DescendantID      uuid.UUID `json:"descendant_id"`
	NearestAncestorID uuid.UUID `json:"nearest_ancestor_id"`
	Level             int       `json:"level"`
	SubjectID         uuid.UUID `json:"subject_id"`
}

type Storage interface {
	CreatePost(ctx context.Context, post Post) error
	GetPostByID(ctx context.Context, postID uuid.UUID) (Post, error)
	ListPosts(ctx context.Context, page, pageSize int) ([]Post, error)
	CreateComment(ctx context.Context, comment Comment) error
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]Comment, error)
	UpdatePost(ctx context.Context, post Post) error
}
