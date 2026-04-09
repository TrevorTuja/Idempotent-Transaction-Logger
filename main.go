package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	sender := flag.String("sender", "", "Sender name")
	receiver := flag.String("receiver", "", "Receiver name")
	amount := flag.Float64("amount", 0, "Transfer amount")
	key := flag.String("key", "", "Idempotency key")
	status := flag.String("status", "", "Query by idempotency key")

	flag.Parse()

	if *status != "" {
		handleStatusQuery(*status)
		return
	}

	if *sender == "" || *receiver == "" || *amount <= 0 || *key == "" {
		fmt.Println("Usage: fintech-go -sender NAME -receiver NAME -amount NUMBER -key KEY")
		fmt.Println("       fintech-go -status KEY")
		os.Exit(1)
	}

	handleTransaction(*sender, *receiver, *amount, *key)
}

func handleTransaction(sender, receiver string, amount float64, key string) {
	ledger, err := LoadLedger()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	existing := FindByIdemptotencyKey(ledger, key)

	if existing != nil {
		msg := fmt.Sprintf(
			"[DUPLICATE DETECTED] Returning original result:\n"+
				"  Transaction ID : %s\n"+
				"  Idempotency Key: %s\n"+
				"  Sender         : %s\n"+
				"  Receiver       : %s\n"+
				"  Amount         : $%.2f\n"+
				"  Status         : %s\n"+
				"  Timestamp      : %s",
			existing.TransactionID,
			existing.IdempotencyKey,
			existing.Sender,
			existing.Receiver,
			existing.Amount,
			existing.Status,
			existing.Timestamp,
		)
		LogTransaction(msg)
		return
	}

	newTx := Transaction{
		TransactionID:  uuid.New().String(),
		IdempotencyKey: key,
		Sender:         sender,
		Receiver:       receiver,
		Amount:         amount,
		Status:         "completed",
		Timestamp:      time.Now().Format(time.RFC3339),
	}

	ledger.Transactions = append(ledger.Transactions, newTx)

	err = SaveLedger(ledger)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	msg := fmt.Sprintf(
		"[NEW TRANSACTION] Successfully processed:\n"+
			"  Transaction ID : %s\n"+
			"  Idempotency Key: %s\n"+
			"  Sender         : %s\n"+
			"  Receiver       : %s\n"+
			"  Amount         : $%.2f\n"+
			"  Status         : %s\n"+
			"  Timestamp      : %s",
		newTx.TransactionID,
		newTx.IdempotencyKey,
		newTx.Sender,
		newTx.Receiver,
		newTx.Amount,
		newTx.Status,
		newTx.Timestamp,
	)
	LogTransaction(msg)
}

func handleStatusQuery(key string) {
	ledger, err := LoadLedger()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	tx := FindByIdemptotencyKey(ledger, key)
	if tx == nil {
		fmt.Printf("No transaction found for key: %s\n", key)
		os.Exit(1)
	}

	fmt.Printf(
		"[TRANSACTION STATUS]\n"+
			" Transaction ID : %s\n"+
			" Idempotency Key: %s\n"+
			" Sender         : %s\n"+
			" Receiver       : %s\n"+
			" Amount         : $%.2f\n"+
			" Status         : %s\n"+
			" Timestamp      : %s",
		tx.TransactionID,
		tx.IdempotencyKey,
		tx.Sender,
		tx.Receiver,
		tx.Amount,
		tx.Status,
		tx.Timestamp,
	)
}
