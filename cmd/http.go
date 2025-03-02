package cmd

import (
	"ewallet-transaction/external"
	"ewallet-transaction/helpers"
	healthcheckHandler "ewallet-transaction/internal/handler/healthcheck"
	transactionHandler "ewallet-transaction/internal/handler/transaction"
	healthcheckRepo "ewallet-transaction/internal/repository/healthcheck"
	transactionRepo "ewallet-transaction/internal/repository/transaction"
	healthcheckSvc "ewallet-transaction/internal/services/healthcheck"
	transactionSvc "ewallet-transaction/internal/services/transaction"
	"ewallet-transaction/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeHttp() {
	r := gin.Default()

	healthcheckRepo := healthcheckRepo.NewRepository()

	healthcheckSvc := healthcheckSvc.NewService(healthcheckRepo)

	external := &external.External{}

	middleware := &middleware.ExternalDependency{
		External: external,
	}

	transactionRepo := transactionRepo.NewRepository(helpers.DB)
	transactionSvc := transactionSvc.NewService(transactionRepo, external)
	transactionHandler := transactionHandler.NewHandler(r, transactionSvc, external, middleware)
	transactionHandler.RegisterRoute()

	healthcheckHandler := healthcheckHandler.NewHandler(r, healthcheckSvc)
	healthcheckHandler.RegisterRoute()

	err := r.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}
}
