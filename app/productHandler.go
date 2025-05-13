package app

import (
	"api-steam/models"
	"api-steam/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	Services services.ProductService
}

// CreateProduct - HTTP POST isteği ile yeni bir oyun oluşturur
func (h ProductHandler) CreateProduct(c echo.Context) error {
	var game models.Game
	if err := c.Bind(&game); err != nil { //c.Bind ile Http nin boudy ksımındaki json esneisi go nesnesine dönüştürürüz
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz istek formatı: " + err.Error()}) //c.JSON =Htttp yantını json formatına dönüştürür htt.StatusBadRequest ile 400 hata kodnunu döneriz map[string] ile inerface{} herhanig bşr nesne demek eror etiketi ile err.error kodunu eşleriz
	}
	result, err := h.Services.ProductInsert(game)
	if err != nil || !result.Status {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Oyun eklenirken hata oluştu: " + err.Error()})
	}
	return c.JSON(http.StatusCreated, result) //dto dakii state -true ve 200 başarı kodunu döneriz
}

// CreateManyProducts - HTTP POST isteği ile birden fazla oyunu toplu olarak ekler
func (h ProductHandler) CreateManyProducts(c echo.Context) error {
	var games []models.Game
	if err := c.Bind(&games); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz istek formatı: " + err.Error()})
	}
	if len(games) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "En az bir oyun göndermelisiniz"})
	}
	result, err := h.Services.ProductInsertMany(games)
	if err != nil || !result.Status {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Toplu oyun eklenirken hata oluştu: " + err.Error()})
	}
	return c.JSON(http.StatusCreated, result) //dto dakii state -true ve 200 başarı kodunu döneriz
}

// GetAllProduct - HTTP GET isteği ile tüm oyunları listeleyerek döndürür
func (h ProductHandler) GetAllProduct(c echo.Context) error {
	result, err := h.Services.ProductGetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Oyunlar listelenirken hata oluştu: " + err.Error()})
	}
	return c.JSON(http.StatusOK, result) //json tipinde []models.Game dizisini ve StatusOK Http kodunu json tipinde döneriz
}

// DeleteProduct - HTTP DELETE isteği ile belirtilen ID'ye sahip oyunu siler
func (h ProductHandler) DeleteProduct(c echo.Context) error {
	query := c.Param("id")
	cnv, err := primitive.ObjectIDFromHex(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz ID formatı: ID bir MongoDB ObjectID olmalıdır"}) //400 hata kodunu Json tipinde öner eror etiketiyle eror mesajını eşlerüiz
	}
	result, err := h.Services.ProductDelete(cnv)
	if err != nil || result == false {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"state": false, "error": "Oyun silinirken hata oluştu veya oyun bulunamadı"}) //işlem gerçekleşemediği için json tipinde StatusBadRequest hata kodnu ve state i false olarak döneriz
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"state": true, "message": "Oyun başarıyla silindi"}) //200 işlem başarılı kodunu döneriz  state :true ile true mesajı döneriz
}

// UpdateProduct - HTTP PUT isteği ile belirtilen ID'ye sahip oyunu tamamen günceller
func (h ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")                            //Url deki id parametrisini alırız
	objectID, err := primitive.ObjectIDFromHex(id) //aldığımız string tipindeki ıd değerini monodb id tipine dönüştürür
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz ID formatı: ID bir MongoDB ObjectID olmalıdır"}) ///400 hata kodunu Json tipinde öner eror etiketiyle eror mesajını eşlerüiz
	}
	var updatedGame models.Game
	if err := c.Bind(&updatedGame); err != nil { //c.Bind http den gelen boudy yi gyani game nesnesinin json tipini &updategame in referansına atayabilirzse  err bil döner dmnemezse err hata mesajı döner
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"state": false, "error": "Geçersiz istek formatı: " + err.Error()}) //err  hata kodunu Json tipinde döner işlem gerçekleşmediği için statei false yaparız
	}
	result, err := h.Services.ProductUptade(objectID, updatedGame)
	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"state": false, "error": "Oyun güncellenirken hata oluştu veya oyun bulunamadı"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"state": true, "message": "Oyun başarıyla güncellendi"})
}

// PatchProduct - HTTP PATCH isteği ile belirtilen ID'ye sahip oyunun belirli alanlarını günceller
func (h ProductHandler) PatchProduct(c echo.Context) error {
	id := c.Param("id")                            //Url mizdeki strin id değerini alır
	objectID, err := primitive.ObjectIDFromHex(id) //ObjectIDFromHex string alınan id değerini mongodb id tipine dönüştürür
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"state": false, "error": "Geçersiz ID formatı: ID bir MongoDB ObjectID olmalıdır"})
	}
	// İstek gövdesini map olarak ayrıştır
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"state": false, "error": "İstek gövdesi ayrıştırılamadı: " + err.Error()})
	}
	// Servis katmanını çağır
	result, err := h.Services.ProductPatch(objectID, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"state": false, "error": "Güncelleme sırasında hata oluştu: " + err.Error()}) //400 hata kodunu Json tipinde öner eror etiketiyle eror mesajını eşlerüiz
	}
	if !result {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"state": false, "message": "Güncellenecek ürün bulunamadı"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"state": true, "message": "Oyun alanları başarıyla güncellendi"})
}

