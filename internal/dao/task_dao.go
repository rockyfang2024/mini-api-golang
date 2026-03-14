package dao

import (
    "gorm.io/gorm"
)

// Task represents the task model
type Task struct {
    ID    uint   `gorm:"primaryKey"`
    Title string `gorm:"column:title;type:varchar(100)"`
    Done  bool   `gorm:"column:done"`
}

// TaskDAO provides methods to access Task data
type TaskDAO struct {
    db *gorm.DB
}

// NewTaskDAO creates a new TaskDAO
func NewTaskDAO(db *gorm.DB) *TaskDAO {
    return &TaskDAO{db}
}

// CreateTask creates a new Task
func (dao *TaskDAO) CreateTask(task *Task) error {
    return dao.db.Create(task).Error
}

// GetTask retrieves a Task by ID
func (dao *TaskDAO) GetTask(id uint) (*Task, error) {
    var task Task
    err := dao.db.First(&task, id).Error
    return &task, err
}

// UpdateTask updates an existing Task
func (dao *TaskDAO) UpdateTask(task *Task) error {
    return dao.db.Save(task).Error
}

// DeleteTask deletes a Task by ID
func (dao *TaskDAO) DeleteTask(id uint) error {
    return dao.db.Delete(&Task{}, id).Error
}