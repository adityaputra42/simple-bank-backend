package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serve http request for our banking service
type Server struct {
	store db.Store
	route *gin.Engine
}

// NerServer create a new HTTP Server and Setup Routing
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.creatAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.route = router
	return server

}

// Start Run the HTTP server on specific address
func (server *Server) Start(address string) error {
	return server.route.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
