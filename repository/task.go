package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	// fmt.Println("ini task (ctx) : ", ctx)
	// fmt.Println("ini task (id) : ", id)
	result := []entity.Task{}
	err := r.db.Table("tasks").Select("tasks.*").Where("user_id = ?", id).Scan(&result)
	if err.Error != nil {
		return []entity.Task{}, err.Error
	}
	// fmt.Println("repo getTask tasks setelah filter : ", result)
	return result, nil // TODO: replace this
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	// fmt.Println("ini repo task StoreTask : ", task)
	res := r.db.Create(&task)
	if res.Error != nil {
		return 0, res.Error
	}
	return taskId, nil // TODO: replace this
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	result := entity.Task{}
	err := r.db.Table("tasks").Select("tasks.*").Where("id = ?", id).Scan(&result)
	if err.Error != nil {
		return entity.Task{}, err.Error
	}
	// fmt.Println("ini isi database task : ", result)
	return result, nil // TODO: replace this
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	result := []entity.Task{}
	err := r.db.Table("tasks").Select("tasks.*").Where("id = ?", catId).Scan(&result)
	if err.Error != nil {
		return nil, err.Error
	}
	// fmt.Println("ini isi database task : ", result)
	return result, nil // TODO: replace this
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	// fmt.Println("ini repo task UpdateTask : ", task)
	res := r.db.Table("tasks").Where("id = ?", task.ID).Updates(task)
	if res.Error != nil {
		return res.Error
	}
	return nil // TODO: replace this
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	fmt.Println("ini repo task DeleteTask (id) : ", id)
	taskByID, err := r.GetTaskByID(ctx, id)
	// fmt.Println("ini repo task DeleteTask (id) : ", taskByID.ID)
	if err != nil {
		panic(err)
	}
	// fmt.Println("repo category  DeleteCategory : ", categoriesByID)
	res := r.db.Table("tasks").Where("id = ?", taskByID.ID).Delete(taskByID.ID)
	if res.Error != nil {
		return res.Error
	}
	return nil // TODO: replace this
}
