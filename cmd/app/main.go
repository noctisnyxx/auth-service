package main

//	@title						Auth Core REST API Docs
//	@version					2.0
//	@description				In progress
//	@BasePath					/api
//	@contact.name				Dinar
//	@contact.email				dinar.hadiyanto@outlook.com
//	@securityDefinitions.basic	BasicAuth
import (
	"errors"
	"fmt"
	"net/http"

	"github.com/noctisnyxx/auth-service/configs"
	_ "github.com/noctisnyxx/auth-service/docs"
	"github.com/noctisnyxx/auth-service/internal/repository"
	"github.com/noctisnyxx/auth-service/internal/transport"
	"github.com/noctisnyxx/auth-service/internal/usecase"
)

func main() {
	logger := configs.SetupLogger()
	sqlDB, err := configs.SetupDatabase()
	if err != nil {
		logger.Error("failed to setup database", "error", err)
		return
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	userRepo := repository.NewUserRepository(sqlDB, logger)
	userUsecase := usecase.NewUserUsecase(logger, userRepo)
	userHandler := transport.NewUserHandler(logger, userUsecase)

	clientRepo := repository.NewClientRepository(logger, sqlDB)
	clientUsecase := usecase.NewClientUsecase(logger, clientRepo)
	clientHandler := transport.NewClientHandler(logger, clientUsecase)

	router := transport.NewRouter(userHandler, clientHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.AppPort),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to listen and serve", "error", err)
	}
}
