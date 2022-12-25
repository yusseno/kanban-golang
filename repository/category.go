package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	// fmt.Println("repo GetCategoriesByUserId (id) : ", id)
	result := []entity.Category{}
	err := r.db.Table("categories").Select("*").Where("user_id = ?", id).Scan(&result)
	if err.Error != nil {
		return nil, err.Error
	}
	// fmt.Println("repo GetCategoriesByUserId : ", result)
	return result, nil // TODO: replace this
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	// fmt.Println("create category : ", category)
	res := r.db.Create(category)
	if res.Error != nil {
		return 0, res.Error
	}
	return category.ID, nil // TODO: replace this
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	res := r.db.Create(&categories)
	if res.Error != nil {
		return res.Error
	}
	return nil // TODO: replace this
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	result := entity.Category{}
	err := r.db.Table("categories").Select("*").Where("id = ?", id).Scan(&result)
	if err.Error != nil {
		return entity.Category{}, err.Error
	}
	// fmt.Println("ini isi database users : ", result)
	return result, nil // TODO: replace this
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	res := r.db.Table("categories").Where("id = ?", category.ID).Updates(category)
	if res.Error != nil {
		return res.Error
	}
	return nil // TODO: replace this
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	categoriesByID, err := r.GetCategoryByID(ctx, id)
	if err != nil {
		panic(err)
	}
	// fmt.Println("repo category  DeleteCategory : ", categoriesByID)
	res := r.db.Table("categories").Where("id = ?", categoriesByID.ID).Delete(categoriesByID.ID)
	if res.Error != nil {
		return res.Error
	}
	return nil // TODO: replace this
}
