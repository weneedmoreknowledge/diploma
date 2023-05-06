package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Asset struct {
  University      string `json:"University"`
  ID              string `json:"ID"`
  StudentName     string `json:"StudentName"`
  Credit          int    `json:"Credit"`
  GradStatus      bool   `json:"GradStatus"`
}

const (
 mspID        = "Org1MSP"
 cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com"
 certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
 keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
 tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
 peerEndpoint = "localhost:7051"
 gatewayPeer  = "peer0.org1.example.com"

 mspID2        = "Org2MSP"
 cryptoPath2   = "../../test-network/organizations/peerOrganizations/org2.example.com"
 certPath2     = cryptoPath2 + "/users/User1@org2.example.com/msp/signcerts/cert.pem"
 keyPath2      = cryptoPath2 + "/users/User1@org2.example.com/msp/keystore/"
 tlsCertPath2  = cryptoPath2 + "/peers/peer0.org2.example.com/tls/ca.crt"
 peerEndpoint2 = "localhost:9051"
 gatewayPeer2  = "peer0.org2.example.com"

 mspID3        = "Org3MSP"
 cryptoPath3   = "../../test-network/organizations/peerOrganizations/org3.example.com"
 certPath3     = cryptoPath3 + "/users/User1@org3.example.com/msp/signcerts/User1@org3.example.com-cert.pem"
 keyPath3      = cryptoPath3 + "/users/User1@org3.example.com/msp/keystore/"
 tlsCertPath3  = cryptoPath3 + "/peers/peer0.org3.example.com/tls/ca.crt"
 peerEndpoint3 = "localhost:11051"
 gatewayPeer3  = "peer0.org3.example.com"

 // Override default values for chaincode and channel name as they may differ in testing contexts.
 chaincodeName = "basic"
 channelName = "mychannel"
)

//couchDB   
// http://localhost:5984/_utils  
// http://localhost:7984/_utils
// http://localhost:9984/_utils

//admin
//adminpw

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	initLedger(contract) // Run initLedger initially.

	fmt.Print("Server is Running...\n")
	fmt.Print("\n http://localhost:5984/_utils \n")
	fmt.Print("\n http://localhost:7984/_utils \n")
	fmt.Print("\n http://localhost:9984/_utils \n")

	handleRoutes()
}

func handleRoutes()  {
	router:= mux.NewRouter()
	router.HandleFunc("/",initPage).Methods("GET")

	router.HandleFunc("/get",getPage).Methods("GET")
	router.HandleFunc("/get2",getPage2).Methods("GET")
	router.HandleFunc("/get3",getPage3).Methods("GET")
	
	router.HandleFunc("/createAsset",createPage).Methods("PUT")

	router.HandleFunc("/updateAsset",updatePage).Methods("PUT")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Minute, // Set the read timeout to 10 seconds
		WriteTimeout: 5 * time.Minute, // Set the write timeout to 10 seconds
	}

	log.Fatal(server.ListenAndServe())
}


func initPage(w http.ResponseWriter, r *http.Request)  {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	initLedger(contract) // Run initLedger initially.

	fmt.Println("\n\nComplete the initial")
}


// ----------------------------------------------------------------------------------

// ----------------------------------------------------------------------------------
//getPages----------------------------------------------------------------------------------
func getPage(w http.ResponseWriter, r *http.Request)  {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	result := getAllAssets(contract)

	w.Header().Set("Content-Type","application/json")

	w.Write(result)

	fmt.Println("\n\nComplete the GetAssets")
}
func getPage2(w http.ResponseWriter, r *http.Request)  {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection2()
	defer clientConnection.Close()

	id := newIdentity2()
	sign := newSign2()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	result := getAllAssets(contract)

	w.Header().Set("Content-Type","application/json")

	w.Write(result)

	fmt.Println("\n\nComplete the GetAssets")
}
func getPage3(w http.ResponseWriter, r *http.Request)  {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection3()
	defer clientConnection.Close()

	id := newIdentity3()
	sign := newSign3()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	result := getAllAssets(contract)

	w.Header().Set("Content-Type","application/json")

	w.Write(result)

	fmt.Println("\n\nComplete the GetAssets")
}

// ----------------------------------------------------------------------------------
//createPage ----------------------------------------------------------------------------------
func createPage(w http.ResponseWriter, r *http.Request)  {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()
	
	id := newIdentity()
	sign := newSign()
	
		// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()
	
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)
	
	type Req struct {
		University      string `json:"university"`
		ID              string `json:"id"`
		StudentName     string `json:"studentName"`
	}

	var asset Req
	err = json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		// Invalid request payload
		print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate request data
	if asset.ID == "" || asset.University==""||asset.StudentName=="" {
		// Required fields missing
		http.Error(w, "Name and University and ID are required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("post sent Successfully")
	
	createAsset(contract, asset.University, asset.ID, asset.StudentName, 0, false)
	
	fmt.Println("\n\nComplete the CreateAssets")
}

// ----------------------------------------------------------------------------------
//updatePage ----------------------------------------------------------------------------------
func updatePage(w http.ResponseWriter, r *http.Request)  {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	type Req struct {
		ID      		string  `json:"id"`
		Credit          string  `json:"credit"`
		GradStatus     	string  `json:"gradStatus"`
	}

	var asset Req
	err = json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("post sent Successfully")

	StudentId := asset.ID
	Credits,err := strconv.Atoi(asset.Credit)
	Graduate,err := strconv.ParseBool(asset.GradStatus)

	updateAsset(contract,StudentId,Credits,Graduate)

	fmt.Println("\n\nComplete the CreateAssets")
}


// ----------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------

// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func initLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}
	

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction to query ledger state.
func getAllAssets(contract *client.Contract) []byte {
	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
	
	return evaluateResult
}


// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createAsset(contract *client.Contract,university string, id string, studentName string, credit int, gradStatus bool) {
	Credits:=strconv.Itoa(credit)
	GradStatus:=strconv.FormatBool(gradStatus)

	_, err := contract.SubmitTransaction("CreateAsset", university,id,studentName,Credits,GradStatus)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}


// Evaluate a transaction by assetID to query ledger state.
func readAssetByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", assetId)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func updateAsset(contract *client.Contract, id string, credit int, gradStatus bool) {
	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, updates")

	Credits:=strconv.Itoa(credit)
	GradStatus:=strconv.FormatBool(gradStatus)

	submitResult, commit, err := contract.SubmitAsync("TransferAsset", client.WithArguments(id,Credits,GradStatus))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Mark. \n", string(submitResult))
	fmt.Println("*** Waiting for transaction commit.")

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}
}




// ----------------------------------------------------------------------------------
// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}
func newGrpcConnection2() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath2)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer2)

	connection, err := grpc.Dial(peerEndpoint2, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}
func newGrpcConnection3() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath3)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer3)

	connection, err := grpc.Dial(peerEndpoint3, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}


// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}
func newIdentity2() *identity.X509Identity {
	certificate, err := loadCertificate(certPath2)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID2, certificate)
	if err != nil {
		panic(err)
	}

	return id
}
func newIdentity3() *identity.X509Identity {
	certificate, err := loadCertificate(certPath3)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID3, certificate)
	if err != nil {
		panic(err)
	}

	return id
}


func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}


// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}
func newSign2() identity.Sign {
	files, err := os.ReadDir(keyPath2)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath2, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}
func newSign3() identity.Sign {
	files, err := os.ReadDir(keyPath3)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath3, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}


// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}