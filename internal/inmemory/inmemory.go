package inmemory

import (
	"context"
	"errors"
	"ozon-test/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
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

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:        make(map[uuid.UUID]models.Post),
		comments:     make(map[uuid.UUID]models.Comment),
		structure:    make(map[uuid.UUID][]models.StructureTree),
		postOrder:    []uuid.UUID{},
		commentOrder: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (s *InMemoryStorage) CreatePost(ctx context.Context, post models.Post) error {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	s.posts[post.ID] = post
	s.postOrder = append(s.postOrder, post.ID)
	return nil
}

func (s *InMemoryStorage) GetPostByID(ctx context.Context, postID uuid.UUID) (models.Post, error) {
	s.postsMutex.RLock()
	defer s.postsMutex.RUnlock()

	post, exists := s.posts[postID]
	if !exists {
		return models.Post{}, errors.New("post not found")
	}
	return post, nil
}

func (s *InMemoryStorage) ListPosts(ctx context.Context, page, pageSize int) ([]models.Post, error) {
	s.postsMutex.RLock()
	defer s.postsMutex.RUnlock()

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(s.postOrder) {
		return []models.Post{}, nil
	}
	if end > len(s.postOrder) {
		end = len(s.postOrder)
	}

	posts := []models.Post{}
	for _, postID := range s.postOrder[start:end] {
		posts = append(posts, s.posts[postID])
	}

	return posts, nil
}

func (s *InMemoryStorage) CreateComment(ctx context.Context, comment models.Comment) error {
	s.commentsMutex.Lock()
	defer s.commentsMutex.Unlock()

	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	s.comments[comment.ID] = comment

	// Update structure tree
	var ancestorID uuid.UUID
	var level int

	if comment.ParentID == nil {
		ancestorID = comment.ID
		level = 0
	} else {
		parentComment, exists := s.comments[*comment.ParentID]
		if !exists {
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

	return nil
}

// GetCommentByID retrieves a comment by its ID
func (s *InMemoryStorage) GetCommentByID(ctx context.Context, commentID uuid.UUID) (models.Comment, error) {
	s.commentsMutex.RLock()
	defer s.commentsMutex.RUnlock()

	comment, exists := s.comments[commentID]
	if !exists {
		return models.Comment{}, errors.New("comment not found")
	}
	return comment, nil
}
func (s *InMemoryStorage) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]models.Comment, error) {
	s.commentsMutex.RLock()
	defer s.commentsMutex.RUnlock()

	commentIDs, exists := s.commentOrder[postID]
	if !exists {
		return []models.Comment{}, nil
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(commentIDs) {
		return []models.Comment{}, nil
	}
	if end > len(commentIDs) {
		end = len(commentIDs)
	}

	comments := []models.Comment{}
	for _, commentID := range commentIDs[start:end] {
		comments = append(comments, s.comments[commentID])
	}

	return comments, nil
}

func (s *InMemoryStorage) UpdatePost(ctx context.Context, post models.Post) error {
	s.postsMutex.Lock()
	defer s.postsMutex.Unlock()

	_, exists := s.posts[post.ID]
	if !exists {
		return errors.New("post not found")
	}

	s.posts[post.ID] = post
	return nil
}
