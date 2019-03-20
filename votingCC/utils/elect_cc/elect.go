package main

import (
	"fmt"

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

func (s *ElectChaincode) getVotingResults(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error(msg.GetErrMsg("COM_ERR_01", []string{"getVotingResults", "4"}))
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(ElectChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
