package main

import (
	"fmt"
	"os"
)

const logFile = "transactions.log"

func LogTransaction(message string) error {
	fmt.Println(message)

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, message)
	if err != nil {
		return fmt.Errorf("failed to write to log: %w", err)
	}

	return nil
}
