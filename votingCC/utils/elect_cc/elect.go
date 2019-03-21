package elect_cc

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	c "../constants"
	u "../keyUtils"
	msg "../msg"
)

var logger = shim.NewLogger("elect_cc")

func (s *ElectChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

func (s *ElectChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "giveVote" {
		return s.giveVote(stub, args)

	} else if function == "getVotingResults" {
		return s.getVotingResults(stub, args)
	}

	return shim.Error(msg.GetErrMsg("COM_ERR_11", []string{function}))
}

// args[0] : ssn
// args[1] : candidatePublic Key
// args[2] : electionType
// args[3] : today Date
func (s *ElectChaincode) giveVote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"giveVote", "4"}))
	}

	err := u.CreateCompKey(stub, c.VOTING_CHOICE, []string{args[2], args[1], args[3], args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}

	result, _ := u.MarshalData(fmt.Sprintf(`{"VoterSSN": "%s", "Candidate":"%s","ElectionType":"%s","ElectionDate":"%s", "TxID": "%s"}`, args[0], args[1], args[2], args[3], stub.GetTxID()), VotingChoice{})

	return shim.Success(result)
}

// @note: demonstrate pagination
// args[0] : electionType
// args[1] : bookmark
// args[2] : page size
// args[3] :
func (s *ElectChaincode) getVotingResults(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"getVotingResults", "4"}))
	}

	electionType := args[0]
	bookmark := args[1]

	pageSize, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		return shim.Error(msg.GetErrMsg("COM_ERR_20", []string{args[2], err.Error()}))
	}

	dataIterator, metadata, err := stub.GetStateByPartialCompositeKeyWithPagination(c.VOTING_CHOICE, []string{electionType}, int32(pageSize), bookmark)
	if err != nil {
		return shim.Error(msg.GetErrMsg("ELECT_ERR_01", []string{err.Error()}))
	}

	defer dataIterator.Close()

	for dataIterator.HasNext() {
		keyIterator, err := dataIterator.Next()
		if err != nil {
			return shim.Error(msg.GetErrMsg("COM_ERR_13", []string{err.Error()}))
		}

		_, keyParts, err := stub.SplitCompositeKey(keyIterator.Key)
		if err != nil {
			return shim.Error(msg.GetErrMsg("COM_ERR_07", []string{err.Error()}))
		}

		fmt.Println(keyParts[1])
	}

	fmt.Println(metadata.Bookmark)

	return shim.Success(nil)
}

/*
func main() {
	err := shim.Start(new(ElectChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
*/
