package main

type ElectChaincode struct {
}

type VotingChoice struct {
	VoterSSN     string `json:"VoterSSN"`
	Candidate    string `json:"Candidate"`
	ElectionType string `json:"ElectionType"`
	ElectionDate string `json:"ElectionDate"`
	TxID         string `json:"TxID"`
}
