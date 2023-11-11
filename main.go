package main

import (
    "bufio"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strings"

    "golang.org/x/crypto/ripemd160"
    "github.com/cosmos/btcutil/bech32"
)

// PublicKeyToAddress converts a public key to a bech32 Tendermint/Cosmos based address
func PublicKeyToAddress(addressPrefix string, pubKeyBytes []byte) string {
    // Hash pubKeyBytes as: RIPEMD160(SHA256(public_key_bytes))
    pubKeySha256Hash := sha256.Sum256(pubKeyBytes)
    ripemd160hash := ripemd160.New()
    ripemd160hash.Write(pubKeySha256Hash[:])
    addressBytes := ripemd160hash.Sum(nil)

    // Convert addressBytes into a bech32 string using the provided addressPrefix
    address, err := toBech32(addressPrefix, addressBytes)
    if err != nil {
        panic(err)
    }

    return address
}

// toBech32 converts a byte slice into a bech32 encoded string with the given prefix
func toBech32(addrPrefix string, addrBytes []byte) (string, error) {
    converted, err := bech32.ConvertBits(addrBytes, 8, 5, true)
    if err != nil {
        return "", err
    }

    addr, err := bech32.Encode(addrPrefix, converted)
    if err != nil {
        return "", err
    }

    return addr, nil
}

// decodePublicKeyString decodes a public key string into a byte slice based on the given format
func decodePublicKeyString(pubKey string, format string) ([]byte, error) {
    var pubKeyBytes []byte
    var err error

    switch format {
    case "Ed25519":
        pubKeyBytes, err = base64.StdEncoding.DecodeString(pubKey)
    case "Secp256k1":
        pubKeyBytes, err = hex.DecodeString(pubKey)
    // Add cases for other formats as needed
    default:
        err = fmt.Errorf("unsupported public key format")
    }

    return pubKeyBytes, err
}

// JSONData structure to match the JSON file
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
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Choose operation mode (1 - Process JSON file, 2 - Convert single address): ")
    mode, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input:", err)
        os.Exit(1)
    }
    mode = strings.TrimSpace(mode)

    switch mode {
    case "1":
        fmt.Print("Enter the path to your JSON file: ")
        jsonFilePath, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input:", err)
            os.Exit(1)
        }
        jsonFilePath = strings.TrimSpace(jsonFilePath)

        jsonData, err := ioutil.ReadFile(jsonFilePath)
        if err != nil {
            fmt.Println("Error reading JSON file:", err)
            os.Exit(1)
        }

        var data JSONData
        err = json.Unmarshal(jsonData, &data)
        if err != nil {
            fmt.Println("Error unmarshalling JSON data:", err)
            os.Exit(1)
        }

        for _, validator := range data.Result.Validators {
            pubKeyValue := validator.PubKey.Value
            pubKeyBytes, err := decodePublicKeyString(pubKeyValue, "Ed25519") // Assuming Ed25519 format
            if err != nil {
                fmt.Println("Error decoding public key:", err)
                continue
            }
            bech32address := PublicKeyToAddress("nomic", pubKeyBytes)
            fmt.Println(bech32address)
        }

    case "2":
        fmt.Print("Enter the public key: ")
        publicKey, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input:", err)
            os.Exit(1)
        }
        publicKey = strings.TrimSpace(publicKey)

        fmt.Print("Enter the public key format (Ed25519, Secp256k1, etc.): ")
        keyFormat, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input:", err)
            os.Exit(1)
        }
        keyFormat = strings.TrimSpace(keyFormat)

        fmt.Print("Enter the address prefix: ")
        addressPrefix, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input:", err)
            os.Exit(1)
        }
        addressPrefix = strings.TrimSpace(addressPrefix)

        pubKeyBytes, err := decodePublicKeyString(publicKey, keyFormat)
        if err != nil {
            fmt.Println("Error decoding public key:", err)
            os.Exit(1)
        }
        bech32address := PublicKeyToAddress(addressPrefix, pubKeyBytes)
        fmt.Println(bech32address)

    default:
        fmt.Println("Invalid mode selected")
        os.Exit(1)
    }
}

