package domain

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	// Ensure timestamps are always stored in UTC
	now := time.Now().UTC()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	return err
}

func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	// Ensure UpdatedAt is always stored in UTC
	m.UpdatedAt = time.Now().UTC()
	return err
}

type BigInt struct {
	*big.Int
}

func NewBigInt(value *big.Int) BigInt {
	return BigInt{Int: value}
}

func (b BigInt) Unwrap() *big.Int {
	return b.Int
}

// Scan implements the sql.Scanner interface
func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b.Int = new(big.Int)
		return nil
	}
	switch v := value.(type) {
	case []byte:
		b.Int = new(big.Int)
		b.SetString(string(v), 10)
		return nil
	case string:
		b.Int = new(big.Int)
		b.SetString(v, 10)
		return nil
	case int64:
		b.Int = big.NewInt(v)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into BigInt", value)
	}
}

// Value implements the driver.Valuer interface
func (b BigInt) Value() (driver.Value, error) {
	if b.Int == nil {
		return nil, nil //nolint:nilnil // nil value is intentional: unset BigInt has no DB representation
	}
	return b.String(), nil
}

// GormDataType omits explicit `gorm:"type:numeric"` tags on fields using BigInt.
func (b BigInt) GormDataType() string {
	return "numeric"
}

var StringToTokenStatus = map[string]int{
	"new":      0,
	"active":   1,
	"inactive": 2,
}

var TokenStatusToString = map[int]string{
	0: "new",
	1: "active",
	2: "inactive",
}
