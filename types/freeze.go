package types

// FreezeAction represents the type of freeze operation (freeze or unfreeze)
type FreezeAction uint8

const (
	// FreezeActionUnfreeze indicates an unfreeze operation
	FreezeActionUnfreeze FreezeAction = 0
	// FreezeActionFreeze indicates a freeze operation
	FreezeActionFreeze FreezeAction = 1
)

// String returns the string representation of the FreezeAction
func (f FreezeAction) String() string {
	switch f {
	case FreezeActionUnfreeze:
		return "unfreeze"
	case FreezeActionFreeze:
		return "freeze"
	default:
		return "unknown"
	}
}
