package main

import (
	"log"
	"os"

	"invoice_project/pkg/infrastructure"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"

	"invoice_project/internal/auth/delivery/http"
	authModel "invoice_project/internal/auth/domain"
	authRepo "invoice_project/internal/auth/repository"
	authUC "invoice_project/internal/auth/usecase"

	invHandler "invoice_project/internal/invoice/delivery/http"
	invModel "invoice_project/internal/invoice/domain"
	invRepo "invoice_project/internal/invoice/repository"
	invUC "invoice_project/internal/invoice/usecase"

	merchantHTTP "invoice_project/internal/merchant/delivery/http"
	merchModel "invoice_project/internal/merchant/domain"
	merchRepo "invoice_project/internal/merchant/repository"
	merchUC "invoice_project/internal/merchant/usecase"

	logModel "invoice_project/internal/log/domain"
)

func main() {
	// โหลด config จากไฟล์หรือ ENV
	configPath := "configs/config.yaml"
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		configPath = env
	}
	cfg, err := infrastructure.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	// สร้าง connection DB
	db, err := infrastructure.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	// Migrate tables: User, RefreshToken, Invoice
	infrastructure.Migrate(db,
		&authModel.User{},
		&authModel.RefreshToken{},
		&invModel.Invoice{},
		&merchModel.Merchant{},
		&merchModel.Store{},
		&merchModel.StoreAddress{},
		&merchModel.MerchantContact{},
		&merchModel.PersonMerchant{},
		&merchModel.CompanyMerchant{},
		&logModel.UserLog{},
	)

	// สร้าง Fiber app พร้อม ErrorHandler กลาง
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Logger middleware
	app.Use(middleware.Logger(db))
	// Global JWT middleware except for auth routes
	app.Use(middleware.JWTMiddlewareExcept(cfg.Auth.JWTSecret, "/auth"))

	// ตระเตรียม Auth module
	authRepository := authRepo.NewAuthRepository(db)
	authUsecase := authUC.NewAuthUsecase(
		authRepository,
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTExpiryAccessMin,
		cfg.Auth.JWTExpiryRefreshHours,
	)
	authHandler := http.NewAuthHandler(authUsecase, cfg.Auth.JWTSecret)
	authHandler.RegisterRoutes(app)

	// ตระเตรียม Invoice module
	invoiceRepository := invRepo.NewInvoiceRepository(db)
	invoiceUsecase := invUC.NewInvoiceUsecase(invoiceRepository)
	invoiceHandler := invHandler.NewInvoiceHandler(invoiceUsecase)
	invoiceHandler.RegisterRoutes(app)

	// Merchant module
	merchRepository := merchRepo.NewMerchantRepository(db)
	merchUsecase := merchUC.NewMerchantUsecase(merchRepository)
	merchantHandler := merchantHTTP.NewMerchantHandler(merchUsecase)
	merchantHandler.RegisterRoutes(app)

	// สตาร์ทเซิร์ฟเวอร์
	log.Printf("Server is running on port %s\n", cfg.Server.Port)
	if err := app.Listen(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
