package app

import (
	"api-steam/models"
	"api-steam/services"
	"net/http"

	"github.com/gofiber/fiber/v2" //web framework'ünde bir HTTP isteğini ve cevabını temsil eden context (bağlam) nesnesidir.
	"go.mongodb.org/mongo-driver/bson/primitive"
	/*İstek verilerini okuma (parametreler, gövde vb.)
	  Cevap oluşturma (JSON, metin, dosya vb.)
	  HTTP durumunu ve başlıklarını yönetme
	  İstek/cevap döngüsünün akışını kontrol etme*/)

type ProductHandler struct {
	Services services.ProductService //servis katmanına erişim sağlamamız için
}

// ürün yaratırken kulanılan
func (h ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var game models.Game
	if err := c.BodyParser(&game); err != nil { //fadesi, Fiber web framework'ünde HTTP isteğinin gövdesini (body) Go struct'ına dönüştürmek için kullanılır.
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	} //go nesnesini json veri tipine dönüştürür
	result, err := h.Services.ProductInsert(game)
	if err != nil || !result.Status {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(http.StatusCreated).JSON(result)
}

// ürün listesini getirirken
func (h ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	result, err := h.Services.ProductGetAll() // Burada h.GetAllProduct() yerine h.Services.GetAll() kullanılmalı
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(result)
	//statuesOk değeri go da 200 değerine sahiptir yani işlemin tamamlandığını belirtir .json ile istediğimiz veriyi json tipinde görüürüz
	//eğer log lama yaparken hata almazsak json tipinde verimizi
}

// silme işlem yaparken
func (h ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	query := c.Params("id")                    //bizim url e verdğimiz yani appRoute.Delete("/api/game/:id", productHandler.DeleteProduct) id yi çeker ve bunu atar
	cnv, _ := primitive.ObjectIDFromHex(query) // ID'yi MongoDB ObjectID'ye dönüştür
	result, err := h.Services.ProductDelete(cnv)
	if err != nil || result == false {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"State": false})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"State": true})
}

// İd ye göre güncelme yaparken
func (h ProductHandler) UpdateProduct(c *fiber.Ctx) error {

	id := c.Params("id")                           //bizim url e verdğimiz yani appRoute.Delete("/api/game/:id", productHandler.DeleteProduct) id yi çeker ve bunu atar
	objectID, err := primitive.ObjectIDFromHex(id) /// ID'yi MongoDB ObjectID'ye dönüştür

	var updatedGame models.Game
	if err := c.BodyParser(&updatedGame); err != nil { //fadesi, Fiber web framework'ünde HTTP isteğinin gövdesini (body) Go struct'ına dönüştürmek için kullanılır.
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"state": false,
		})
	}
	result, err := h.Services.ProductUptade(objectID, updatedGame)
	if err != nil || result == false {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"state": false,
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"state": true,
	})
}

// PatchProduct, PATCH isteği ile bir ürünün belirli alanlarını günceller
func (h ProductHandler) PatchProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id) // ID'yi MongoDB ObjectID'ye dönüştür
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"state": false, "error": "Geçersiz ID formatı"})
	}

	// İstek gövdesini map olarak ayrıştır
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"state": false, "error": "İstek gövdesi ayrıştırılamadı: " + err.Error()})
	}

	// Servis katmanını çağır
	result, err := h.Services.ProductPatch(objectID, updates)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"state": false, "error": "Güncelleme sırasında hata oluştu: " + err.Error()})
	}

	if !result {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"state": false, "message": "Güncellenecek ürün bulunamadı"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"state": true,
	})
}

// ID ye  göre veri getirme
func (h ProductHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")                           //url lemdki id kısmını alıcak
	objectID, err := primitive.ObjectIDFromHex(id) //id değerini mongodb id değerine göre dönüştürür
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz ID formatı"})
	}
	result, err := h.Services.ProductGetByID(objectID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Oyun bulunamadı"})
	}
	return c.Status(http.StatusOK).JSON(result)
}

// Fiyata göre sıralama
func (h ProductHandler) GetGamesSorted(c *fiber.Ctx) error {
	query := c.Query("field", "price.amount") // field: Hangi alana göre sıralama yapılacağını belirtir (örn. "price.amount")
	orderType := c.Query("order", "asc")      //url den gelen talep
	var order int
	if orderType == "asc" {
		order = 1
	} else {
		order = -1
	}
	games, err := h.Services.ProductGetSorted(query, order)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(games)
}
func (h ProductHandler) GetGamesByExactName(c *fiber.Ctx) error {
	name := c.Query("name") // url den name parametresini alır eğer name parametresi yoksa empty döner buna göre logladım
	if name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "İsim parametresi gereklidir"})
	}

	result, err := h.Services.ProductGetByExactName(name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Oyunlar getirilirken bir hata oluştu"})
	}

	return c.Status(http.StatusOK).JSON(result)
}

// Kısmi isme göre sıralma
func (h ProductHandler) GetGamesByPartialName(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "İsim parametresi gereklidir"})
	}

	result, err := h.Services.ProductGetByPartialName(name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Oyunlar getirilirken bir hata oluştu"})
	}

	return c.Status(http.StatusOK).JSON(result)
}
