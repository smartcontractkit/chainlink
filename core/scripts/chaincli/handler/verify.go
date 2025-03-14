package handler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func (k *Keeper) changeToContractsDirectory() error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Check if hardhat.config.ts exists in the current directory, return if it does
	if _, err := os.Stat(filepath.Join(currentDir, "hardhat.config.ts")); err == nil {
		return nil
	}

	// Command should run from core/scripts/chaincli, so we need to change directory to contracts
	// Calculate the absolute path of the target directory
	absPath := filepath.Join(currentDir, "../../../contracts")

	// Change directory
	if err := os.Chdir(absPath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	// Check if hardhat.config.ts exists in the current directory
	if _, err := os.Stat(filepath.Join(absPath, "hardhat.config.ts")); err != nil {
		return errors.New("hardhat.config.ts not found in the current directory")
	}

	log.Printf("Successfully changed to directory %s\n", absPath)

	return nil
}

func (k *Keeper) runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
