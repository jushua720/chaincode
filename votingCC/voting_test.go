package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"./utils/elect_cc"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Init(test *testing.T) *shim.MockStub {
	stub := shim.NewMockStub("VotingCCTestStub", new(VotingChaincode))
	result := stub.MockInit("000", nil)

	if result.Status != shim.OK {
		test.FailNow()
	}
	return stub
}

func Invoke(test *testing.T, stub *shim.MockStub, function string, args ...string) []byte {
	ccArgs := make([][]byte, 1+len(args))
	ccArgs[0] = []byte(function)

	for i, arg := range args {
		ccArgs[i+1] = []byte(arg)
	}

	txID := rand.Int()
	result := stub.MockInvoke(strconv.Itoa(txID), ccArgs)

	fmt.Println("Call: 		 ", function, "(", strings.Join(args, ","), ")")
	fmt.Println("ResStatus:  ", result.Status)
	fmt.Println("ResMsg 	 ", result.Message)
	fmt.Println("ResPayload: ", string(result.Payload))
	fmt.Println()

	if result.Status != shim.OK {
		test.FailNow()
	}
	return result.Payload
}

func TestInvokeElectCC(test *testing.T) {

	stub := Init(test)
	stub.Invokables["elect_cc"] = shim.NewMockStub("elect_cc", new(elect_cc.ElectChaincode))

	Invoke(test, stub.Invokables["elect_cc"], "giveVote", "a", "b", "c", "d")

}

/*
Function List :
1	| registerUser
2	| registerElection
3	| registerCandidate
4	| registerVoter
5	| getUser
*/

func TestCCFunctions(test *testing.T) {
	stub := Init(test)

	var test_users []NewUser
	user := NewUser{}

	userKeys := make([]string, 0)
	userSSNs := make([]string, 0)

	fmt.Println("= Register User =")

	for i := 0; i < 2; i++ {
		ssn := fmt.Sprint("SSN_" + strconv.Itoa(rand.Int()))
		firstName := fmt.Sprint("FirstName" + strconv.Itoa(rand.Int()))
		lastName := fmt.Sprint("LastName" + strconv.Itoa(rand.Int()))

		res := Invoke(test, stub, "registerUser", ssn, firstName, lastName, "1992/02/24", "M")
		json.Unmarshal(res, &user)
		test_users = append(test_users, user)

		userSSNs = append(userSSNs, user.SSN)
		userKeys = append(userKeys, user.PublicKey)

	}

	fmt.Println("= User SSN = ")
	fmt.Println("User 0 SSN ", userSSNs[0])
	fmt.Println("User 1 SSN ", userSSNs[1])

	fmt.Println("= Get User By Key =")

	for i := 0; i < len(userSSNs); i++ {
		Invoke(test, stub, "getUser", "USERKEY", userKeys[i])
	}

	fmt.Println("= Get User By ID =")
	Invoke(test, stub, "getUser", "identity", userSSNs[0])

	fmt.Println("= Register Election =")
	Invoke(test, stub, "registerElection", "primary", "ElectionID", "2019/03/12", "2019/03/20")

	fmt.Println("= Register Candidate =")
	Invoke(test, stub, "registerCandidate", "primary", userSSNs[0], userKeys[0])

	fmt.Println("= Register Voter =")
	Invoke(test, stub, "registerVoter", userSSNs[1], "primary")

	fmt.Println("= Register Voter Candidate Case =")
	Invoke(test, stub, "registerVoter", userSSNs[0], "primary")

	fmt.Println("= Vote =")
	Invoke(test, stub, "vote", userSSNs[1], "primary", userKeys[1])

	fmt.Println("= Get User Info After Election =")
	Invoke(test, stub, "getUser", "identity", userSSNs[1])

	// @notice unit test for key history not implemented

}
