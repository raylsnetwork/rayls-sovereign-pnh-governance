package types

// HeaderFlagReason represents the reason for flagging a participant header
type HeaderFlagReason uint8

const (
	HeaderFlagReasonLiveliness HeaderFlagReason = iota
)

var HeaderFlagReasonToString = map[uint8]string{
	uint8(HeaderFlagReasonLiveliness): "liveliness check failed, header not submitted in time",
}

var StringToHeaderFlagReason = map[string]uint8{
	"liveliness check failed, header not submitted in time": uint8(HeaderFlagReasonLiveliness),
}

// HeaderFlagInitiator represents who initiated the flagging
type HeaderFlagInitiator uint8

const (
	HeaderFlagInitiatorAutomaticSystem HeaderFlagInitiator = iota
)

var HeaderFlagInitiatorToString = map[uint8]string{
	uint8(HeaderFlagInitiatorAutomaticSystem): "AUTOMATIC_SYSTEM_CHECK",
}

var StringToHeaderFlagInitiator = map[string]uint8{
	"AUTOMATIC_SYSTEM_CHECK": uint8(HeaderFlagInitiatorAutomaticSystem),
}
