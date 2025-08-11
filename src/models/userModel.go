package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UserModel struct {
	ID                   uint           `json:"id,omitempty" gorm:"primarykey"`
	Uuid                 uuid.UUID      `json:"uuid,omitempty" gorm:"type:uuid; uniqueIndex" filter:"true"`
	IsActive             bool           `json:"is_active,omitempty" gorm:"type:bool; default:true" filter:"true" like:"true"`
	FirstName            string         `json:"first_name,omitempty" gorm:"type:varchar(255); default:null" filter:"true" like:"true" sort:"true"`
	LastName             string         `json:"last_name,omitempty" gorm:"type:varchar(255); default:null" filter:"true" like:"true" sort:"true"`
	FullName             string         `json:"full_name,omitempty" gorm:"-" filter:"true" sort:"true"`
	FatherName           string         `json:"father_name,omitempty" gorm:"type:varchar(255); default:null" filter:"true" like:"true" sort:"true"`
	ProfileImage         string         `json:"profile_image,omitempty" gorm:"type:varchar(255); default:null" filter:"true" sort:"true"`
	Password             []byte         `json:"-" gorm:"type:text; default:null"`
	NationalIdentityCode string         `json:"national_identity_code,omitempty" gorm:"type:varchar(255); uniqueIndex; default:null" filter:"true" like:"true"`
	Mobile               string         `json:"mobile,omitempty" gorm:"type:varchar(100); uniqueIndex; not null" filter:"true" like:"true"`
	Email                string         `json:"email,omitempty" gorm:"type:varchar(100); default:null" filter:"true" like:"true"`
	CreatedAt            time.Time      `json:"created_at,omitempty" sort:"true"`
	UpdatedAt            time.Time      `json:"updated_at,omitempty" sort:"true"`
	DeletedAt            gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" sort:"true"`
}

func (*UserModel) TableName() string {
	return "users"
}
