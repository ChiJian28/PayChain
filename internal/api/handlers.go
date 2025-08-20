package api

import (
	"net/http"
	"paychain/internal/account"
	"paychain/internal/blockchain"
	"paychain/internal/kafka"
	txpool "paychain/internal/pool"
	"paychain/pkg/utils"

	"github.com/gin-gonic/gin"
)

type transferReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func RegisterRoutes(r *gin.Engine, prod *kafka.Producer, acct *account.Store, chain *blockchain.Chain, pool *txpool.Pool) {
	r.POST("/transfer", func(c *gin.Context) {
		var req transferReq
		if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 || req.To == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		tx := blockchain.Transaction{From: req.From, To: req.To, Amount: req.Amount, Time: utils.NowUnix()}
		if err := prod.PublishTransaction(tx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "queued"})
	})

	r.GET("/balance/:user", func(c *gin.Context) {
		user := c.Param("user")
		c.JSON(http.StatusOK, gin.H{"user": user, "balance": acct.GetBalance(user)})
	})

	r.GET("/blockchain", func(c *gin.Context) {
		c.JSON(http.StatusOK, chain.All())
	})

	r.GET("/pending", func(c *gin.Context) {
		c.JSON(http.StatusOK, pool.List())
	})
}
