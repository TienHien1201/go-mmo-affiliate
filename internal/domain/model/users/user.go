package users

import "time"

type User struct {
	ID          int        `json:"id" gorm:"primary_key" example:"1"`
	TanantID    int        `json:"tanant_id"`
	Email       string     `json:"email" gorm:"unique" example:"abc@gmail.com"`
	FullName    string     `json:"full_name" example:"John Doe"`
	Avatar      string     `json:"avatar,omitempty" example:"https://avatar.com/abc.jpg"`
	Password    string     `json:"password" example:"password"`
	FacebookID  *string    `json:"facebook_id" gorm:"unique"`
	GoogleID    *string    `json:"google_id" gorm:"unique"`
	TotpSecret  *string    `json:"totp_secret" example:"secret"`
	RoleID      int        `json:"role_id" example:"2"`
	LastLoginAt *time.Time `json:"last_login_at" example:"2020-09-06T10:10:10Z"`
	DeletedBy   int        `json:"deleted_by" example:"1"`
	IsDeleted   bool       `json:"is_deleted" example:"0"`
	IsActive    int        `json:"is_active" example:"1"`
	DeletedAt   *time.Time `json:"deleted_at" example:"2020-09-06T10:10:10Z"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-01-01T10:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2025-01-01T10:00:00Z"`
}
