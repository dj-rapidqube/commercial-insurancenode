/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strconv"
    "time"
    //"strings"
    //"reflect"



    "github.com/hyperledger/fabric/core/chaincode/shim"
)

var userIndexStr = "_userindex"


type Claim struct {
    Id string `json:"userid"`
    ClaimNo int `json:"claimno"`
    ExaminerId string `json:"examinerid"`
    ClaimAdjusterId string `json:"claimadjusterid"`
    PublicAdjusterId string `json:"publicadjusterid"`
    Status string `json:"status"`
    Title string `json:"title"`
    DamageDetails string `json:"damagedetails"`
    TotalDamageValue int `json:"totaldamagevalue"`
    TotalClaimValue int `json:"totalclaimvalue"`
    Documents[] Document `json:"document"`
    ClaimNotifiedDate time.Time `json:"claimnotifieddate"`
    ClaimSubmittedDate time.Time `json:"claimsubmitteddate"`

    AssessedDamageValue int `json:"assesseddamagevalue"`
    AssessedClaimValue int `json:"assessedclaimvalue"`
    ClaimExaminedDate time.Time `json:"claimexamineddate"`
    ClaimValidatedDate time.Time `json:"claimvalidateddate"`
    Negotiationvalue[] Negotiation `json:"negotiationlist"`
    ApprovedClaim int `json:"approvedclaim"`
    ClaimApprovedDate time.Time `json:"claimapproveddate"`
    ClaimSettledDate time.Time `json:"claimsettleddate"`

}

type ClaimList struct {
    Claimlist[] Claim `json:"claimlist"` // contains array of claims
}

type Document struct {

    ClaimId int `json:"claimid"`
    FIRCopy string `json:"fircopy"` //the fieldtags of User Document hashvalue are needed to store in the ledger
    Photos string `json:"photos"`
    Certificates string `json:"certificates"`


}

type Negotiation struct {
    Id string `json:"id"`

    Negotiations int `json:"negotiationvalue"` //the fieldtags of claim Negotiation are needed to store in the ledger
    AsPerTerm2B string `json:"asperterm"`
}

type ExaminedUpdate struct {
    Id string `json:"id"`
    ClaimId int `json:"claimid"`
    AssessedDamageValue int `json:"assesseddamagevalue"` //the field tags of examiner
    AssessedClaimValue int `json:"assessedclaimvalue"`

}

type SimpleChaincode struct {}

// Main Function

