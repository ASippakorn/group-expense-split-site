package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleOwner             = "owner"
	RoleParticipant       = "participant"
	DefaultCurrency       = "THB"
	SplitTypeEqual        = "equal"
	SplitTypeManualAmount = "manual_amount"
	SplitTypePercentage   = "percentage"
	SplitTypeTag          = "tag"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}

func (s *Session) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type Group struct {
	ID              uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name            string        `gorm:"not null"`
	DefaultCurrency string        `gorm:"type:char(3);not null;default:THB"`
	Description     string        `gorm:"not null;default:''"`
	OwnerID         uuid.UUID     `gorm:"type:uuid;not null;index"`
	Owner           User          `gorm:"foreignKey:OwnerID"`
	Participants    []Participant `gorm:"foreignKey:GroupID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (g *Group) BeforeCreate(_ *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	if g.DefaultCurrency == "" {
		g.DefaultCurrency = DefaultCurrency
	}
	return nil
}

type Participant struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_participant_group_user"`
	Group     Group
	UserID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_participant_group_user"`
	User      User
	Role      string `gorm:"not null"`
	Active    bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Participant) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.Role == "" {
		p.Role = RoleParticipant
	}
	return nil
}

type Expense struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	GroupID            uuid.UUID `gorm:"type:uuid;not null;index"`
	Group              Group
	PayerParticipantID uuid.UUID `gorm:"type:uuid;not null;index"`
	PayerParticipant   Participant
	Description        string `gorm:"not null"`
	AmountMinor        int64  `gorm:"not null"`
	Currency           string `gorm:"type:char(3);not null"`
	ExpenseDate        time.Time
	SplitType          string     `gorm:"not null"`
	TagID              *uuid.UUID `gorm:"type:uuid"`
	Tag                *Tag
	Splits             []ExpenseSplit
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Tag struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	GroupID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_tag_group_name"`
	Group        Group
	Name         string        `gorm:"not null;uniqueIndex:idx_tag_group_name"`
	Participants []Participant `gorm:"many2many:tag_participants;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t *Tag) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (e *Expense) BeforeCreate(_ *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.Currency == "" {
		e.Currency = DefaultCurrency
	}
	if e.SplitType == "" {
		e.SplitType = SplitTypeEqual
	}
	return nil
}

type ExpenseSplit struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ExpenseID             uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_expense_split_participant"`
	Expense               Expense
	ParticipantID         uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_expense_split_participant"`
	Participant           Participant
	AmountMinor           int64 `gorm:"not null"`
	PercentageBasisPoints int64 `gorm:"not null;default:0"`
	CreatedAt             time.Time
}

type Settlement struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	GroupID               uuid.UUID `gorm:"type:uuid;not null;index"`
	PayerParticipantID    uuid.UUID `gorm:"type:uuid;not null;index"`
	PayerParticipant      Participant
	ReceiverParticipantID uuid.UUID `gorm:"type:uuid;not null;index"`
	ReceiverParticipant   Participant
	AmountMinor           int64     `gorm:"not null"`
	Currency              string    `gorm:"type:char(3);not null"`
	SettlementDate        time.Time `gorm:"not null"`
	Note                  string    `gorm:"not null;default:''"`
	CreatedAt             time.Time
}

func (s *Settlement) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Currency == "" {
		s.Currency = DefaultCurrency
	}
	return nil
}

func (s *ExpenseSplit) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
