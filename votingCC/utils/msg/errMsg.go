package msg

import (
	"fmt"
)

var errCodeMap = map[string]string{

	"COM_ERR_01": "Incorrect Number of Arguments for \"%s\" Expected %s",

	"COM_ERR_02": "Failed to Unmarshal : %s",
	"COM_ERR_03": "Failed to Marshal : %s",

	"COM_ERR_04": "GetStateByPartialCompositeKey Failed : %s",
	"COM_ERR_05": "Composite Key HasNext() Failed : %s",
	"COM_ERR_06": "Composite Key Next() Failed : %s",
	"COM_ERR_07": "Composit Key Split Failed : %s",
	"COM_ERR_08": "Failed to Create Composite Key \"%s\" for \"%s\": %s",

	"COM_ERR_09": "Put State for \"%s\" Failed : %s",
	"COM_ERR_10": "Get State for \"%s\" Failed : %s",

	"COM_ERR_11": "Invalid Function Name : \"%s \"",

	"COM_ERR_12": "Invalid Query Type \"%s\" expercting \"%s\"",

	"COM_ERR_13": "Failed to iterate : %s",
	"COM_ERR_14": "User SSN \"%s\" does not exists",

	"COM_ERR_15": "Failed to convert \"%s\" : %s",
	"COM_ERR_16": "GetStateByRangeWithPagination Failed : %s",
	"COM_ERR_17": "Error calling chaincode \"%s\" : %s",

	"COM_ERR_18": "Invalid Argument \"%s\" : %s",
	"COM_ERR_19": "Failed To Get History for \"%s\" : %s",
	"COM_ERR_20": "Failed to Parse Value \"%s\" : %s",

	"VOT_ERR_01": "Duplicated SSN : \"%s\"",
	"VOT_ERR_02": "Failed to Register New User : %s",
	"VOT_ERR_03": "Failed to Generate Keys : %s",
	"VOT_ERR_04": "Invalid Election Type: \"%s\"",
	"VOT_ERR_05": "Invalid Election Period: \"%s\" - \"%s\"",
	"VOT_ERR_06": "Election Exists: %s",
	"VOT_ERR_07": "Election \"%s\" is Not Available for Registration",
	"VOT_ERR_08": "Keys mismatch: \"%s\"  and  \"%s\"",
	"VOT_ERR_09": "Candidate is already registered : \"%s\"",
	"VOT_ERR_10": "Voter \"%s\" is already registered",
	"VOT_ERR_11": "Not allowed to Vote : %s",
	"VOT_ERR_12": "Invalid Candidate \"%s\" : %s",
	"VOT_ERR_13": "%s Not \"%s\" Election Period : %s",
	"VOT_ERR_14": "%s Has Already Voted",
	"VOT_ERR_15": "Election \"%s\" Not Exist",
	"VOT_ERR_16": "Invalid Voting Method : %s",
	"VOT_ERR_17": "Election \"%s\" is Not Over: %s , %s ",

	"ELECT_ERR_01": "GetStateByPartialCompositeKeyWithPagination Failed : %s",
}

func GetErrMsgParams(arr []string) []interface{} {
	params := make([]interface{}, len(arr))
	for i, p := range arr {
		params[i] = p
	}

	return params
}

func GetErrMsg(msgCode string, params []string) string {
	msgMap := errCodeMap
	msgBody := msgMap[msgCode]
	msgParam := GetErrMsgParams(params)

	msg := fmt.Sprintf(msgBody, msgParam...)

	return msg
}
