package main

import (
	"log"
	"os"

	"invoice_project/pkg/infrastructure"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

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

	customerHTTP "invoice_project/internal/customer/delivery/http"
	customerModel "invoice_project/internal/customer/domain"
	customerRepo "invoice_project/internal/customer/repository"
	customerUC "invoice_project/internal/customer/usecase"

	productHTTP "invoice_project/internal/product/delivery/http"
	productModel "invoice_project/internal/product/domain"
	productRepo "invoice_project/internal/product/repository"
	productUC "invoice_project/internal/product/usecase"

	locationHTTP "invoice_project/internal/location/delivery/http"
	locationModel "invoice_project/internal/location/domain"
	locationRepo "invoice_project/internal/location/repository"
	locationUC "invoice_project/internal/location/usecase"

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
		&customerModel.Customer{},
		&customerModel.CompanyCustomer{},
		&customerModel.PersonCustomer{},
		&customerModel.CustomerAddress{},
		&customerModel.CustomerContact{},
		&productModel.Product{},
		&productModel.ProductImage{},
		&locationModel.Province{},
		&locationModel.District{},
		&locationModel.SubDistrict{},
		&logModel.UserLog{},
	)

	// สร้าง Fiber app พร้อม ErrorHandler กลาง
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// ✅ เปิดใช้งาน CORS สำหรับทุก origin (เพื่อเริ่มต้น)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Logger middleware
	app.Use(middleware.Logger(db))
	// Global JWT middleware except for auth routes
	app.Use(middleware.JWTMiddlewareExcept(cfg.Auth.JWTSecret, "/auth"))

	// Merchant module
	merchRepository := merchRepo.NewMerchantRepository(db)
	merchUsecase := merchUC.NewMerchantUsecase(merchRepository)
	merchantHandler := merchantHTTP.NewMerchantHandler(merchUsecase)
	merchantHandler.RegisterRoutes(app)

	// ตระเตรียม Auth module
	authRepository := authRepo.NewAuthRepository(db)
	authUsecase := authUC.NewAuthUsecase(
		authRepository,
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTExpiryAccessMin,
		cfg.Auth.JWTExpiryRefreshHours,
	)
	authHandler := http.NewAuthHandler(authUsecase, merchUsecase, cfg.Auth.JWTSecret)
	authHandler.RegisterRoutes(app)

	// ตระเตรียม Invoice module
	invoiceRepository := invRepo.NewInvoiceRepository(db)
	invoiceUsecase := invUC.NewInvoiceUsecase(invoiceRepository)
	invoiceHandler := invHandler.NewInvoiceHandler(invoiceUsecase)
	invoiceHandler.RegisterRoutes(app)

	// Customer module
	customerRepository := customerRepo.NewCustomerRepository(db)
	customerUseCase := customerUC.NewCustomerUseCase(customerRepository)
	customerHandler := customerHTTP.NewCustomerHandler(customerUseCase)
	customerHandler.RegisterRoutes(app)

	// Product module
	productRepository := productRepo.NewProductRepository(db)
	productUsecase := productUC.NewProductUseCase(productRepository)
	productHandler := productHTTP.NewProductHandler(productUsecase)
	productHandler.RegisterRoutes(app)

	// Location module
	locationRepository := locationRepo.NewLocationRepository(db)
	locationUsecase := locationUC.NewLocationUseCase(locationRepository)
	locationHandler := locationHTTP.NewLocationHandler(locationUsecase)
	locationHandler.RegisterRoutes(app)

	// สตาร์ทเซิร์ฟเวอร์
	log.Printf("Server is running on port %s\n", cfg.Server.Port)
	if err := app.Listen(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
