package main

import (
	"os"

	"github.com/drossan/claude-init/cmd/root"
	"github.com/drossan/claude-init/internal/logger"
)

func main() {
	// Crear logger con salida a stdout y nivel INFO
	log := logger.New(os.Stdout, logger.INFOLevel)

	// Ejecutar el comando raíz
	if err := root.Execute(log); err != nil {
		os.Exit(1)
	}
}
