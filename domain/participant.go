package domain

import (
	"strings"
	"time"
)

// MemberStatus represents the status of a participant
type MemberStatus int

const (
	MemberStatusNew MemberStatus = iota
	MemberStatusActive
	MemberStatusInactive
	MemberStatusFrozen
)

var StringToMemberStatus = map[string]MemberStatus{
	"new":      MemberStatusNew,
	"active":   MemberStatusActive,
	"inactive": MemberStatusInactive,
	"frozen":   MemberStatusFrozen,
}

var MemberStatusToString = map[int]string{
	int(MemberStatusNew):      "new",
	int(MemberStatusActive):   "active",
	int(MemberStatusInactive): "inactive",
	int(MemberStatusFrozen):   "frozen",
}

// MemberRole represents the role of a participant
type MemberRole int

const (
	MemberRoleParticipant MemberRole = iota
	MemberRoleIssuer
	MemberRoleAuditor
)

var StringToMemberRole = map[string]MemberRole{
	"participant": MemberRoleParticipant,
	"issuer":      MemberRoleIssuer,
	"auditor":     MemberRoleAuditor,
}

var MemberRoleToString = map[int]string{
	int(MemberRoleParticipant): "participant",
	int(MemberRoleIssuer):      "issuer",
	int(MemberRoleAuditor):     "auditor",
}

// GetAllowedMemberStatuses returns a comma-separated string of allowed member statuses
func GetAllowedMemberStatuses() string {
	statuses := make([]string, 0, len(StringToMemberStatus))
	for status := range StringToMemberStatus {
		statuses = append(statuses, status)
	}
	return strings.Join(statuses, ", ")
}

// GetAllowedMemberRoles returns a comma-separated string of allowed member roles
func GetAllowedMemberRoles() string {
	roles := make([]string, 0, len(StringToMemberRole))
	for role := range StringToMemberRole {
		roles = append(roles, role)
	}
	return strings.Join(roles, ", ")
}

type Participant struct {
	Model
	ChainId            *uint  `form:"chainId,omitempty" json:"chainId"`
	Name               string `form:"name,omitempty"    json:"name"`
	OwnerId            string `form:"ownerId,omitempty" json:"ownerId"              example:"0x29a97C6EfFB8A411DABc6aDEEfaa84f5067C8bbe"`
	Status             uint8  `                         json:"-"                                                                         gorm:"column:status"`
	StatusStr          string `form:"status,omitempty"  json:"status"                                                                    gorm:"-"`
	Role               uint8  `                         json:"-"                                                                         gorm:"column:role"`
	RoleStr            string `form:"role,omitempty"    json:"role"                                                                      gorm:"-"`
	AllowedToBroadcast bool   `                         json:"allowedToBroadcast"`

	// Flagging fields
	IsFlagged  bool       `json:"isFlagged"`
	FlagReason *string    `json:"flagReason,omitempty"`
	FlaggedAt  *time.Time `json:"flaggedAt,omitempty"`

	// Query filters
	CreatedAfter  *time.Time `gorm:"-" form:"createdAfter,omitempty"  time_format:"2006-01-02" json:"-"`
	CreatedBefore *time.Time `gorm:"-" form:"createdBefore,omitempty" time_format:"2006-01-02" json:"-"`
}

func (p *Participant) GetStatusDescription() string {
	return MemberStatusToString[int(p.Status)]
}

func (p *Participant) GetRoleDescription() string {
	return MemberRoleToString[int(p.Role)]
}
