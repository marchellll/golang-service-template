// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameTodo = "todos"

// Todo mapped from table <todos>
type Todo struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	Text      string         `gorm:"column:text;type:text;not null" json:"text"`
	CreatedAt *time.Time     `gorm:"column:created_at;type:timestamp with time zone;index:idx_created_at,priority:1;default:now()" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;type:timestamp with time zone;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp with time zone;index:idx_deleted_at,priority:1" json:"deleted_at"`
}

// TableName Todo's table name
func (*Todo) TableName() string {
	return TableNameTodo
}