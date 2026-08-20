package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/adapters"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/adapters/handlers"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/adapters/repositories"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/docs"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/infrastructure"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/middleware"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

const (
	corsMaxAge = 12 * time.Hour
)

func Run(configPath string) error {
	// Setup infrastructure layer
	infra, err := infrastructure.SetupInfrastructure(configPath)
	if err != nil {
		return fmt.Errorf("failed to setup infrastructure: %w", err)
	}
	defer func() {
		if cleanupErr := infra.Cleanup(); cleanupErr != nil {
			logger.Error("Error cleaning up infrastructure", "error", cleanupErr)
		}
	}()

	// Extract from infrastructure
	conf := infra.Config
	dbClient := infra.DBClient

	// Secondary Adapters
	txRepo := repositories.NewTransactionRepository(dbClient)
	tokenRepo := repositories.NewTokenRepository(dbClient)
	participantRepo := repositories.NewParticipantRepository(dbClient)
	headerProofRepo := repositories.NewHeaderProofRepository(dbClient)
	pnRepo := repositories.NewPrivateNetworkRepository(dbClient)
	balanceRepo := repositories.NewBalanceRepository(dbClient)
	metadataService := adapters.NewTokenMetadataService()
	log := logger.NewLogger()

	// Core Services
	transactionService := core.NewTransactionService(txRepo, metadataService, log)
	participantService := core.NewParticipantService(participantRepo, log)
	tokenServiceCore := core.NewTokenService(tokenRepo, log)
	headerProofService := core.NewHeaderProofService(headerProofRepo, log)
	authService := core.NewAuthService(pnRepo, []byte(conf.JWTSecret), log)
	balanceService := core.NewBalanceService(balanceRepo, log)

	// Primary Adapters
	transactionHandler := handlers.NewTransactionHandler(transactionService, log)
	participantHandler := handlers.NewParticipantHandler(participantService, log)
	tokenHandler := handlers.NewTokenHandler(tokenServiceCore, log)
	headerProofHandler := handlers.NewHeaderProofHandler(headerProofService, log)
	authHandler := handlers.NewAuthHandler(authService, log)
	balanceHandler := handlers.NewBalanceHandler(balanceService, log)

	// Create a new Gin router
	router := gin.Default()

	// Swagger Config
	docs.SwaggerInfo.Title = "Rayls Audit API"

	allowedOrigins := strings.Split(conf.CorsUrls, ";")

	// CORS Configs
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           corsMaxAge,
	}))

	router.Use(middleware.ValidateQueryEncoding())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	auditGroup := router.Group("/audit")
	{
		transactionsRoutes := auditGroup.Group("/transactions")
		{
			transactionsRoutes.GET("/enygma/batch/:batchId", transactionHandler.GetTransactionsByEnygmaBatchId)
			transactionsRoutes.GET("/batch/:batchId", transactionHandler.GetTransactionsByRegularBatchId)
			transactionsRoutes.GET(
				"",
				middleware.ValidateQueryParams(dto.MergedTransactionsFilters{}),
				transactionHandler.GetTransactions,
			)
			transactionsRoutes.GET("/:messageId", transactionHandler.GetTransactionByMessageId)
			transactionsRoutes.GET("/dvp/:transactionId", transactionHandler.GetTransactionByTransactionId)
			transactionsRoutes.GET("/dvp/swap/:sharedId", transactionHandler.GetTransactionsBySharedId)
		}

		participantsRoutes := auditGroup.Group("/participants")
		{
			participantsRoutes.GET("/:chainId", participantHandler.GetParticipantByChainId)
			participantsRoutes.GET(
				"",
				middleware.ValidateQueryParams(dto.ParticipantListFilters{}),
				participantHandler.GetParticipantList,
			)
		}

		tokensRoutes := auditGroup.Group("/tokens")
		{
			tokensRoutes.GET("/:resourceId", tokenHandler.GetTokenByResourceId)
			tokensRoutes.GET("", middleware.ValidateQueryParams(dto.TokenListFilters{}), tokenHandler.GetTokenList)
		}

		headerProofsRoutes := auditGroup.Group("/header-proofs")
		{
			headerProofsRoutes.GET(
				"",
				middleware.ValidateQueryParams(dto.HeaderProofFilters{}),
				headerProofHandler.GetHeaderProofsList,
			)
		}
	}

	router.GET("/flagged", transactionHandler.GetFlaggedTransactions)

	router.POST("/private-network/signup", authHandler.SignUp)

	router.POST("/private-network/login", authHandler.Login)

	authMiddleware := middleware.RequireAuth(authService)
	// Supports wildcard: /resources/:chainid/ (all) or /resources/:chainid/:resourceid (specific)
	router.GET("/resources/:chainid/*resourceid", authMiddleware, balanceHandler.GetBalancesInChain)

	router.GET("/resource_info_all_chains/:resourceid", authMiddleware, balanceHandler.GetBalanceAcrossAllChains)

	router.POST("/resource_info_list_chains", authMiddleware, balanceHandler.GetBalanceAcrossSpecificChains)

	router.GET("/participant_info/:chainId", authMiddleware, participantHandler.GetParticipantByChainId)

	router.GET("/token_status/:resourceid", authMiddleware, tokenHandler.GetTokenRegistryStatus)

	// Run the application on port 8080
	if err := router.Run(":8080"); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
