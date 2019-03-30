package constants

const (
	PRIMARY = "primary"
	GENERAL = "general"
	LOCAL   = "local"
)

const (
	SSNKEY        = "ssn~publicKey"
	ELECTION      = "electionType~startDate~endDate~electionID"
	CANDIDATE     = "electionType~ssn"
	VOTING_CHOICE = "electionType~candidate~date~ssn"
)

const (
	IDENTITY = "identity"
	USERKEY  = "userkey"
)

const (
	REGISTERED = "registered"
	VOTED      = "voted"
)

const (
	CANDIDATE_MIN_AGE = 25
	VOTER_MIN_AGE     = 18
)

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
