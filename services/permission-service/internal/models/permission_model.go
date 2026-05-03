package models

type Permission struct {
	UserID string
	TaskID string
	Role   string // owner
}
