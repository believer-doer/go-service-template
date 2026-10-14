package task

import "time"

// Task is a deliberately small example domain model. Replace it with the
// resource required by your application while keeping the surrounding layers.
type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null;size:200"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"not null;size:20;default:todo"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateInput = CreateInput
