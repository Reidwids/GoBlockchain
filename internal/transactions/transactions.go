package transactions

import "github.com/google/uuid"

type Transaction struct {
	ID        uuid.UUID
	Sender    string
	Recipient string
	Amount    float32
}
