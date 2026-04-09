package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const ledgerFile = "ledger.json"

func LoadLedger() (Ledger, error) {
	var ledger Ledger

	data, err := os.ReadFile(ledgerFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ledger, nil
		}
		return ledger, fmt.Errorf("failed to read ledger: %w", err)
	}

	err = json.Unmarshal(data, &ledger)
	if err != nil {
		return ledger, fmt.Errorf("failed to parse ledger: %w", err)
	}

	return ledger, nil
}

func SaveLedger(ledger Ledger) error {
	data, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ledger: %w", err)
	}

	err = os.WriteFile(ledgerFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write ledger: %w", err)
	}

	return nil
}

func FindByIdemptotencyKey(ledger Ledger, key string) *Transaction {
	for _, tx := range ledger.Transactions {
		if tx.IdempotencyKey == key {
			return &tx
		}
	}
	return nil
}
