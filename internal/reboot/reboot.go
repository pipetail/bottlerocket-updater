package reboot

import (
	"github.com/pipetail/bottlerocket-updater/internal/common"
	"github.com/pipetail/bottlerocket-updater/pkg/bottlerocket"
	"log"
	"os"
)

type Config struct {
	SocketPath string
}

func RealMain(config Config) {
	// prepare HTTP client with the special UDS configuration
	client := common.GetHTTPClient(config.SocketPath)
	err := bottlerocket.Reboot(client)
	if err != nil {
		log.Printf("could not reboot OS: %s", err.Error())
		os.Exit(1)
	}

	log.Println("rebbot request sent")
	os.Exit(0)
}
