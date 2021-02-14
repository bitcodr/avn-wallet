package main

import (
	"log"
	"os"

	"github.com/bitcodr/avn-wallet/provider"
)

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	if err := os.Setenv("WALLET_SERVICE_ROOT_DIR", currentDir); err != nil {
		log.Fatalln(err)
	}

	provider.HTTP()

	provider.GRPC()
}
