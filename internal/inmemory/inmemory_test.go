package inmemory_test

import (
	"context"
	"ozon-test/internal/inmemory"
	"ozon-test/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()
	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 1, "There should be one post")

	createdPost := posts[0]
	assert.Equal(t, post.Title, createdPost.Title, "Titles should match")
	assert.Equal(t, post.Content, createdPost.Content, "Contents should match")
	assert.Equal(t, post.UserID, createdPost.UserID, "User IDs should match")
	assert.True(t, createdPost.AllowComments, "AllowComments should be true")
	assert.WithinDuration(t, time.Now(), createdPost.CreatedAt, time.Second, "CreatedAt should be recent")
}

func TestGetPostByID(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()
	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	fetchedPost, err := storage.GetPostByID(context.Background(), createdPost.ID)
	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, createdPost, fetchedPost, "Posts should match")
}

func TestUpdatePost(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()
	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	createdPost.Title = "Updated Test Post"
	createdPost.AllowComments = false

	err = storage.UpdatePost(context.Background(), createdPost)
	assert.NoError(t, err, "Error should be nil")

	updatedPost, err := storage.GetPostByID(context.Background(), createdPost.ID)
	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, "Updated Test Post", updatedPost.Title, "Titles should match")
	assert.False(t, updatedPost.AllowComments, "AllowComments should be false")
}

func TestCreateComment(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()
	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	comment := models.Comment{
		PostID:  createdPost.ID,
		Content: "This is a test comment.",
		UserID:  uuid.New(),
	}

	err = storage.CreateComment(context.Background(), comment)
	assert.NoError(t, err, "Error should be nil")

	comments, err := storage.GetCommentsByPostID(context.Background(), createdPost.ID, 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, comments, 1, "There should be one comment")

	createdComment := comments[0]
	assert.Equal(t, comment.Content, createdComment.Content, "Contents should match")
	assert.Equal(t, comment.UserID, createdComment.UserID, "User IDs should match")
	assert.WithinDuration(t, time.Now(), createdComment.CreatedAt, time.Second, "CreatedAt should be recent")
}

func TestGetCommentsByPostID(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()
	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	comment1 := models.Comment{
		PostID:  createdPost.ID,
		Content: "This is a test comment 1.",
		UserID:  uuid.New(),
		ID:      uuid.New(),
	}

	err = storage.CreateComment(context.Background(), comment1)
	assert.NoError(t, err, "Error should be nil")

	comment2 := models.Comment{
		PostID:  createdPost.ID,
		Content: "This is a test comment 2.",
		UserID:  uuid.New(),
		ID:      uuid.New(),
	}

	err = storage.CreateComment(context.Background(), comment2)
	assert.NoError(t, err, "Error should be nil")

	comments, err := storage.GetCommentsByPostID(context.Background(), createdPost.ID, 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, comments, 2, "There should be two comments")

	assert.Equal(t, "This is a test comment 1.", comments[0].Content, "First comment's content should match")
	assert.Equal(t, "This is a test comment 2.", comments[1].Content, "Second comment's content should match")
}

func TestPagination(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()

	for i := 0; i < 25; i++ {
		post := models.Post{
			Title:         "Test Post",
			Content:       "This is test post content.",
			UserID:        uuid.New(),
			AllowComments: true,
		}
		err := storage.CreatePost(context.Background(), post)
		assert.NoError(t, err, "Error should be nil")
	}

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 10, "There should be ten posts on the first page")

	posts, err = storage.ListPosts(context.Background(), 2, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 10, "There should be ten posts on the second page")

	posts, err = storage.ListPosts(context.Background(), 3, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 5, "There should be five posts on the third page")
}

func TestCreateAndRetrievePostWithComments(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()

	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}
	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 1, "There should be one post")

	createdPost := posts[0]
	assert.Equal(t, post.Title, createdPost.Title, "Titles should match")
	assert.Equal(t, post.Content, createdPost.Content, "Contents should match")

	comment1 := models.Comment{
		PostID:  createdPost.ID,
		Content: "This is the first test comment.",
		UserID:  uuid.New(),
		ID:      uuid.New(),
	}
	err = storage.CreateComment(context.Background(), comment1)
	assert.NoError(t, err, "Error should be nil")

	comment2 := models.Comment{
		PostID:  createdPost.ID,
		Content: "This is the second test comment.",
		UserID:  uuid.New(),
	}
	err = storage.CreateComment(context.Background(), comment2)
	assert.NoError(t, err, "Error should be nil")

	comments, err := storage.GetCommentsByPostID(context.Background(), createdPost.ID, 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, comments, 2, "There should be two comments")

	assert.Equal(t, "This is the first test comment.", comments[0].Content, "First comment's content should match")
	assert.Equal(t, "This is the second test comment.", comments[1].Content, "Second comment's content should match")
}

