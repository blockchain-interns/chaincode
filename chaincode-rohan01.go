//First Custom Chaincode
//Transactions require two values: Age and Name

package main

//Imports
import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//Chaincode Structure
type SimpleChaincode struct {
}

//Init function (runs when chaincode is initialized)
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mycc-new Init")
	_, args := stub.GetFunctionAndParameters()
	var age, name string    // VARIABLE (ENTITY) NAMES
	var ageVal int          //AGE
	var nameVal string      //NAME
	var err error

	//Ensure 4 arguments are passed into Init function
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4.")
	}

	// Initialize the chaincode
	//Define "age" variable name
	age = args[0]

	//Check if value passed for ageVal is an int
	ageVal, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	//Define "name" variable name
	name = args[2]

	//Set value for nameVal
	nameVal = args[3]

	//Print set values to CLI
	fmt.Printf("ageVal = %d, nameVal = %s\n", ageVal, nameVal)

	// Write the state to the ledger
	err = stub.PutState(age, []byte(strconv.Itoa(ageVal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(name, []byte(nameVal))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//Invoke function defines what to do when various functions are invoked
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("mycc-new Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Assigns new values to ageVal and nameVal
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// Gets entity states from the ledger
		return t.query(stub, args)
	}

	//Check if function invoked is valid
	return shim.Error("Invalid function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// "invoke" function sets ageVal and nameVal to user-defined values
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var age, name string    // VARIABLE (ENTITY) NAMES
	var ageVal int		//AGE
	var nameVal string	//NAME
	var X int    		//New AGE value
	var N string 		//New NAME value
	var err error

	//Ensure 4 arguments are passed into function
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4. Did you include age and name?")
	}

	//Define new names for "age" and "name" variables
	age = args[0]
	name = args[1]

	// Get the state from the ledger (ensure not in "err" state)
	ageValbytes, err := stub.GetState(age)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if ageValbytes == nil {
		return shim.Error("Entity not found")
	}
	ageVal, _ = strconv.Atoi(string(ageValbytes))

	nameValbytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if nameValbytes == nil {
		return shim.Error("Entity not found")
	}
	nameVal = string(nameValbytes)

	// Set new values passed into function
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting an integer value...")
	}
	N = args[3]
	ageVal = X
	nameVal = N
	fmt.Printf("ageVal = %d, nameVal = %s\n", ageVal, nameVal)

	// Write the state back to the ledger
	err = stub.PutState(age, []byte(strconv.Itoa(ageVal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	stub.PutState(name, []byte(nameVal))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// "delete" function deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// "query" function gets the state of an entity from the ledger
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // ENTITY VARIABLE
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	//Set value passed into function as name of entity to check
	A = args[0]

	// Get the state from the ledger
	AValbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	//If entity does not exist in ledger...
	if AValbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "or entity does not exist in ledger" + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Value\":\"" + string(AValbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(AValbytes)
}

//Main() function
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
