package inmemory

import (
	"context"
	"errors"
	"ozon-test/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type InMemoryStorage struct {
	posts         map[uuid.UUID]models.Post
	comments      map[uuid.UUID]models.Comment
	structure     map[uuid.UUID][]models.StructureTree
	postOrder     []uuid.UUID
	commentOrder  map[uuid.UUID][]uuid.UUID
	postsMutex    sync.RWMutex
	commentsMutex sync.RWMutex
}

// NewInMemoryStorage creates a new instance of InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:        make(map[uuid.UUID]models.Post),
		comments:     make(map[uuid.UUID]models.Comment),
		structure:    make(map[uuid.UUID][]models.StructureTree),
		postOrder:    []uuid.UUID{},
		commentOrder: make(map[uuid.UUID][]uuid.UUID),
	}
}

// CreatePost adds a new post to the in-memory storage.
func (s *InMemoryStorage) CreatePost(ctx context.Context, post models.Post) error {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	s.posts[post.ID] = post
	s.postOrder = append(s.postOrder, post.ID)

	slog.Info("Post created", "postID", post.ID)
	return nil
}

// GetPostByID retrieves a post by its ID from the in-memory storage.
func (s *InMemoryStorage) GetPostByID(ctx context.Context, postID uuid.UUID) (models.Post, error) {
	s.postsMutex.RLock()
	defer s.postsMutex.RUnlock()

	post, exists := s.posts[postID]
	if !exists {
		slog.Warn("Post not found", "postID", postID)
		return models.Post{}, errors.New("post not found")
	}
	return post, nil
}

// ListPosts retrieves a paginated list of posts from the in-memory storage.
func (s *InMemoryStorage) ListPosts(ctx context.Context, page, pageSize int) ([]models.Post, error) {
	if page <= 0 || pageSize <= 0 {
		slog.Warn("Invalid page or pageSize parameter", "page", page, "pageSize", pageSize)
		return nil, errors.New("invalid page or pageSize parameter")
	}

	s.postsMutex.RLock()
	defer s.postsMutex.RUnlock()

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(s.postOrder) {
		return []models.Post{}, nil
	}
	if end > len(s.postOrder) {
		end = len(s.postOrder)
	}

	posts := []models.Post{}
	for _, postID := range s.postOrder[start:end] {
		posts = append(posts, s.posts[postID])
	}

	slog.Info("Listed posts", "page", page, "pageSize", pageSize)
	return posts, nil
}

// CreateComment adds a new comment to the in-memory storage.
func (s *InMemoryStorage) CreateComment(ctx context.Context, comment models.Comment) error {
	s.commentsMutex.Lock()
	defer s.commentsMutex.Unlock()

	comment.CreatedAt = time.Now()
	s.comments[comment.ID] = comment

	var ancestorID uuid.UUID
	var level int

	if comment.ParentID == nil {
		ancestorID = comment.ID
		level = 0
	} else {
		parentComment, exists := s.comments[*comment.ParentID]
		if !exists {
			slog.Warn("Parent comment not found", "parentID", comment.ParentID)
			return errors.New("parent comment not found")
		}
		ancestorID = parentComment.ID
		level = 1
	}

	s.structure[comment.PostID] = append(s.structure[comment.PostID], models.StructureTree{
		AncestorID:        ancestorID,
		DescendantID:      comment.ID,
		NearestAncestorID: ancestorID,
		Level:             level,
		SubjectID:         comment.PostID,
	})

	s.commentOrder[comment.PostID] = append(s.commentOrder[comment.PostID], comment.ID)

	slog.Info("Comment created", "commentID", comment.ID, "postID", comment.PostID)
	return nil
}

// GetCommentByID retrieves a comment by its ID from the in-memory storage.
func (s *InMemoryStorage) GetCommentByID(ctx context.Context, commentID uuid.UUID) (models.Comment, error) {
	s.commentsMutex.RLock()
	defer s.commentsMutex.RUnlock()

	comment, exists := s.comments[commentID]
	if !exists {
		slog.Warn("Comment not found", "commentID", commentID)
		return models.Comment{}, errors.New("comment not found")
	}
	return comment, nil
}

// GetCommentsByPostID retrieves a paginated list of comments for a given postID from the in-memory storage.
func (s *InMemoryStorage) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]models.Comment, error) {
	if page <= 0 || pageSize <= 0 {
		slog.Warn("Invalid page or pageSize parameter", "page", page, "pageSize", pageSize)
		return nil, errors.New("invalid page or pageSize parameter")
	}

	s.commentsMutex.RLock()
	defer s.commentsMutex.RUnlock()

	commentIDs, exists := s.commentOrder[postID]
	if !exists {
		return []models.Comment{}, nil
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(commentIDs) {
		return []models.Comment{}, nil
	}
	if end > len(commentIDs) {
		end = len(commentIDs)
	}

	comments := []models.Comment{}
	for _, commentID := range commentIDs[start:end] {
		comments = append(comments, s.comments[commentID])
	}

	slog.Info("Listed comments for post", "postID", postID, "page", page, "pageSize", pageSize)
	return comments, nil
}

// UpdatePost updates an existing post in the in-memory storage.
func (s *InMemoryStorage) UpdatePost(ctx context.Context, post models.Post) error {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	_, exists := s.posts[post.ID]
	if !exists {
		slog.Warn("Post not found", "postID", post.ID)
		return errors.New("post not found")
	}

	s.posts[post.ID] = post
	slog.Info("Post updated", "postID", post.ID)
	return nil
}
