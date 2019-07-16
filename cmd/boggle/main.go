package main

import (
	"os"

	"github.com/adambullmer/go-boggle/cmd/boggle/internal/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := handlers.Run(); err != nil {
		log.Info("Emergency Shut Down", err)
		os.Exit(1)
	}
}
