package constants

// ElectionTypes
const (
	PRIMARY = "primary"
	GENERAL = "general"
	LOCAL   = "local"
)

// Composite Key Type
const (
	SSN           = "ssn"
	SSNKEY        = "ssn~publicKey"
	ELECTION      = "electionType~startDate~endDate~electionID"
	CANDIDATE     = "electionType~ssn"
	VOTING_CHOICE = "electionType~candidate~date~ssn"
)

// User Query Type
const (
	IDENTITY = "identity"
	USERKEY  = "userkey"
)

const (
	REGISTERED = "registered"
	VOTED      = "voted"
)

// elect_cc
const (
	CCNAME    = "elect_cc"
	CHANNELID = "mychannel"
)

const (
	PLURALITY   = "plurality"
	BORDA       = "borda"
	ELIMINATION = "elimination"
)
const SEPARATOR = "-"

const Base58Table = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
