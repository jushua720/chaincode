package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	c "./utils/constants"
	u "./utils/keyUtils"
	msg "./utils/msg"
)

var logger = shim.NewLogger("voting_cc")

func (s *VotingChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

// @dev: init - access control

func (s *VotingChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "registerUser" {
		return s.registerUser(stub, args)

	} else if function == "registerElection" {
		return s.registerElection(stub, args)
	} else if function == "registerCandidate" {
		return s.registerCandidate(stub, args)
	} else if function == "registerVoter" {
		return s.registerVoter(stub, args)

	} else if function == "getUser" {
		return s.getUser(stub, args)
	}

	return shim.Error(msg.GetErrMsg("COM_ERR_11", []string{function}))
}

// args[0] : SSN
// args[1] : FirstName
// args[2] : LastName
// args[3] : DateOfBirth
// args[4] : Gender
func (s *VotingChaincode) registerUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"registerUser", "5"}))
	}

	ssn := args[0]
	gender := args[4]
	registrationDate := string(time.Now().UTC().Format("2006/01/02 15:04:05"))

	found, _ := u.FindUserBySSN(stub, ssn)

	if found == true {
		return shim.Error(msg.GetErrMsg("VOT_ERR_01", []string{ssn}))
	}

	isValid := u.ValidateArgument(gender)
	if isValid != true {
		return shim.Error(msg.GetErrMsg("COM_ERR_18", []string{"gender", gender}))
	}

	pubKey, privKey, err := u.GenerateKeys()
	if err != nil {
		return shim.Error(msg.GetErrMsg("VOT_ERR_03", []string{err.Error()}))
	}

	userAsBytes, _ := u.MarshalData(fmt.Sprintf(`{"SSN": "%s", "PublicKey":"%s","FirstName":"%s","LastName":"%s","DateOfBirth":"%s","Gender":"%s","VotingChoice":"%s","RegistrationDate":"%s"}`,
		ssn, pubKey, args[1], args[2], args[3], gender, "", registrationDate), User{})

	err = stub.PutState(pubKey, userAsBytes)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_09", []string{pubKey, err.Error()}))
	}

	err = u.CreateCompKey(stub, c.SSN, []string{c.SSNKEY, ssn, pubKey})
	if err != nil {
		return shim.Error(err.Error())
	}

	result, _ := u.MarshalData(fmt.Sprintf(`{"ssn": "%s", "PublicKey":"%s","PrivateKey":"%s","RegistrationDAte":"%s"}`, ssn, pubKey, privKey, registrationDate), NewUser{})

	return shim.Success(result)
}

// args[0] : election Type
// args[1] : electionID
// args[2] : start date
// args[3] : end date
func (s *VotingChaincode) registerElection(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"registerElection", "4"}))
	}

	electionType := args[0]

	electionID := args[1]
	startDate := args[2]
	endDate := args[3]

	if electionType != c.PRIMARY && electionType != c.GENERAL && electionType != c.LOCAL {
		return shim.Error(msg.GetErrMsg("VOT_ERR_04", []string{electionType}))
	}

	isValid := u.ValidateElectionPeriod(startDate, endDate)
	if isValid != true {
		return shim.Error(msg.GetErrMsg("VOT_ERR_05", []string{startDate, endDate}))
	}

	registeredElection, err := u.FindCompositeKey(stub, c.ELECTION, []string{})
	if registeredElection != "" {
		return shim.Error(msg.GetErrMsg("VOT_ERR_06", []string{registeredElection}))
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	err = u.CreateCompKey(stub, c.ELECTION, []string{electionType, startDate, endDate, electionID})
	if err != nil {
		return shim.Error(err.Error())
	}

	newElection := NewElection{electionType, electionID, startDate, endDate, stub.GetTxID()}
	newElectionJSON, _ := json.Marshal(newElection)

	return shim.Success(newElectionJSON)
}

// args[0] : election type
// args[1] : ssn
// args[2] : pubKey
func (s *VotingChaincode) registerCandidate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"registerCandidate", "3"}))
	}

	electionType := args[0]
	ssn := args[1]
	pubKey := args[2]

	if electionType != c.PRIMARY && electionType != c.GENERAL && electionType != c.LOCAL {
		return shim.Error(msg.GetErrMsg("VOT_ERR_04", []string{electionType}))
	}

	election, _ := u.FindCompositeKey(stub, c.ELECTION, []string{electionType})
	if election == "" {
		return shim.Error(msg.GetErrMsg("VOT_ERR_07", []string{electionType}))
	}

	candidateCompKey := fmt.Sprintf("\x00" + c.CANDIDATE + "\x00" + electionType + "\x00" + ssn + "\x00")
	candidateKeyAsBytes, _ := stub.GetState(candidateCompKey)
	if candidateKeyAsBytes != nil {
		return shim.Error(msg.GetErrMsg("VOT_ERR_09", []string{candidateCompKey}))
	}

	found, userPubKey := u.FindUserBySSN(stub, ssn)
	if !found {
		return shim.Error(msg.GetErrMsg("COM_ERR_14", []string{ssn}))
	}

	// @dev - check if registrationPeriod is open

	userAsBytes, err := stub.GetState(userPubKey)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_10", []string{userPubKey, err.Error()}))
	}

	user := User{}
	json.Unmarshal(userAsBytes, &user)

	if userPubKey != pubKey {
		return shim.Error(msg.GetErrMsg("", []string{userPubKey, user.PublicKey}))
	}

	// @dev check candidate age

	err = u.CreateCompKey(stub, c.CANDIDATE, []string{electionType, ssn})
	if err != nil {
		return shim.Error(err.Error())
	}

	// @optimize
	_, keyParts, err := stub.SplitCompositeKey(election)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_07", []string{election}))
	}

	electionPeriod := fmt.Sprint(keyParts[1] + " - " + keyParts[2])

	newCandidate := NewCandidate{ssn, pubKey, user.FirstName, user.LastName, user.DateOfBirth, electionType, electionPeriod, stub.GetTxID()}
	newCandidateJSON, _ := json.Marshal(newCandidate)

	return shim.Success(newCandidateJSON)

}

// args[0] : SSN
// args[] : electionType
func (s *VotingChaincode) registerVoter(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

// args[0]: search criteria [identity (ssn) / publickey]
// args[1]: ssn (national id)/ pub Key
func (s *VotingChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"getUser", "2"}))
	}

	var found bool
	queryType := args[0]
	user := args[1]

	if queryType != c.IDENTITY && queryType != c.USERKEY {
		shim.Error(msg.GetErrMsg("COM_ERR_12", []string{queryType, fmt.Sprintf(c.IDENTITY + " or " + c.USERKEY)}))
	}

	if queryType == c.IDENTITY {
		found, user = u.FindUserBySSN(stub, user)
		if !found {
			return shim.Error(msg.GetErrMsg("COM_ERR_14", []string{args[1]}))
		}
	}

	userAsBytes, err := stub.GetState(user)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_10", []string{user, err.Error()}))
	}

	return shim.Success(userAsBytes)

}

func main() {
	err := shim.Start(new(VotingChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
