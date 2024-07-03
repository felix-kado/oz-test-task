package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	model1 "ozon-test/internal/gql/model"
	"ozon-test/internal/models"
	"time"

	"github.com/google/uuid"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, userID string) (*model1.Post, error) {
	post := models.Post{
		ID:            uuid.New(),
		Title:         title,
		Content:       content,
		UserID:        uuid.MustParse(userID),
		AllowComments: true, // Default value, can be changed later
		CreatedAt:     time.Now(),
	}

	err := r.Storage.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return &model1.Post{
		ID:            post.ID.String(),
		Title:         post.Title,
		Content:       post.Content,
		UserID:        post.UserID.String(),
		AllowComments: post.AllowComments,
		CreatedAt:     post.CreatedAt.Format(time.RFC3339),
	}, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string, userID string) (*model1.Comment, error) {
	comment := models.Comment{
		ID:        uuid.New(),
		PostID:    uuid.MustParse(postID),
		Content:   content,
		UserID:    uuid.MustParse(userID),
		CreatedAt: time.Now(),
	}

	if parentID != nil {
		parsedParentID := uuid.MustParse(*parentID)
		comment.ParentID = &parsedParentID
	}

	err := r.Storage.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &model1.Comment{
		ID:        comment.ID.String(),
		PostID:    comment.PostID.String(),
		ParentID:  parentID,
		Content:   comment.Content,
		UserID:    comment.UserID.String(),
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}, nil
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, title *string, content *string, allowComments *bool) (*model1.Post, error) {
	postID := uuid.MustParse(id)
	post, err := r.Storage.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if title != nil {
		post.Title = *title
	}
	if content != nil {
		post.Content = *content
	}
	if allowComments != nil {
		post.AllowComments = *allowComments
	}

	err = r.Storage.UpdatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return &model1.Post{
		ID:            post.ID.String(),
		Title:         post.Title,
		Content:       post.Content,
		UserID:        post.UserID.String(),
		AllowComments: post.AllowComments,
		CreatedAt:     post.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model1.Post, error) {
	postID := uuid.MustParse(id)
	post, err := r.Storage.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return &model1.Post{
		ID:            post.ID.String(),
		Title:         post.Title,
		Content:       post.Content,
		UserID:        post.UserID.String(),
		AllowComments: post.AllowComments,
		CreatedAt:     post.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, page int, pageSize int) ([]*model1.Post, error) {
	posts, err := r.Storage.ListPosts(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var result []*model1.Post
	for _, post := range posts {
		result = append(result, &model1.Post{
			ID:            post.ID.String(),
			Title:         post.Title,
			Content:       post.Content,
			UserID:        post.UserID.String(),
			AllowComments: post.AllowComments,
			CreatedAt:     post.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, page int, pageSize int) ([]*model1.Comment, error) {
	comments, err := r.Storage.GetCommentsByPostID(ctx, uuid.MustParse(postID), page, pageSize)
	if err != nil {
		return nil, err
	}

	var result []*model1.Comment
	for _, comment := range comments {
		parentID := ""
		if comment.ParentID != nil {
			parentID = comment.ParentID.String()
		}

		result = append(result, &model1.Comment{
			ID:        comment.ID.String(),
			PostID:    comment.PostID.String(),
			ParentID:  &parentID,
			Content:   comment.Content,
			UserID:    comment.UserID.String(),
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }