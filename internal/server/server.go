package server

import "github.com/gin-gonic/gin"

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

// TODO : s.handleConversion

func (s *Server) handleConversion(c *gin.Context) {

}