// GetByID - HTTP GET isteği ile belirtilen ID'ye sahip oyunu getirir
func (h ProductHandler) GetByID(c echo.Context) error {
	id := c.Param("id")                            //url deki id veri tipini alır
	objectID, err := primitive.ObjectIDFromHex(id) //alınan id strngini mongodbid tiine dönüştürür
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz ID formatı: ID bir MongoDB ObjectID olmalıdır"}) //400 hata kodunu Json tipinde öner eror mesajı eror mesajını eşlerüiz
	}
	result, err := h.Services.ProductGetByID(objectID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Belirtilen ID'ye sahip oyun bulunamadı"}) //400 hata kodunu Json tipinde öner eror etiketiyle eror mesajını eşlerüiz
	}
	return c.JSON(http.StatusOK, result) //
}

// GetGamesSorted - HTTP GET isteği ile oyunları belirtilen alana ve sıralama yönüne göre sıralar
func (h ProductHandler) GetGamesSorted(c echo.Context) error {
	query := c.QueryParam("field") //QueryParam() sorgu parametresine verilen değeri almak için kulanılr  field a verilen değeri alır bunu artandan azalana yada azalandan artana sıralamak için kulanırız
	if query == "" {
		query = "price.amount" // fiayata göre sıralayacaımız için query price.amount olarak default olarak ayarlanır
	}
	orderType := c.QueryParam("order")
	if orderType == "" {
		orderType = "asc" // Varsayılan değer
	}
	var order int
	if orderType == "asc" {
		order = 1 //artanda-azalana 1
	} else {
		order = -1 //azalandan artana -1
	}
	games, err := h.Services.ProductGetSorted(query, order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Oyunlar sıralanırken hata oluştu: " + err.Error()})
	}
	return c.JSON(http.StatusOK, games)
}

// GetGamesByExactName - HTTP GET isteği ile tam isim eşleşmesine göre oyunları arar
func (h ProductHandler) GetGamesByExactName(c echo.Context) error {
	name := c.QueryParam("name") //url deki name etiketine  verilen değeri çekme için kulanılır
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "İsim parametresi gereklidir: ?name=<oyun adı> formatında gönderilmelidir"})
	}
	result, err := h.Services.ProductGetByExactName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Oyunlar isimle aranırken hata oluştu: " + err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetGamesByPartialName - HTTP GET isteği ile kısmi isim eşleşmesine göre oyunları arar
func (h ProductHandler) GetGamesByPartialName(c echo.Context) error {
	name := c.QueryParam("name") //url deki name etiketine  verilen değeri çekme için kulanılır
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "İsim parametresi gereklidir: ?name=<oyun adının bir parçası> formatında gönderilmelidir"})
	}
	result, err := h.Services.ProductGetByPartialName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Oyunlar kısmi isimle aranırken hata oluştu: " + err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h ProductHandler) GetGamesByPriceRange(c echo.Context) error {
	minPriceStr := c.QueryParam("min")
	maxPriceStr := c.QueryParam("max")

	// Varsayılan değerler
	minPrice := 0.0
	maxPrice := 1000.0 // Yüksek bir default değer

	//Yapay zeka
	// Min fiyat parametresi varsa parse et
	if minPriceStr != "" {
		var err error
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Geçersiz minimum fiyat değeri: " + err.Error(),
			})
		}
	}

	// Max fiyat parametresi varsa parse et
	if maxPriceStr != "" {
		var err error
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Geçersiz maksimum fiyat değeri: " + err.Error(),
			})
		}
	}

	// Min değer max değerden büyük olamaz
	if minPrice > maxPrice {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Minimum fiyat, maksimum fiyattan büyük olamaz",
		})
	}
	//yapay zeka

	result, err := h.Services.ProductGetByPriceRange(minPrice, maxPrice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Fiyat aralığına göre oyunlar getirilirken hata oluştu: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

//objectID, err := primitive.ObjectIDFromHex(id)//aldığımız string tipindeki ıd değerini monodb id tipine dönüştürür
//err := c.Bind(&updatedGame); err != nil {//c.Bind http den gelen boudy yi gyani game nesnesinin json tipini &updategame in referansına atayabilirzse  err bil döner dmnemezse err hata mesajı döner
//query := c.QueryParam("field")//QueryParam() sorgu parametresine verilen değeri almak için kulanılr  field a verilen değeri alır bunu artandan azalana yada azalandan artana sıralamak için kulanırız
