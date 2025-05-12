package main

import (
	"api-steam/app"
	"api-steam/configs"
	"api-steam/repository"
	"api-steam/services"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance oluştur
	e := echo.New()

	// DB bağlantısı configs.DB üzerinden zaten kurulmuş durumda
	dbClient := configs.GetCollection(configs.DB, "games")            //tabloya bağlanmak için
	productRepositoryDB := repository.NewProductRepository(dbClient)  //Repistory katmanına bağlantı nesnesini veririz
	productService := services.NewProductService(productRepositoryDB) // servis katmanında repistory katmanındakifonksiyonlara erişmek için
	productHandler := app.ProductHandler{Services: productService}    //handlerda kulancağımız servis elamanları için handlera servis den bir nesne veiriz

	//endpointi
	e.POST("/api/game", productHandler.CreateProduct)                // Yeni bir oyun oluşturur
	e.GET("/api/games", productHandler.GetAllProduct)                // Tüm oyunları listeler
	e.DELETE("/api/game/:id", productHandler.DeleteProduct)          // ID'ye göre oyun siler
	e.PUT("/api/game/:id", productHandler.UpdateProduct)             // ID'ye göre oyunu tamamen günceller
	e.PATCH("/api/game/:id", productHandler.PatchProduct)            // ID'ye göre oyunun belirli alanlarını günceller
	e.GET("/api/game/:id", productHandler.GetByID)                   // ID'ye göre oyun getirir
	e.GET("/api/games/sorted", productHandler.GetGamesSorted)        // Oyunları belirtilen alana göre sıralar (asc/desc)
	e.GET("/api/games/exact", productHandler.GetGamesByExactName)    // Tam isim eşleşmesine göre oyun arar
	e.GET("/api/games/search", productHandler.GetGamesByPartialName) // Kısmi isim eşleşmesine göre oyun arar

	// Sunucuyu başlat
	log.Println("Server 8080 portunda başlatılıyor...")
	log.Fatal(e.Start(":8080"))
}
