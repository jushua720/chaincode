package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	c "../constants"
	msg "../msg"

	a "../access"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func ConvertToBytes(str []string) []byte {

	stringByte := "\x00" + strings.Join(str, "\x20\x00")

	return []byte(stringByte)
}

func ValidateAge(dob, dateFormat, startDate, endDate string, minAge int) (string, bool) {

	isAdult := true
	dateOfBirth, _ := time.Parse(dateFormat, dob)

	age := time.Now().Year() - dateOfBirth.Year()
	dobMonthDay, _ := strconv.Atoi(strconv.Itoa(int(dateOfBirth.Month())) + strconv.Itoa(dateOfBirth.Day()))
	nowMonthDay, _ := strconv.Atoi(strconv.Itoa(int(time.Now().Month())) + strconv.Itoa(time.Now().Day()))

	if dobMonthDay > nowMonthDay {
		age = age - 1
	}

	adultYear := dateOfBirth.Year() + 18

	adultDate, _ := time.Parse(dateFormat, strconv.Itoa(adultYear)+"/"+fmt.Sprintf("%02d", int(dateOfBirth.Month()))+"/"+strconv.Itoa(dateOfBirth.Day()))

	if age < 18 {
		isAdult = IsWithinRange(adultDate.Format(dateFormat), startDate, endDate, dateFormat)
	}

	return strconv.Itoa(age), isAdult
}

func IsWithinRange(dateToCheck, startDate, endDate, dateFormat string) bool {

	date, _ := time.Parse(dateFormat, dateToCheck)
	date1, _ := time.Parse(dateFormat, startDate)
	date2, _ := time.Parse(dateFormat, endDate)

	if date.Before(date2) && date.After(date1) || dateToCheck == startDate || dateToCheck == endDate {
		return true
	}

	return false
}

func GetAllCompositeKeys(stub shim.ChaincodeStubInterface, objType string, args []string) ([]string, error) {
	var foundKeys []string

	keySearchIterator, err := stub.GetStateByPartialCompositeKey(objType, args)
	if err != nil {
		return []string{}, errors.New(msg.GetErrMsg("COM_ERR_04", []string{objType, err.Error()}))
	}

	defer keySearchIterator.Close()

	if !keySearchIterator.HasNext() {
		return []string{}, nil
	}

	for keySearchIterator.HasNext() {
		keyRange, err := keySearchIterator.Next()
		if err != nil {
			return []string{}, errors.New(msg.GetErrMsg("COM_ERR_06", []string{err.Error()}))
		}

		foundKeys = append(foundKeys, keyRange.Key)
	}

	return foundKeys, nil
}

func FindCompositeKey(stub shim.ChaincodeStubInterface, objType string, args []string) (string, error) {

	var foundKey string

	keySearchIterator, err := stub.GetStateByPartialCompositeKey(objType, args)
	if err != nil {
		return "", errors.New(msg.GetErrMsg("COM_ERR_04", []string{objType, err.Error()}))
	}

	defer keySearchIterator.Close()

	if !keySearchIterator.HasNext() {
		return "", nil
	}

	for keySearchIterator.HasNext() {
		keyRange, err := keySearchIterator.Next()
		if err != nil {
			return "", errors.New(msg.GetErrMsg("COM_ERR_06", []string{err.Error()}))
		}

		foundKey = keyRange.Key
	}

	return foundKey, nil

}

func ValidateElectionPeriod(date1, date2 string) bool {

	startDate, err := time.Parse("2006/01/02", date1)
	endDate, err := time.Parse("2006/01/02", date2)
	if err != nil {
		return false
	}

	startMonthDay, _ := strconv.Atoi(strconv.Itoa(int(startDate.Month())) + strconv.Itoa(startDate.Day()))
	endMonthDay, _ := strconv.Atoi(strconv.Itoa(int(endDate.Month())) + strconv.Itoa(endDate.Day()))

	if (endMonthDay - startMonthDay) <= 0 {
		return false
	}

	return true
}

func ValidateArgument(arg string) bool {
	var rxPat = regexp.MustCompile(`^[M, m, Male, male, MALE, F, f, Female, female, FEMALE, O, o, Other, other, OTHER]{1}$`)

	if !rxPat.MatchString(arg) {
		return false
	}

	return true
}

func ArrayToChaincodeArgs(args []string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func ConvertStrToInt(args []string) ([]int, error) {
	res := make([]int, 0)
	for i := 0; i < len(args); i++ {
		val, err := strconv.Atoi(args[i])
		if err != nil {
			return []int{}, errors.New(msg.GetErrMsg("COM_ERR_15", []string{args[i], err.Error()}))
		}
		res = append(res, val)
	}

	return res, nil

}

func MarshalData(data string, dataStruct interface{}) ([]byte, error) {

	if err := json.Unmarshal([]byte(data), &dataStruct); err != nil {
		return []byte{0x00}, errors.New(msg.GetErrMsg("COM_ERR_02", []string{err.Error()}))
	}

	js, err := json.Marshal(dataStruct)
	if err != nil {
		return []byte{0x00}, errors.New(msg.GetErrMsg("COM_ERR_03", []string{err.Error()}))
	}

	return js, nil
}

func FindUserBySSN(stub shim.ChaincodeStubInterface, ssn string) (bool, string) {

	keyResultsIterator, err := stub.GetStateByPartialCompositeKey(c.SSNKEY, []string{ssn})
	if err != nil {
		return false, ""
	}
	defer keyResultsIterator.Close()

	if !keyResultsIterator.HasNext() {
		return false, ""
	}

	// Key: "userID-userPubKey
	var i int
	for i = 0; keyResultsIterator.HasNext(); i++ {
		responseRange, err := keyResultsIterator.Next()
		if err != nil {
			return false, ""
		}

		_, keyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return false, ""
		}
		return true, keyParts[1]

	}

	return false, ""
}

func CreateCompKey(stub shim.ChaincodeStubInterface, objType string, args []string) error {

	compositeKey, err := stub.CreateCompositeKey(objType, args)
	if err != nil {
		return errors.New(msg.GetErrMsg("COM_ERR_08", []string{objType, args[1], err.Error()}))
	}

	err = stub.PutState(compositeKey, []byte{0x00})
	if err != nil {
		return errors.New(msg.GetErrMsg("COM_ERR_09", []string{compositeKey, err.Error()}))
	}
	return nil
}

func GenerateKeys() (string, string, error) {

	keys, err := a.GenerateKeys()
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	return keys.PrivateKey, keys.PublicKey, nil
}

func VerifyUser(key, data, R, S, x, y string) (bool, string, error) {

	hash := a.GetHash(data)
	pubKey := a.GetPubKeyFromXY(x, y)

	if key != a.GenerateAccount(pubKey) {
		return false, "", errors.New(msg.GetErrMsg("COM_ERR_21", []string{key, x, y}))
	}

	isVerified := a.Verify(pubKey, hash, R, S)

	return isVerified, hash, nil
}
