package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	result := entity.User{}
	err := r.db.Table("users").Select("*").Where("id = ?", id).Scan(&result)
	if err.Error != nil {
		return entity.User{}, err.Error
	}
	fmt.Println("ini isi database users : ", result)
	return result, nil // TODO: replace this
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	// fmt.Println("ini get user by email(ctx) : ", ctx)
	// fmt.Println("ini get user by email(emial) : ", email)
	result := entity.User{}
	err := r.db.Table("users").Select("*").Where("email = ?", email).Scan(&result)
	if err.Error != nil {
		return entity.User{}, err.Error
	}
	// fmt.Println("ini isi database users : ", result)
	return result, nil // TODO: replace this
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	// fmt.Println("create user : ", user)
	res := r.db.Create(&user)
	if res.Error != nil {
		return entity.User{}, res.Error
	}
	return user, nil // TODO: replace this
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	return entity.User{}, nil // TODO: replace this
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	return nil // TODO: replace this
}
