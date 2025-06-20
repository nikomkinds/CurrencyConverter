package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nikomkinds/CurrencyConverter/internal/converter"
	"net/http"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()
	s := &Server{router: r}
	s.configureRoutes()
	return s
}

func (s *Server) configureRoutes() {
	s.router.GET("/convert", s.handleConversion)
}

func (s *Server) Run() {
	s.router.Run(":8080")
}

func (s *Server) handleConversion(c *gin.Context) {

	from := c.Query("from")
	to := c.Query("to")
	amount := c.Query("amount")
	dateReq := c.Query("date_req")

	if from == "" || to == "" || amount == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all three params 'from', 'to' and 'amount' are required"})
		return
	}

	result, err := converter.Convert(from, to, amount, dateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
