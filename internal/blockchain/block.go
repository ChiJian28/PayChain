package blockchain

type Transaction struct {
	From   string
	To     string
	Amount int
	Time   int64
}

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}
