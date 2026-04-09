package main

type Transaction struct {
	TransactionID  string  `json:"transaction_id"`
	IdempotencyKey string  `json:"idempotency_key"`
	Sender         string  `json:"sender"`
	Receiver       string  `json:"receiver"`
	Amount         float64 `json:"amount"`
	Timestamp      string  `json:"timestamp"`
	Status         string  `json:"status"`
}

type Ledger struct {
	Transactions []Transaction `json:"transactions"`
}
