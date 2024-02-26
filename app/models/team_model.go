package models

import (
	"idealist/pkg/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Team model and relations definition
type Team struct {
	ID string `gorm:"primaryKey;column:team_id" json:"id"`

	MemberID string `json:"memberId" gorm:"column:mem_id"`
	Member   Member `json:"member" gorm:"foreignKey:MemberID;references:ID"`

	Name string `gorm:"column:team_name" json:"name" validate:"required;lte=45"`
}

func (t *Team) TableName() string {
	return "team"
}

// Before creating the team, create its UUID and member
func (team *Team) BeforeCreate(tx *gorm.DB) (err error) {
	// Create UUID
	team.ID = uuid.New().String()

	// Create member
	member := Member{
		ID:           uuid.New().String(),
		SubClassID:   team.ID,
		SubClassType: repository.MemberIsUser,
	}

	// Reference member to user
	team.MemberID = member.ID
	team.Member = member

	tx.Create(&member)

	return
}

type TeamRole struct {
	// One-to-one relationship with Member table
	ID   string `gorm:"primaryKey;column:team_role_id" json:"id"`
	Name string `gorm:"column:team_role_name" json:"name" validate:"required;lte=45"`
}

func (t *TeamRole) TableName() string {
	return "teamrole"
}

// Before creating the user, create its UUID and member
func (t *TeamRole) BeforeCreate(tx *gorm.DB) (err error) {
	// Create UUID
	t.ID = uuid.New().String()
	return
}

// Intermediary table for many-to-many relationship between Team and Member
type TeamHasMember struct {
	ID string `gorm:"primaryKey;column:team_mem_id" json:"id"`

	// One-to-many relationship with Team table
	TeamID string `gorm:"column:team_id" json:"teamId"`
	Team   Team   `json:"team" gorm:"foreignKey:TeamID;references:ID"`

	// One-to-many relationship with TeamRole table
	TeamRoleID string   `gorm:"column:team_role_id" json:"teamRoleId"`
	TeamRole   TeamRole `json:"teamRole" gorm:"foreignKey:TeamRoleID;references:ID"`

	// One-to-many relationship with Member table
	MemberID string `gorm:"column:mem_id" json:"memId"`
	Member   Member `json:"member" gorm:"foreignKey:MemberID;references:ID"`

	TeamMemberName string `gorm:"column:team_mem_name" json:"memTeamName" validate:"lte=45"`
}

func (thm *TeamHasMember) TableName() string {
	return "team_has_member"
}

// Before creating the user, create its UUID and member
func (thm *TeamHasMember) BeforeCreate(tx *gorm.DB) (err error) {
	// Create UUID
	thm.ID = uuid.New().String()
	return
}
