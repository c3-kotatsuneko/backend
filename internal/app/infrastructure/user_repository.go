package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
	"github.com/c3-kotatsuneko/backend/internal/domain/repository"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepository struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) repository.IUserRepository {
	return &UserRepository{
		client: client,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	// 最小限のバリデーション
	if user.ID == "" {
		return errors.New("user ID is empty")
	}
	if user.Name == "" {
		return errors.New("user name is empty")
	}
	if len(user.Name) > 100 {
		return errors.New("user name is too long")
	}
	existingUser, _ := r.GetUserByName(ctx, user.Name)
	if existingUser != nil {
		return errors.New("user name already exists")
	}

	// メインの処理
	_, err := r.client.Collection("users").Doc(user.ID).Set(ctx,
		map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"password": user.Password,
		})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	doc, err := r.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	var user entity.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	iter := r.client.Collection("users").
		Where("name", "==", name).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by name: %w", err)
	}

	var user entity.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}

	return &user, nil
}
