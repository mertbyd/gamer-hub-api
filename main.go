package main

import (
	"api-steam/app"
	"api-steam/configs"
	"api-steam/repository"
	"api-steam/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Fiber app oluştur
	appRoute := fiber.New(fiber.Config{
		AppName: "API-Steam Game API",
	})

	// Yapay Zeka
	// DB bağlantısı configs.DB üzerinden zaten kurulmuş durumda
	dbClient := configs.GetCollection(configs.DB, "games")
	ProductRepositoryDB := repository.NewProductRepository(dbClient)
	productService := services.NewProductService(ProductRepositoryDB)
	productHandler := app.ProductHandler{Services: productService}
	// Yapay Zeka BİTİŞ

	// Game endpoint'i
	// CRUD ve Arama Endpointleri
	appRoute.Post("/api/game", productHandler.CreateProduct)                // Yeni bir oyun oluşturur
	appRoute.Get("/api/games", productHandler.GetAllProduct)                // Tüm oyunları listeler
	appRoute.Delete("/api/game/:id", productHandler.DeleteProduct)          // ID'ye göre oyun siler
	appRoute.Put("/api/game/:id", productHandler.UpdateProduct)             // ID'ye göre oyunu tamamen günceller
	appRoute.Patch("/api/game/:id", productHandler.PatchProduct)            // ID'ye göre oyunun belirli alanlarını günceller
	appRoute.Get("/api/game/:id", productHandler.GetByID)                   // ID'ye göre oyun getirir
	appRoute.Get("/api/games/sorted", productHandler.GetGamesSorted)        // Oyunları belirtilen alana göre sıralar (asc/desc)
	appRoute.Get("/api/games/exact", productHandler.GetGamesByExactName)    // Tam isim eşleşmesine göre oyun arar
	appRoute.Get("/api/games/search", productHandler.GetGamesByPartialName) // Kısmi isim eşleşmesine göre oyun arar
	// Sunucuyu başlat
	log.Println("Server 8080 portunda başlatılıyor...")
	log.Fatal(appRoute.Listen(":8080"))
}
