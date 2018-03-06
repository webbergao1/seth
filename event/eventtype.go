package event

// TypeEvent event type
type TypeEvent int16

const (
	// EventTxEnterPool event for transaction enter txpool
	EventTxEnterPool TypeEvent = 0
	// EventNewMinedBlock event for new mined block
	EventNewMinedBlock TypeEvent = 1
	// EventNodeDisconnect event for node disconnect
	EventNodeDisconnect TypeEvent = 2
)
