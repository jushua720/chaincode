package main

type VotingChaincode struct {
}

type User struct {
	SSN              string `json:"SocalSecurityNumber"`
	PublicKey        string `json:"PublicKey"`
	FirstName        string `json:"FirstName"`
	LastName         string `json:"LastName"`
	DateOfBirth      string `json:"DateOfBirth"`
	Gender           string `json:"Gender"`
	VotingChoice     string `json:"VotingChoice"`
	RegistrationDate string `json:"RegistrationDate"`
}

type Candidate struct {
	SSN            string `json:"SSN"`
	PublicKey      string `json:"PublicKey"`
	FirstName      string `json:"FirstName"`
	LastName       string `json:"LastName"`
	Status         string `json:"Status"`
	ElectionType   string `json:"ElectionType"`
	ElectionPeriod string `json:"ElectionPeriod"`
	ElectionResult string `json:"ElectionResult"`
}

type Election struct {
	ID             string `json:"ID"`
	PublicKey      string `json:"PublicKey"`
	ElectionType   string `json:"ElectionType"`
	ElectionPeriod string `json:"ElectionPeriod"`
	ElectionResult string `json:"ElectionResult"`
}

type NewUser struct {
	SSN              string `json:"SSN"`
	PublicKey        string `json:"PublicKey"`
	PrivateKey       string `json:"PrivateKey"`
	RegistrationDate string `json:"RegistrationDate"`
}

type NewElection struct {
	ElectionType string `json:"ElectionType`
	ElectionID   string `json:"ElectionID"`
	StartDate    string `json:"StartDate"`
	EndDate      string `json:"EndDate"`
	TxID         string `json:"TxID"`
}
type NewCandidate struct {
	SSN            string `json:"SSN"`
	PublicKey      string `json:"PublicKey"`
	FirstName      string `json:"FirstName"`
	LastName       string `json:"LastName"`
	DateOfBirth    string `json:"DateOfBirth"`
	ElectionType   string `json:"ElectionType"`
	ElectionPeriod string `json:"ElectionPeriod"`
	TxID           string `json:"TxID"`
}

type Result struct {
	Res string `json:"Result"`
	Msg string `json:"Message"`
}
