package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	} else if function == "getUserVotingHistory" {
		return s.getUserVotingHistory(stub, args)
	} else if function == "getAllUsers" {
		return s.getAllUsers(stub, args)

	} else if function == "vote" {
		return s.vote(stub, args)

	} else if function == "countVotes" {
		return s.countVotes(stub, args)
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

	userAsBytes, _ := u.MarshalData(fmt.Sprintf(`{"SSN": "%s", "PublicKey":"%s","FirstName":"%s","LastName":"%s","DateOfBirth":"%s","Gender":"%s","Election":"%s","RegistrationDate":"%s"}`,
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

	registeredElection, err := u.FindCompositeKey(stub, c.ELECTION, []string{electionType})
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

	userAsBytes, err := stub.GetState(userPubKey)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_10", []string{userPubKey, err.Error()}))
	}

	user := User{}
	json.Unmarshal(userAsBytes, &user)

	if userPubKey != pubKey {
		return shim.Error(msg.GetErrMsg("", []string{userPubKey, user.PublicKey}))
	}

	err = u.CreateCompKey(stub, c.CANDIDATE, []string{electionType, ssn})
	if err != nil {
		return shim.Error(err.Error())
	}

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
// args[1] : electionType
// args[2] : candidate
func (s *VotingChaincode) registerVoter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"registerVoter", "2"}))
	}

	var isCandidate bool
	ssn := args[0]
	electionType := args[1]

	found, userPubKey := u.FindUserBySSN(stub, ssn)
	if !found {
		return shim.Error(msg.GetErrMsg("COM_ERR_14", []string{ssn}))
	}

	userAsBytes, err := stub.GetState(userPubKey)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_10", []string{ssn, err.Error()}))
	}

	user := User{}
	json.Unmarshal(userAsBytes, &user)

	election, _ := u.FindCompositeKey(stub, c.ELECTION, []string{electionType})
	if election == "" {
		return shim.Error(msg.GetErrMsg("VOT_ERR_07", []string{electionType}))
	}

	_, keyParts, err := stub.SplitCompositeKey(election)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_07", []string{election}))
	}

	electionStartDate := keyParts[1]
	electionEndDate := keyParts[2]

	isRegistered := strings.Contains(user.Election, fmt.Sprint(c.REGISTERED+
		c.SEPARATOR+electionType+c.SEPARATOR+electionStartDate+
		c.SEPARATOR+electionEndDate))

	if isRegistered == true {
		return shim.Error(msg.GetErrMsg("VOT_ERR_10", []string{ssn}))
	}

	age, isEligibleToVote := u.ValidateAge(user.DateOfBirth, "2006/01/02", electionStartDate, electionEndDate)

	candidateCompKey := fmt.Sprintf("\x00" + c.CANDIDATE + "\x00" + electionType + "\x00" + ssn + "\x00")
	candidateKeyAsBytes, _ := stub.GetState(candidateCompKey)

	if candidateKeyAsBytes != nil {
		isCandidate = true
	}

	electionInfo := fmt.Sprintf(c.REGISTERED +
		c.SEPARATOR + electionType + c.SEPARATOR + electionStartDate +
		c.SEPARATOR + electionEndDate +
		c.SEPARATOR + strconv.FormatBool(isCandidate) +
		c.SEPARATOR + age +
		c.SEPARATOR + strconv.FormatBool(isEligibleToVote))

	user.Election = electionInfo

	userAsBytes, _ = json.Marshal(user)

	err = stub.PutState(userPubKey, userAsBytes)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_09", []string{userPubKey, err.Error()}))
	}

	newVoter := NewVoter{
		ssn,
		user.FirstName, user.LastName,
		user.DateOfBirth, age, isEligibleToVote,
		isCandidate, electionType,
		fmt.Sprint(electionStartDate + "-" + electionEndDate)}

	newVoterJSON, _ := json.Marshal(newVoter)

	return shim.Success(newVoterJSON)
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

