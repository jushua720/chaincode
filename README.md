# Voting Chaincode

The code is created to demonstrate GoLang chaincode functionality and not intended for production. 

&nbsp; 

## Detailed Information on Implemented Functions


### 1. registerUser

| Arguments         		| Payload                | 
| :---            		| :----                  | 
| [0] : UserSSN     		| [0]: UserSSN           | 
| [1] : FirstName  		| [1]: PublicKey         | 
| [2] : LastName    		| [2]: PrivateKey        | 
| [3] : DateOfBirth <br> [ *yyyy/mm/dd* ] | [3]: RegistrationDate     | 
| [4] : Gender <br> [ *M, m, Male, MALE; F, f, Female, FEMALE, O, o, Other, other, OTHER* ]   |      | 


&nbsp; 

Function contains calls to the following sub-functions and methods:

| Function              | Description  |
| :-----                | :-----        | 
| FindUserBySSN()       | Implements *GetStateByPartialCompositeKey* method  | 
| ValidateArgument()    | Checks whether provided argument matches a pattern |
| GenerateKeys()        | Generates ECDSA public and private keys |
| GenerateAccount()     | Shortens ECDSA public key making it 40 characters in length.                              <br> Purpose: to save memory | 
| CreateCompKey()       | Demonstrates composite key creation | 
| MarshalData())        | Demonstrates a way of passing a data struct as a parameter | 


&nbsp; 

### 2. registerElection

| Arguments         		| Payload                | 
| :---            		| :----                  | 
| [0] : ElectionType <br> [ *primary / general / local* ]| [0]: electionType           | 
| [1] : ElectionID  		| [1]: electionID         | 
| [2] : ElectionStartDate <br>   [ *yyyy/mm/dd* ]        | [2]: startDate        | 
| [3] : ElectionEndDate <br> [ *yyyy/mm/dd* ] | [3]: endDate     | 
| [4] : TodayDate <br> [ *yyyy/mm/dd* ]   |   [4]: txID   | 


&nbsp; 

Function contains calls to the following sub-functions and methods:

| Function | Decription |
| :-----  | :----- | 
|FindCompositeKey()  | Implements *GetStateByPartialCompositeKey* method | 
|ValidateElectionPeriod()  | Ensures election lasts more than a day |

&nbsp; 

### 3. registerCandidate

| Arguments         		| Payload                | 
| :---            		    | :----                  | 
| [0] : ElectionType <br> [ *primary / general / local* ]| [0]: UserSSN             | 
| [1] : UserPublicKey  		| [1]: UserPublicKey       | 
| [2] : R                   | [2]: UserFirstName       | 
| [3] : S                   | [3]: UserLastName        | 
| [4] : X                   | [4]: UserDateOfBirth     | 
| [5] : Y                   | [5]: UserAge             | 
|                           | [6]: ElectionType        | 
|                           | [7]: ElectionPeriod      | 
|                           | [8]: TxID                | 

*R, S, X, Y – ecdsa algorithm parameters. Use [ ssilka ]  to generate it*

&nbsp; 

Function contains calls to the following sub-functions and methods:

| Function | Decription     |
| :-----   | :-----         | 
|VerifyUser()         | Constructs ecdsa user public key from X, Y. Verifies ecdsa signature using R, S, public key | 
|SplictCompositeKey()  | [**built-in**] Splits composite keys into attributes. |
|ValidateAge()   | Calculates user age and checks if a user is an adult ( *at least 18 years old* ) |



&nbsp; 

### 4. registerVoter

| Arguments | Payload |
| :-----  | :-----  | 
|[0] : UserSSN  | [0] : UserSSN | 
|[1] : ElectionType <br> [ *primary / general / local* ]  | [1] : UserFirstName |
|   | [2] : UserLastName |
|   | [3] : UserDateOfBirth| 
|   | [4] : UserAge | 
|   | [5] : UserEligibilityToVote  [ *bool* ] |
|   | [6] : IsVoterCandidate [ *bool* ] | 
|   | [7] : ElectionType | 
|   | [8] : ElectionPeriod <br> [ *yyyy/mm/dd-yyyy/mm/dd* ] | 

&nbsp; 

Function contains calls to the following sub-functions and methods:

| Function | Decription     |
| :-----   | :-----         | 
|Contains()         | [ **built-in** ] Checks if the string contains given substring | 





&nbsp; 

### 5. vote

| Arguments | Payload  |
| :-----  | :-----  | 
| [0] : UserSSN                   | [0] : UserSSN |
| [1] : ElectionType <br> [ *primary / general / local* ]  | [1] : UserFirstName |
| [2] : UserPublicKey             | [2] : UserLastName | 
|                                 | [3] : UserAge | 
|                                 | [4] : UserPublicKey |
|                                 | [5] : TodayDate     |
|                                 | [6] : ElectionType | 
|                                 | [7] : TxID | 

&nbsp; 

Function contains calls to the following sub-functions and methods:


| Function              | Description  |
| :-----                | :-----        | 
| FindUserBySSN()       | Implements *GetStateByPartialCompositeKey* method  | 
| ValidateArgument()    | Checks whether the provided argument matches the pattern |
| GenerateKeys()        | Generates ECDSA public and private keys |
| GenerateAccount()     | Shortens ECDSA public key making it 40 characters in length.                              <br> Purpose: save memory | 
| CreateCompKey()       | Demonstrates composite key creation | 
| MarshalData())        | Demonstrates a way of passing a data struct as a parameter | 




&nbsp; 

### 6. countVotes

| Arguments | Payload  |
| :-----  | :-----  | 
| [0] : ElectionType <br>  [ *primary / general / local* ]  | [0] : VotingResult |


&nbsp; 

Function contains calls to the following sub-functions and methods:

| Function | Decription |
| :-----  | :----- | 
|callOtherCC()  | Implements methid to call other chaincode | 
|append  | [ **built-in** ] Used to concatenate two slices |


&nbsp; 

### 7. getUser

| Arguments | Payload  |
| :-----  | :-----  | 
| [0] : SearchCriteria <br>  [ *identity / publickey* ]  | [0] : UserSSN |
| [1] : UserSSN **or** UserPublicKey  | [1] : UserPublicKey |
|                                 | [2] : UserFirstName | 
|                                 | [3] : UserLastName | 
|                                 | [4] : UserDateOfBirth |
|                                 | [5] : UserGender |
|                                 | [6] : UserElectionInfo | 
|                                 | [7] : UserRegistrationDate | 


&nbsp; 

### 8. getUserVotingHistory

| Arguments | Payload |
| :-----  | :-----  | 
|[0] : UserSSN  | [0] : UserSSN | 
|   | [1] : UserPublicKey |
|   | [2] : UserFirstName |
|   | [3] : UserLastName | 
|   | [4] : UserDateOfBirth | 
|   | [5] : UserGender |
|   | [6] : UserElectionInfo | 
|   | [7] : UserRegistrationDate | 

&nbsp; 

Function contains calls to the following sub-functions and methods:

| Function | description  |
| :-----   | :-----       | 
| GetHistoryForKey()  | [**built-in**] Return a history of key values |



&nbsp; 

### 9. getAllUsers

| Arguments | Payload |
| :-----  | :-----  | 
|[0] : Bookmark  | [0] : UsersPublicKeys | 
| [1] : PageSize  |  |
| 