func main() {
    err: = shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

// Init Function - reset all the things
func(t * SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {

    //_, args := stub.GetFunctionAndParameters()
    var Aval int
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    // Initialize the chaincode
    Aval, err = strconv.Atoi(args[0])
    if err != nil {
        return nil, errors.New("Expecting integer value for asset holding")
    }

    // Write the state to the ledger
    err = stub.PutState("abc", [] byte(strconv.Itoa(Aval))) //making a test var "abc" to read/write into ledger to test the network
    if err != nil {
        return nil, err
    }

    var empty[] string
    jsonAsBytes, _: = json.Marshal(empty) //marshal an emtpy array of strings to clear the index
    err = stub.PutState(userIndexStr, jsonAsBytes)
    if err != nil {
        return nil, err
    }


    return nil, nil
}

// Invoke is ur entry point to invoke a chaincode function
func(t * SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if
    function == "init" { //initialize the chaincode state, used as reset
        return t.Init(stub, "init", args)
    } else if
    function == "write" {
        return t.write(stub, args) //writes a value to the chaincode state

    } else if
    function == "notifyClaim" { //writes claim details with status notified in ledger
        return t.notifyClaim(stub, args)

    } else if
    function == "createClaim" { //writes  claim details with status approved in ledger
        return t.createClaim(stub, args)

    } else if
    function == "Delete" { //deletes an entity from its state
        return t.Delete(stub, args)

    } else if
    function == "UploadDocuments" { //upload the dcument hash value 
        return t.UploadDocuments(stub, args)

    } else if
    function == "ExamineClaim" { //Examine and updtaes the claim with status examined
        return t.ExamineClaim(stub, args)

    } else if
    function == "ClaimNegotiation" { //claim negotiations takes place between public adjuster and claim adjuster
        return t.ClaimNegotiation(stub, args)

    } else if
    function == "approveClaim" { //after negotiation claim amount is finalised and approved
        return t.approveClaim(stub, args)

    } else if
    function == "settleClaim" { //after negotiation claim amount is finalised and approved
        return t.settleClaim(stub, args)

    }

    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}

// write - invoke function to write key/value pair
func(t * SimpleChaincode) write(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var key, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }

    key = args[0]
    value = args[1]
    err = stub.PutState(key, [] byte(value)) //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}

// Query is our entry point for queries
func(t * SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if
    function == "readuser" { //read values for particular keys
        return t.readuser(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair

func(t * SimpleChaincode) readuser(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var name, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
    }

    name = args[0]
    valAsbytes, err: = stub.GetState(name) //get the key value from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil //send it onward
}

func(t * SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
        if len(args) != 1 {
            return nil, errors.New("Incorrect number of arguments. Expecting 1")
        }

        name: = args[0]
        err: = stub.DelState(name) //remove the key from chaincode state
        if err != nil {
            return nil, errors.New("Failed to delete state")
        }

        //get the user index
        userAsBytes, err: = stub.GetState(userIndexStr)
        if err != nil {
            return nil, errors.New("Failed to get array index")
        }
        var userIndex[] string
        json.Unmarshal(userAsBytes, & userIndex) //un stringify it aka JSON.parse()

        //remove user from index
        for i, val: = range userIndex {
            fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
            if val == name { //find the correct index

                userIndex = append(userIndex[: i], userIndex[i + 1: ]...) //remove it
                for x: = range userIndex { //debug prints...
                    fmt.Println(string(x) + " - " + userIndex[x])
                }
                break
            }
        }
        jsonAsBytes, _: = json.Marshal(userIndex) //save new index
        err = stub.PutState(userIndexStr, jsonAsBytes)
        return nil, nil
    }
    //notification of claim from insured takes place
func(t * SimpleChaincode) notifyClaim(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
        var err error


        if len(args) != 4 {
            return nil, errors.New("Incorrect number of arguments. Expecting 4")
        }

        //input sanitation
        fmt.Println("- start NotifyClaim")
        if len(args[0]) <= 0 {
            return nil, errors.New("1st argument must be a non-empty string")
        }
        if len(args[1]) <= 0 {
            return nil, errors.New("2nd argument must be a non-empty string")
        }
        if len(args[2]) <= 0 {
            return nil, errors.New("3rd argument must be a non-empty string")
        }
        if len(args[3]) <= 0 {
            return nil, errors.New("4th argument must be a non-empty string")
        }


        claim: = Claim {}
        claim.Id = args[0]

        claim.ClaimNo, err = strconv.Atoi(args[1])
        if err != nil {
            return nil, errors.New("Failed to get ClaimNo as cannot convert it to int")
        }

        claim.Title = args[2]
        claim.DamageDetails = args[3]
        claim.ClaimNotifiedDate = time.Now()
        claim.Status = "Notified"

        fmt.Println("claim", claim)
            //get claims empty[]
        UserAsBytes, err: = stub.GetState("getclaims")
        if err != nil {
            return nil, errors.New("Failed to get claims")
        }
        var claimlist ClaimList
        json.Unmarshal(UserAsBytes, & claimlist) //un stringify it aka JSON.parse()

        claimlist.Claimlist = append(claimlist.Claimlist, claim);
        fmt.Println("campaignallusers", claimlist.Claimlist) //append each claim to claimlist[]
        fmt.Println("! appended cuser to campaignallusers")
        jsonAsBytes, _: = json.Marshal(claimlist)
        fmt.Println("json", jsonAsBytes)
        err = stub.PutState("getclaims", jsonAsBytes) //rewrite claimlist[]
        if err != nil {
            return nil, err
        }
        fmt.Println("- end claimlist")
        return nil, nil
    }
    //application of claim from insuured takes place after notification
func(t * SimpleChaincode) createClaim(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error


    if len(args) != 4 {
        return nil, errors.New("Incorrect number of arguments. Expecting 4")
    }

    //input sanitation
    fmt.Println("- start createClaim")
    if len(args[0]) <= 0 {
        return nil, errors.New("1st argument must be a non-empty string")
    }
    if len(args[1]) <= 0 {
        return nil, errors.New("2nd argument must be a non-empty string")
    }
    if len(args[2]) <= 0 {
        return nil, errors.New("1st argument must be a non-empty string")
    }
    if len(args[3]) <= 0 {
        return nil, errors.New("1st argument must be a non-empty string")
    }




    ClaimId, err: = strconv.Atoi(args[0])
    if err != nil {
        return nil, errors.New("Failed to get ClaimId as cannot convert it to int")
    }
    TotalDamageValue, err: = strconv.Atoi(args[1])
    if err != nil {
        return nil, errors.New("Failed to get TotalDamageValue as cannot convert it to int")
    }
    TotalClaimValue, err: = strconv.Atoi(args[2])
    if err != nil {
        return nil, errors.New("Failed to get TotalClaimValue as cannot convert it to int")
    }
    Status: = "Submitted"
    PublicAdjusterId: = args[3]
    ClaimSubmittedDate: = time.Now()


    UserAsBytes, err: = stub.GetState("getclaims")
    if err != nil {
        return nil, errors.New("Failed to get claims")
    }

    var claimlist ClaimList
    json.Unmarshal(UserAsBytes, & claimlist) //un stringify it aka JSON.parse()


    for i: = 0;
    i < len(claimlist.Claimlist);
    i++{


        if (claimlist.Claimlist[i].ClaimNo == ClaimId) {

            claimlist.Claimlist[i].TotalDamageValue = TotalDamageValue
            claimlist.Claimlist[i].TotalClaimValue = TotalClaimValue

            claimlist.Claimlist[i].Status = Status
            claimlist.Claimlist[i].PublicAdjusterId = PublicAdjusterId
            claimlist.Claimlist[i].ClaimSubmittedDate = ClaimSubmittedDate
        }


        jsonAsBytes, _: = json.Marshal(claimlist)
        fmt.Println("json", jsonAsBytes)
        err = stub.PutState("getclaims", jsonAsBytes)
        if err != nil {
            return nil, err
        }
    }
    fmt.Println("- end claimlist")
    return nil, nil
}

//upload documents of insured in form of hash takes place			

func(t * SimpleChaincode) UploadDocuments(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
        var err error


        if len(args) != 4 {
            return nil, errors.New("Incorrect number of arguments. Expecting 4")
        }

        //input sanitation
        fmt.Println("- start registration")
        if len(args[0]) <= 0 {
            return nil, errors.New("1st argument must be a non-empty string")
        }
        if len(args[1]) <= 0 {
            return nil, errors.New("2nd argument must be a non-empty string")
        }
        if len(args[2]) <= 0 {
            return nil, errors.New("3rd argument must be a non-empty string")
        }
        if len(args[3]) <= 0 {
            return nil, errors.New("4th argument must be a non-empty string")
        }


        document: = Document {}

        document.ClaimId, err = strconv.Atoi(args[0])
        if err != nil {
            return nil, errors.New("Failed to get ClaimId as cannot convert it to int")
        }
        document.FIRCopy = args[1]

        document.Photos = args[2]
        document.Certificates = args[3]


        fmt.Println("document", document)

        UserAsBytes, err: = stub.GetState("getclaims")
        if err != nil {
            return nil, errors.New("Failed to get claims")
        }

        var claimlist ClaimList
        json.Unmarshal(UserAsBytes, & claimlist) //un stringify it aka JSON.parse()


        for i: = 0;
        i < len(claimlist.Claimlist);
        i++{


            if (claimlist.Claimlist[i].ClaimNo == document.ClaimId) {

                claimlist.Claimlist[i].Documents = append(claimlist.Claimlist[i].Documents, document);

            }


            jsonAsBytes, _: = json.Marshal(claimlist)
            fmt.Println("json", jsonAsBytes)
            err = stub.PutState("getclaims", jsonAsBytes)
            if err != nil {
                return nil, err
            }
        }

        fmt.Println("- end uploaddocumen")
        return nil, nil
    }
    //examination of claim takes place from examiner
func(t * SimpleChaincode) ExamineClaim(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
        var err error


        if len(args) != 4 {
            return nil, errors.New("Incorrect number of arguments. Expecting 4")
        }

        //input sanitation
        fmt.Println("- start ExamineClaim")
        if len(args[0]) <= 0 {
            return nil, errors.New("1st argument must be a non-empty string")
        }
        if len(args[1]) <= 0 {
            return nil, errors.New("2nd argument must be a non-empty string")
        }
        if len(args[2]) <= 0 {
            return nil, errors.New("3rd argument must be a non-empty string")
        }
        if len(args[3]) <= 0 {
            return nil, errors.New("3rd argument must be a non-empty string")
        }


        examine: = ExaminedUpdate {}
        examine.Id = args[0]

        examine.ClaimId, err = strconv.Atoi(args[1])
        if err != nil {
            return nil, errors.New("Failed to get ClaimId as cannot convert it to int")
        }
        examine.AssessedDamageValue, err = strconv.Atoi(args[2])
        if err != nil {
            return nil, errors.New("Failed to get AssessedDamageValue as cannot convert it to int")
        }
        examine.AssessedClaimValue, err = strconv.Atoi(args[3])
        if err != nil {
            return nil, errors.New("Failed to get AssessedClaimValue as cannot convert it to int")
        }
        Status: = "Examined"
        ClaimExaminedDate: = time.Now()


        fmt.Println("examine", examine)

        UserAsBytes, err: = stub.GetState("getclaims")
        if err != nil {
            return nil, errors.New("Failed to get claims")
        }

        var claimlist ClaimList
        json.Unmarshal(UserAsBytes, & claimlist) //un stringify it aka JSON.parse()


        for i: = 0;
        i < len(claimlist.Claimlist);
        i++{


            if (claimlist.Claimlist[i].ClaimNo == examine.ClaimId) {

                claimlist.Claimlist[i].AssessedDamageValue = examine.AssessedDamageValue
                claimlist.Claimlist[i].AssessedClaimValue = examine.AssessedClaimValue
                claimlist.Claimlist[i].Status = Status
                claimlist.Claimlist[i].ExaminerId = examine.Id
                claimlist.Claimlist[i].ClaimExaminedDate = ClaimExaminedDate
            }


            jsonAsBytes, _: = json.Marshal(claimlist)
            fmt.Println("json", jsonAsBytes)
            err = stub.PutState("getclaims", jsonAsBytes)
            if err != nil {
                return nil, err
            }
        }

        fmt.Println("- end ExaminedDocument")
        return nil, nil
    }
    //claim negotiation between public adjuster and claim adjuster takes place
func(t * SimpleChaincode) ClaimNegotiation(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
        var err error


        if len(args) != 4 {
            return nil, errors.New("Incorrect number of arguments. Expecting 4")
        }

        //input sanitation
        fmt.Println("- start ClaimNegotiation")
        if len(args[0]) <= 0 {
            return nil, errors.New("1st argument must be a non-empty string")
        }
        if len(args[1]) <= 0 {
            return nil, errors.New("2nd argument must be a non-empty string")
        }
        if len(args[2]) <= 0 {
            return nil, errors.New("3rd argument must be a non-empty string")
        }
        if len(args[3]) <= 0 {
            return nil, errors.New("3rd argument must be a non-empty string")
        }



        negotiation: = Negotiation {}
        negotiation.Id = args[0]

        ClaimId, err: = strconv.Atoi(args[1])
        if err != nil {
            return nil, errors.New("Failed to get ClaimId as cannot convert it to int")
        }
        negotiation.Negotiations, err = strconv.Atoi(args[2])
        if err != nil {
            return nil, errors.New("Failed to get Negotiations as cannot convert it to int")
        }
        negotiation.AsPerTerm2B = args[3]

        Status: = "Validated"
        ClaimValidatedDate: = time.Now()



        fmt.Println("negotiation", negotiation)

        UserAsBytes, err: = stub.GetState("getclaims")
        if err != nil {
            return nil, errors.New("Failed to get claims")
        }

        var claimlist ClaimList
        json.Unmarshal(UserAsBytes, & claimlist) //un stringify it aka JSON.parse()


        for i: = 0;
        i < len(claimlist.Claimlist);
        i++{


            if (claimlist.Claimlist[i].ClaimNo == ClaimId) {
                if (claimlist.Claimlist[0].Negotiationvalue == nil) {
                    claimlist.Claimlist[i].Status = Status
                    claimlist.Claimlist[i].ClaimAdjusterId = negotiation.Id
                    claimlist.Claimlist[i].Negotiationvalue = append(claimlist.Claimlist[i].Negotiationvalue, negotiation);
                    claimlist.Claimlist[i].ClaimValidatedDate = ClaimValidatedDate
                } else {
                    claimlist.Claimlist[i].Status = Status
                    claimlist.Claimlist[i].Negotiationvalue = append(claimlist.Claimlist[i].Negotiationvalue, negotiation);
                    claimlist.Claimlist[i].ClaimValidatedDate = ClaimValidatedDate
                }
            }
            jsonAsBytes, _: = json.Marshal(claimlist)
            fmt.Println("json", jsonAsBytes)
            err = stub.PutState("getclaims", jsonAsBytes)
            if err != nil {
                return nil, err
            }
        }

        fmt.Println("- end Negotiation")
        return nil, nil
    }
    //after negotiation claim amount will be finalised and approved
func(t * SimpleChaincode) approveClaim(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error


    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    //input sanitation
    fmt.Println("- start approveClaim")
    if len(args[0]) <= 0 {
        return nil, errors.New("1st argument must be a non-empty string")
    }




    ClaimId, err: = strconv.Atoi(args[0])
    if err != nil {
        return nil, errors.New("Failed to get ClaimId as cannot convert it to int")
    }


    Status: = "Approved"
    ClaimApprovedDate: = time.Now()



    UserAsBytes, err: = stub.GetState("getclaims")
    if err != nil {
        return nil, errors.New("Failed to get claims")
    }

    var claimlist ClaimList
    json.Unmarshal(UserAsBytes, & claimlist) //un stringify it aka JSON.parse()


    for i: = 0;
    i < len(claimlist.Claimlist);
    i++{


        if (claimlist.Claimlist[i].ClaimNo == ClaimId) {
            if (claimlist.Claimlist[i].Negotiationvalue[(len(claimlist.Claimlist[i].Negotiationvalue) - 1)].Negotiations == claimlist.Claimlist[i].Negotiationvalue[(len(claimlist.Claimlist[i].Negotiationvalue) - 2)].Negotiations) {
                claimlist.Claimlist[i].Status = Status
                claimlist.Claimlist[i].ClaimApprovedDate = ClaimApprovedDate
                lastindex: = (len(claimlist.Claimlist[i].Negotiationvalue) - 1)
                lastnegotiation: = claimlist.Claimlist[i].Negotiationvalue[lastindex]
                claimlist.Claimlist[i].ApprovedClaim = lastnegotiation.Negotiations


            }
        }
        jsonAsBytes, _: = json.Marshal(claimlist)
        fmt.Println("json", jsonAsBytes)
        err = stub.PutState("getclaims", jsonAsBytes)
        if err != nil {
            return nil, err
        }
    }

    fmt.Println("- end approve claim")
    return nil, nil
}


func(t * SimpleChaincode) settleClaim(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error


    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    //input sanitation
    fmt.Println("- start settleClaim")
    if len(args[0]) <= 0 {
        return nil, errors.New("1st argument must be a non-empty string")
    }




    ClaimId, err: = strconv.Atoi(args[0])
    if err != nil {
        return nil, errors.New("Failed to get ClaimId as cannot convert it to int")
    }


    Status: = "Settled"
    ClaimSettledDate: = time.Now()



    UserAsBytes, err: = stub.GetState("getclaims")
    if err != nil {
        return nil, errors.New("Failed to get claims")
    }

    var claimlist ClaimList
    json.Unmarshal(UserAsBytes, & claimlist) //un stringify it aka JSON.parse()




    for i: = 0;
    i < len(claimlist.Claimlist);
    i++{


        if (claimlist.Claimlist[i].ClaimNo == ClaimId) {


            claimlist.Claimlist[i].Status = Status

            claimlist.Claimlist[i].ClaimSettledDate = ClaimSettledDate
        }
        jsonAsBytes, _: = json.Marshal(claimlist)
        fmt.Println("json", jsonAsBytes)
        err = stub.PutState("getclaims", jsonAsBytes)
        if err != nil {
            return nil, err
        }
    }

    fmt.Println("- end settled claim")
    return nil, nil
}