func TestCreateUpdateAndRetrievePost(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()

	post := models.Post{
		Title:         "Initial Title",
		Content:       "Initial content.",
		UserID:        uuid.New(),
		AllowComments: true,
	}
	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	createdPost.Title = "Updated Title"
	createdPost.Content = "Updated content."
	err = storage.UpdatePost(context.Background(), createdPost)
	assert.NoError(t, err, "Error should be nil")

	updatedPost, err := storage.GetPostByID(context.Background(), createdPost.ID)
	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, "Updated Title", updatedPost.Title, "Title should be updated")
	assert.Equal(t, "Updated content.", updatedPost.Content, "Content should be updated")
}

func TestCreatePostsAndCommentsWithPagination(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()

	for i := 0; i < 25; i++ {
		post := models.Post{
			Title:         "Test Post",
			Content:       "This is test post content.",
			UserID:        uuid.New(),
			AllowComments: true,
		}
		err := storage.CreatePost(context.Background(), post)
		assert.NoError(t, err, "Error should be nil")
	}

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 10, "There should be ten posts on the first page")

	posts, err = storage.ListPosts(context.Background(), 2, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 10, "There should be ten posts on the second page")

	posts, err = storage.ListPosts(context.Background(), 3, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, posts, 5, "There should be five posts on the third page")

	firstPost := posts[0]
	for i := 0; i < 15; i++ {
		comment := models.Comment{
			PostID:  firstPost.ID,
			Content: "Test comment content.",
			UserID:  uuid.New(),
		}
		err := storage.CreateComment(context.Background(), comment)
		assert.NoError(t, err, "Error should be nil")
	}

	comments, err := storage.GetCommentsByPostID(context.Background(), firstPost.ID, 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, comments, 10, "There should be ten comments on the first page")

	comments, err = storage.GetCommentsByPostID(context.Background(), firstPost.ID, 2, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, comments, 5, "There should be five comments on the second page")
}

func TestNestedComments(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()

	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}
	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	comment1 := models.Comment{
		PostID:  createdPost.ID,
		Content: "This is a top-level comment.",
		UserID:  uuid.New(),
	}
	err = storage.CreateComment(context.Background(), comment1)
	assert.NoError(t, err, "Error should be nil")

	comments, err := storage.GetCommentsByPostID(context.Background(), createdPost.ID, 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdComment1 := comments[0]

	comment2 := models.Comment{
		PostID:   createdPost.ID,
		ParentID: &createdComment1.ID,
		Content:  "This is a nested comment.",
		UserID:   uuid.New(),
	}
	err = storage.CreateComment(context.Background(), comment2)
	assert.NoError(t, err, "Error should be nil")

	comments, err = storage.GetCommentsByPostID(context.Background(), createdPost.ID, 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, comments, 2, "There should be two comments")

	assert.Equal(t, createdComment1.ID, comments[0].ID, "First comment ID should match")
	assert.Equal(t, comment2.Content, comments[1].Content, "Nested comment content should match")
	assert.Equal(t, createdComment1.ID, *comments[1].ParentID, "Nested comment's ParentID should match first comment's ID")
}

func TestInvalidPaginationParameters(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()

	for i := 0; i < 10; i++ {
		post := models.Post{
			Title:         "Test Post",
			Content:       "This is test post content.",
			UserID:        uuid.New(),
			AllowComments: true,
		}
		err := storage.CreatePost(context.Background(), post)
		assert.NoError(t, err, "Error should be nil")
	}

	tests := []struct {
		page     int
		pageSize int
	}{
		{page: 0, pageSize: 10},
		{page: 1, pageSize: 0},
		{page: -1, pageSize: 10},
		{page: 1, pageSize: -10},
	}

	for _, tc := range tests {
		posts, err := storage.ListPosts(context.Background(), tc.page, tc.pageSize)
		assert.Error(t, err, "Error should not be nil for invalid page or pageSize")
		assert.Nil(t, posts, "Posts should be nil for invalid page or pageSize")
	}

	commentTests := []struct {
		page     int
		pageSize int
	}{
		{page: 0, pageSize: 10},
		{page: 1, pageSize: 0},
		{page: -1, pageSize: 10},
		{page: 1, pageSize: -10},
	}

	post := models.Post{
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
	}
	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 1)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	for i := 0; i < 10; i++ {
		comment := models.Comment{
			PostID:  createdPost.ID,
			Content: "This is test comment content.",
			UserID:  uuid.New(),
		}
		err := storage.CreateComment(context.Background(), comment)
		assert.NoError(t, err, "Error should be nil")
	}

	for _, tc := range commentTests {
		comments, err := storage.GetCommentsByPostID(context.Background(), createdPost.ID, tc.page, tc.pageSize)
		assert.Error(t, err, "Error should not be nil for invalid page or pageSize")
		assert.Nil(t, comments, "Comments should be nil for invalid page or pageSize")
	}
}
