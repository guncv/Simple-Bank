package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/guncv/Simple-Bank/db/sqlc"
)

// Server serves HTTP requests for out banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// New Server creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	// add routes to router
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
