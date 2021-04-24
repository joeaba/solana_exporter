package collectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	HttpTimeout = 5 * time.Second
)

func GetKeys() ([]byte, error) {
	jsonFile, err := os.Open("//Users/anmolrajarora/Downloads/Solana_Home/solana_exporter/config.local.json")
	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened file as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	return byteValue, err
}