// args[0] : ssn
// args[1] : election type
// args[2] : candidate pub key
func (s *VotingChaincode) vote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"vote", "3"}))
	}

	todayDate := string(time.Now().UTC().Format("2006/01/02"))

	voterSSN := args[0]
	electionType := args[1]
	candidatePubKey := args[2]

	election, err := u.FindCompositeKey(stub, c.ELECTION, []string{electionType})
	if err != nil {
		return shim.Error(err.Error())
	}

	if election == "" {
		return shim.Error(msg.GetErrMsg("VOT_ERR_15", []string{electionType}))
	}

	_, keyParts, err := stub.SplitCompositeKey(election)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_07", []string{election}))
	}

	isElectionPeriod := u.IsWithinRange(todayDate, keyParts[1], keyParts[2], "2006/01/02")
	if isElectionPeriod != true {
		return shim.Error(msg.GetErrMsg("VOT_ERR_13", []string{todayDate, electionType, fmt.Sprint(keyParts[1] + "-" + keyParts[2])}))
	}

	found, voterPubKey := u.FindUserBySSN(stub, voterSSN)
	if !found {
		return shim.Error(msg.GetErrMsg("COM_ERR_14", []string{voterSSN}))
	}

	voterAsBytes, err := stub.GetState(voterPubKey)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_10", []string{voterSSN, err.Error()}))
	}

	voter := User{}
	err = json.Unmarshal(voterAsBytes, &voter)

	if err != nil {
		return shim.Error("Failed to unmarshal voter")
	}

	hasVoted := strings.Contains(voter.Election, c.VOTED)
	if hasVoted == true {
		return shim.Error(msg.GetErrMsg("VOT_ERR_14", []string{voterSSN}))
	}

	isRegistered := strings.Contains(voter.Election, c.REGISTERED)
	if isRegistered != true {
		return shim.Error(msg.GetErrMsg("VOT_ERR_11", []string{"Not Registered"}))
	}

	isEligibleToVote := strings.Split(voter.Election, c.SEPARATOR)
	if isEligibleToVote[6] != "true" {
		return shim.Error(msg.GetErrMsg("VOT_ERR_11", []string{"Not Eligible"}))
	}

	voterAge := isEligibleToVote[5]

	candidateAsBytes, err := stub.GetState(candidatePubKey)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_10", []string{voterSSN, err.Error()}))
	}

	candidate := User{}
	json.Unmarshal(candidateAsBytes, &candidate)

	isCandidate := strings.Split(candidate.Election, c.SEPARATOR)
	if isCandidate[4] != "true" {
		return shim.Error(msg.GetErrMsg("VOT_ERR_12", []string{candidatePubKey, "Not Registered"}))
	}

	if candidate.SSN == voter.SSN {
		return shim.Error(msg.GetErrMsg("VOT_ERR_12", []string{candidatePubKey, fmt.Sprint("Same Voter " + voterSSN + " and Candidate " + candidate.SSN)}))
	}

	_, err = s.callOtherCC(stub, c.CCNAME, c.CHANNELID, []string{"giveVote", voter.SSN, candidatePubKey, electionType, todayDate})
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_17", []string{c.CCNAME, err.Error()}))
	}

	voter.Election = strings.Replace(voter.Election, c.REGISTERED, c.VOTED, -1)

	voterAsBytes, _ = json.Marshal(voter)

	err = stub.PutState(voterPubKey, voterAsBytes)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_09", []string{voterPubKey, err.Error()}))
	}

	vote := Vote{
		voterSSN,
		voter.FirstName,
		voter.LastName,
		voterAge,
		candidatePubKey,
		todayDate,
		electionType,
		stub.GetTxID()}

	voteJSON, _ := json.Marshal(vote)

	return shim.Success(voteJSON)
}

func (s *VotingChaincode) callOtherCC(stub shim.ChaincodeStubInterface, ccName string, channelID string, args []string) ([]byte, error) {

	ccInvokeArgs := u.ArrayToChaincodeArgs(args)

	ccInvoke := stub.InvokeChaincode(ccName, ccInvokeArgs, channelID)

	if ccInvoke.Status != shim.OK {
		return []byte{0x00}, errors.New(msg.GetErrMsg("COM_ERR_17", []string{ccName, ccInvoke.Message}))
	}

	return ccInvoke.Payload, nil
}

// args[0] : ssn
func (s *VotingChaincode) getUserVotingHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"getUserVotingHistory", "1"}))
	}

	ssn := args[0]

	found, userPubKey := u.FindUserBySSN(stub, ssn)
	if !found {
		return shim.Error(msg.GetErrMsg("COM_ERR_14", []string{ssn}))
	}

	historyIterator, err := stub.GetHistoryForKey(userPubKey)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_19", []string{ssn, err.Error()}))
	}

	history := make([]User, 0)
	for historyIterator.HasNext() {
		record, err := historyIterator.Next()
		if err != nil {
			return shim.Error(msg.GetErrMsg("COM_ERR_13", []string{err.Error()}))
		}
		var user User

		json.Unmarshal(record.Value, &user)
		history = append(history, user)
	}

	historyAsBytes, _ := json.Marshal(&history)

	return shim.Success(historyAsBytes)
}

// @notice demonstrate pagination
// args[0] : bookmark
// args[1] : page size
func (s *VotingChaincode) getAllUsers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"getAllUsers", "2"}))
	}

	wallets := ""
	separator := " ! "
	bookmark := args[0]

	pageSize, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}

	dataIterator, metadata, err := stub.GetStateByRangeWithPagination("", "", int32(pageSize), bookmark)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer dataIterator.Close()

	logger.Info("MetaData ", metadata)

	for dataIterator.HasNext() {
		key, err := dataIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		logger.Info("Found wallet ", key.Key)

		wallets = fmt.Sprint(wallets + separator + key.Key)

		logger.Info("wallets", wallets)
	}

	return shim.Success([]byte(wallets))

}

// args[0] : voting method
// args[1] : election type
func (s *VotingChaincode) countVotes(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"countVotes", "2"}))
	}

	todayDate := string(time.Now().UTC().Format("2006/01/02"))

	method := args[0]
	electionType := args[1]

	if method != c.PLURALITY && method != c.BORDA && method != c.ELIMINATION {
		return shim.Error(msg.GetErrMsg("COM_ERR_16", []string{method}))
	}

	election, err := u.FindCompositeKey(stub, c.ELECTION, []string{electionType})
	if err != nil {
		return shim.Error(err.Error())
	}

	if election == "" {
		return shim.Error(msg.GetErrMsg("VOT_ERR_15", []string{electionType}))
	}

	_, keyParts, err := stub.SplitCompositeKey(election)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_07", []string{election}))
	}

	electionIsNotOver := u.IsWithinRange(todayDate, keyParts[1], keyParts[2], "2006/01/02")
	if !electionIsNotOver {
		return shim.Error(msg.GetErrMsg("VOT_ERR_17", []string{electionType, fmt.Sprint(keyParts[1] + "-" + keyParts[2]), todayDate}))
	}

	votingRes, err := s.callOtherCC(stub, c.CCNAME, c.CHANNELID, []string{"getVotingResults", electionType})
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_17", []string{c.CCNAME, err.Error()}))
	}

	fmt.Println(votingRes)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(VotingChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
