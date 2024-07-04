package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"ozon-test/internal/models"
	"ozon-test/internal/postgres"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		).WithDeadline(60 * time.Second),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	setupSchema(t, db)
	return db
}

func setupSchema(t *testing.T, db *sqlx.DB) {
	schema := `
    CREATE TABLE posts (
        id UUID PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        user_id UUID NOT NULL,
        allow_comments BOOLEAN NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
    CREATE TABLE comments (
        id UUID PRIMARY KEY,
        post_id UUID NOT NULL REFERENCES posts(id),
        parent_id UUID,
        content TEXT NOT NULL,
        user_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
    CREATE TABLE structure_tree (
        ancestor_id UUID NOT NULL,
        descendant_id UUID NOT NULL,
        nearest_ancestor_id UUID NOT NULL,
        level INT NOT NULL,
        subject_id UUID NOT NULL,
        PRIMARY KEY (ancestor_id, descendant_id)
    );`

	_, err := db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
}

func TestCreateAndRetrievePost(t *testing.T) {
	db := setupTestDB(t)
	storage := postgres.NewPostgresStorage(db)

	post := models.Post{
		ID:            uuid.New(),
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
		CreatedAt:     time.Now(),
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	retrievedPost, err := storage.GetPostByID(context.Background(), post.ID)
	assert.NoError(t, err)
	assert.Equal(t, post.ID, retrievedPost.ID)
	assert.Equal(t, post.Title, retrievedPost.Title)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, post.UserID, retrievedPost.UserID)
	assert.Equal(t, post.AllowComments, retrievedPost.AllowComments)
}

func TestCreateAndRetrieveComment(t *testing.T) {
	db := setupTestDB(t)
	storage := postgres.NewPostgresStorage(db)

	post := models.Post{
		ID:            uuid.New(),
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
		CreatedAt:     time.Now(),
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	comment := models.Comment{
		ID:        uuid.New(),
		PostID:    post.ID,
		Content:   "This is a test comment.",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
	}

	err = storage.CreateComment(context.Background(), comment)
	assert.NoError(t, err)

	retrievedComments, err := storage.GetCommentsByPostID(context.Background(), post.ID, 1, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, retrievedComments, "No comments found")

	if len(retrievedComments) == 0 {
		t.Fatal("No comments found")
	}

	retrievedComment := retrievedComments[0]
	assert.Equal(t, comment.ID, retrievedComment.ID)
	assert.Equal(t, comment.PostID, retrievedComment.PostID)
	assert.Equal(t, comment.Content, retrievedComment.Content)
	assert.Equal(t, comment.UserID, retrievedComment.UserID)
}

func TestUpdatePost(t *testing.T) {
	db := setupTestDB(t)
	storage := postgres.NewPostgresStorage(db)

	post := models.Post{
		ID:            uuid.New(),
		Title:         "Initial Title",
		Content:       "Initial content.",
		UserID:        uuid.New(),
		AllowComments: true,
		CreatedAt:     time.Now(),
	}

	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	post.Title = "Updated Title"
	post.Content = "Updated content."
	post.AllowComments = false

	err = storage.UpdatePost(context.Background(), post)
	assert.NoError(t, err)

	retrievedPost, err := storage.GetPostByID(context.Background(), post.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", retrievedPost.Title)
	assert.Equal(t, "Updated content.", retrievedPost.Content)
	assert.False(t, retrievedPost.AllowComments)
}

func TestPagination(t *testing.T) {
	db := setupTestDB(t)
	storage := postgres.NewPostgresStorage(db)

	for i := 0; i < 25; i++ {
		post := models.Post{
			ID:            uuid.New(),
			Title:         "Test Post",
			Content:       "This is test post content.",
			UserID:        uuid.New(),
			AllowComments: true,
			CreatedAt:     time.Now(),
		}
		err := storage.CreatePost(context.Background(), post)
		assert.NoError(t, err)
	}

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, posts, 10)

	posts, err = storage.ListPosts(context.Background(), 2, 10)
	assert.NoError(t, err)
	assert.Len(t, posts, 10)

	posts, err = storage.ListPosts(context.Background(), 3, 10)
	assert.NoError(t, err)
	assert.Len(t, posts, 5)
}

func TestNestedComments(t *testing.T) {
	db := setupTestDB(t)
	storage := postgres.NewPostgresStorage(db)

	post := models.Post{
		ID:            uuid.New(),
		Title:         "Test Post",
		Content:       "This is a test post.",
		UserID:        uuid.New(),
		AllowComments: true,
		CreatedAt:     time.Now(),
	}
	err := storage.CreatePost(context.Background(), post)
	assert.NoError(t, err, "Error should be nil")

	posts, err := storage.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err, "Error should be nil")
	createdPost := posts[0]

	comment1 := models.Comment{
		ID:        uuid.New(),
		PostID:    createdPost.ID,
		Content:   "This is a top-level comment.",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
	}
	err = storage.CreateComment(context.Background(), comment1)
	assert.NoError(t, err, "Error should be nil")

	comments, err := storage.GetCommentsByPostID(context.Background(), createdPost.ID, 1, 10)
	assert.NoError(t, err, "Error should be nil")
	assert.Len(t, comments, 1, "There should be one comment")
	createdComment1 := comments[0]

	comment2 := models.Comment{
		ID:        uuid.New(),
		PostID:    createdPost.ID,
		ParentID:  &createdComment1.ID,
		Content:   "This is a nested comment.",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
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
