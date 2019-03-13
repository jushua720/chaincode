package constants

// ElectionTypes
const (
	PRIMARY = "primary"
	GENERAL = "general"
	LOCAL   = "local"
)

// Composite Key Type
const (
	SSN       = "ssn"
	SSNKEY    = "ssn~publicKey"
	ELECTION  = "electionType~startDate~endDate~electionID"
	CANDIDATE = "electionType~ssn"
)

// User Query Type
const (
	IDENTITY = "identity"
	USERKEY  = "userkey"
)

const Base58Table = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
