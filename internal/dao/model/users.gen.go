// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID        string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Email     string         `gorm:"column:email;type:text;not null;uniqueIndex:idx_users_email,priority:1" json:"email"`
	Password  string         `gorm:"column:password;type:text;not null" json:"password"`
	CreatedAt *time.Time     `gorm:"column:created_at;type:timestamp with time zone;index:idx_users_created_at,priority:1;default:now()" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;type:timestamp with time zone;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp with time zone;index:idx_users_deleted_at,priority:1" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
