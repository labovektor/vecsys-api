package entity

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	Id            uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventId       string       `json:"event_id"`
	Event         *Event       `json:"event" gorm:"foreignKey:EventId;references:Id"`
	RegionId      string       `json:"region_id"`
	Region        *Region      `json:"region" gorm:"foreignKey:RegionId;references:Id"`
	CategoryId    string       `json:"category_id"`
	Category      *Category    `json:"category" gorm:"foreignKey:CategoryId;references:Id"`
	Name          string       `json:"name"`
	InstitutionId string       `json:"institution_id"`
	Institution   *Institution `json:"institution" gorm:"foreignKey:InstitutionId;references:Id"`
	Email         string       `gorm:"unique" json:"email"`
	Password      string       `json:"-"`
	Biodata       *Biodata     `json:"biodata" gorm:"foreignKey:ParticipantId;references:Id"`
	Payment       *Payment     `json:"payment" gorm:"foreignKey:ParticipantId;references:Id"`
	ProgressStep  string       `json:"progress_step" gorm:"type:participant_progress;default:'registered'"`
	VerifiedAt    *time.Time   `json:"verified_at,omitempty"`
	LockedAt      *time.Time   `json:"locked_at,omitempty"`
	CreatedAt     time.Time    `gorm:"default:now();" json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
}

func (p *Participant) IsVerified() bool {
	return p.VerifiedAt != nil && p.ProgressStep == StepValidatedParticipant
}

func (p *Participant) IsLocked() bool {
	return p.LockedAt != nil && p.ProgressStep == StepLockedParticipant
}

// Enum of participant progress steps 'registered', 'categorized', 'paid', 'validated', 'select_institution', 'fill_biodatas', 'locked'
const (
	StepRegisteredParticipant        = "registered"
	StepCategorizedParticipant       = "categorized"
	StepPaidParticipant              = "paid"
	StepValidatedParticipant         = "validated"
	StepSelectInstitutionParticipant = "select_institution"
	StepFillBiodatasParticipant      = "fill_biodatas"
	StepLockedParticipant            = "locked"
)

func getStepValue(step string) int {
	switch step {
	case StepRegisteredParticipant:
		return 0
	case StepCategorizedParticipant:
		return 1
	case StepPaidParticipant:
		return 2
	case StepValidatedParticipant:
		return 3
	case StepSelectInstitutionParticipant:
		return 4
	case StepFillBiodatasParticipant:
		return 5
	case StepLockedParticipant:
		return 6
	default:
		return 0
	}
}

func (p *Participant) ValidateUserStep(updateStepTarget string) bool {
	updateStepValue := getStepValue(updateStepTarget)

	// Check if user is locked, if locked user data can't be updated
	if p.IsLocked() {
		return false
	}

	// Check if user is verified, if verified user data bellow step 3 can't be updated
	if updateStepValue <= 3 && p.IsVerified() {
		return false
	}

	// Check if user is not verified, if not verified user data above step 3 can't be updated
	if updateStepValue > 3 && !p.IsVerified() {
		return false
	}

	return true
}
