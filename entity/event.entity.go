package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id                uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	AdminId           string     `json:"admin_id"`
	Admin             *Admin     `json:"admin" gorm:"foreignKey:AdminId;references:Id"`
	Name              string     `json:"name"`
	Desc              string     `json:"desc"`
	GroupMemberNum    int        `gorm:"default:3" json:"group_member_num"`
	Icon              string     `json:"icon"`
	ParticipantTarget int        `json:"participant_target"`
	Period            string     `json:"period"`
	Active            *bool      `gorm:"default:true" json:"active"`
	CreatedAt         time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

type EventReq struct {
	Name              string `json:"name" form:"name"`
	Desc              string `json:"desc" form:"desc"`
	GroupMemberNum    int    `json:"group_member_num" form:"group_member_num"`
	ParticipantTarget int    `json:"participant_target" form:"participant_target"`
	Period            string `json:"period" form:"period"`
}
