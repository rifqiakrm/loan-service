package main

import (
	"loan-service/api"
	"loan-service/core/loan"
	"loan-service/email"
)

func main() {
	// Setup repository, email mock, and service
	repo := loan.NewInMemoryLoanRepository()
	mailer := email.NewMockEmailSender()
	service := loan.NewLoanService(repo, mailer)

	// Setup HTTP handler and routes
	handler := api.NewHandler(service)
	router := api.SetupRouter(handler)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
