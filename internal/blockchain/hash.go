package blockchain

import (
	"fmt"
	"paychain/pkg/utils"
	"strconv"
)

func computeTxsString(txs []Transaction) string {
	if len(txs) == 0 {
		return ""
	}
	b := make([]byte, 0, len(txs)*32)
	for _, tx := range txs {
		b = append(b, []byte(tx.From)...)
		b = append(b, '|')
		b = append(b, []byte(tx.To)...)
		b = append(b, '|')
		b = append(b, []byte(fmt.Sprintf("%d|%d", tx.Amount, tx.Time))...)
		b = append(b, ';')
	}
	return string(b)
}

func ComputeBlockHash(b Block) string {
	return utils.HashStrings(
		strconv.Itoa(b.Index),
		strconv.FormatInt(b.Timestamp, 10),
		b.PrevHash,
		strconv.Itoa(b.Nonce),
		computeTxsString(b.Transactions),
	)
}
