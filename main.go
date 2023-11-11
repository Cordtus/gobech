package main

import (
    "bufio"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"

    "golang.org/x/crypto/ripemd160"
    "github.com/cosmos/btcutil/bech32"
)

// PublicKeyToAddress converts secp256k1 public key to a bech32 Tendermint/Cosmos based address
func PublicKeyToAddress(addressPrefix, publicKeyString string) string {
	// Decode public key string
	pubKeyBytes := decodePublicKeyString(publicKeyString)

	// Hash pubKeyBytes as: RIPEMD160(SHA256(public_key_bytes))
	pubKeySha256Hash := sha256.Sum256(pubKeyBytes)
	ripemd160hash := ripemd160.New()
	ripemd160hash.Write(pubKeySha256Hash[:])
	addressBytes := ripemd160hash.Sum(nil)

	// Convert addressBytes into a bech32 string
	address := toBech32("nomic", addressBytes)

	return address
}

// Code courtesy: https://github.com/cosmos/cosmos-sdk/blob/90c9c9a9eb4676d05d3f4b89d9a907bd3db8194f/types/bech32/bech32.go#L10
func toBech32(addrPrefix string, addrBytes []byte) string {
	converted, err := bech32.ConvertBits(addrBytes, 8, 5, true)
	if err != nil {
		panic(err)
	}

	addr, err := bech32.Encode(addrPrefix, converted)
	if err != nil {
		panic(err)
	}

	return addr
}

// decodePublicKeyString decodes a base-64 encoded public key
// into a Byte Array. The logic will differ for other string encodings
func decodePublicKeyString(pubKey string) []byte {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}

	return pubKeyBytes
}


// JSONData structure to match the updated JSON file
type JSONData struct {
    Result struct {
        Validators []struct {
            PubKey struct {
                Type  string `json:"type"`
                Value string `json:"value"`
            } `json:"pub_key"`
        } `json:"validators"`
    } `json:"result"`
}

func main() {
    // Create a new reader, assuming stdin is a terminal
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the path to your JSON file: ")
    jsonFilePath, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input:", err)
        os.Exit(1)
    }

    // Trim the newline character from the input
    jsonFilePath = strings.TrimSpace(jsonFilePath)

    // Read JSON file
    jsonData, err := ioutil.ReadFile(jsonFilePath)
    if err != nil {
        fmt.Println("Error reading JSON file:", err)
        os.Exit(1)
    }

    // Unmarshal JSON data
    var data JSONData
    err = json.Unmarshal(jsonData, &data)
    if err != nil {
        fmt.Println("Error unmarshalling JSON data:", err)
        os.Exit(1)
    }

    // Process each validator's public key
    for _, validator := range data.Result.Validators {
        pubKeyValue := validator.PubKey.Value
        bech32address := PublicKeyToAddress("nomic", pubKeyValue)
        fmt.Println(bech32address)
    }
}
