package app

import (
	"api-steam/models"
	"api-steam/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	Services services.ProductService
}

// CreateProduct - ürün oluşturma handler'ı
func (h ProductHandler) CreateProduct(c echo.Context) error {
	var game models.Game
	if err := c.Bind(&game); err != nil { //c.Bind ile Http nin boudy ksımındaki json esneisi go nesnesine dönüştürürüz
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()}) //c.JSON =Htttp yantını json formatına dönüştürür htt.StatusBadRequest ile 400 hata kodnunu döneriz map[string] ile inerface{} herhanig bşr nesne demek eror etiketi ile err.error kodunu eşleriz
	}

	result, err := h.Services.ProductInsert(game)
	if err != nil || !result.Status {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, result) //dto dakii state -true ve 200 başarı kodunu döneriz
}

// GetAllProduct - tüm ürünleri getirme
func (h ProductHandler) GetAllProduct(c echo.Context) error {
	result, err := h.Services.ProductGetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, result) //json tipinde []models.Game dizisini ve StatusOK Http kodunu json tipinde döneriz
}

// DeleteProduct - ürün silme
func (h ProductHandler) DeleteProduct(c echo.Context) error {
	query := c.Param("id")
	cnv, err := primitive.ObjectIDFromHex(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz ID formatı"}) //400 hata kodunu Json tipinde öner eror etiketiyle eror mesajını eşlerüiz
	}

	result, err := h.Services.ProductDelete(cnv)
	if err != nil || result == false {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"state": false}) //işlem gerçekleşemediği için json tipinde StatusBadRequest hata kodnu ve state i false olarak döneriz
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"state": true}) //200 işlem başarılı kodunu döneriz  state :true ile true mesajı döneriz
}

// UpdateProduct - ürün güncelleme handler'ı
func (h ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")                            //Url deki id parametrisini alırız
	objectID, err := primitive.ObjectIDFromHex(id) //aldığımız string tipindeki ıd değerini monodb id tipine dönüştürür
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz ID formatı"}) ///400 hata kodunu Json tipinde öner eror etiketiyle eror mesajını eşlerüiz
	}

	var updatedGame models.Game
	if err := c.Bind(&updatedGame); err != nil { //c.Bind http den gelen boudy yi gyani game nesnesinin json tipini &updategame in referansına atayabilirzse  err bil döner dmnemezse err hata mesajı döner
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"state": false}) //err  hata kodunu Json tipinde döner işlem gerçekleşmediği için statei false yaparız
	}
	result, err := h.Services.ProductUptade(objectID, updatedGame)
	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"state": false})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"state": true})
}

// PatchProduct - ürün kısmi güncelleme handler'ı
func (h ProductHandler) PatchProduct(c echo.Context) error {
	id := c.Param("id")                            //Url mizdeki strin id değerini alır
	objectID, err := primitive.ObjectIDFromHex(id) //ObjectIDFromHex string alınan id değerini mongodb id tipine dönüştürür
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"state": false, "error": "Geçersiz ID formatı"})
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
	return c.JSON(http.StatusOK, map[string]interface{}{"state": true})
}

// GetByID - ID'ye göre ürün getirme handler'ı
func (h ProductHandler) GetByID(c echo.Context) error {
	id := c.Param("id")                            //url deki id veri tipini alır
	objectID, err := primitive.ObjectIDFromHex(id) //alınan id strngini mongodbid tiine dönüştürür
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Geçersiz ID formatı"}) //400 hata kodunu Json tipinde öner eror mesajı eror mesajını eşlerüiz
	}
	result, err := h.Services.ProductGetByID(objectID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Oyun bulunamadı"}) //400 hata kodunu Json tipinde öner eror etiketiyle eror mesajını eşlerüiz
	}
	return c.JSON(http.StatusOK, result) //
}

// GetGamesSorted - fiyata göre sıralama handler'ı
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
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, games)
}

// GetGamesByExactName - tam isimle ürün arama handler'ı
func (h ProductHandler) GetGamesByExactName(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "İsim parametresi gereklidir"})
	}
	result, err := h.Services.ProductGetByExactName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Oyunlar getirilirken bir hata oluştu"})
	}
	return c.JSON(http.StatusOK, result)
}

// GetGamesByPartialName - kısmi isimle ürün arama handler'ı
func (h ProductHandler) GetGamesByPartialName(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "İsim parametresi gereklidir"})
	}
	result, err := h.Services.ProductGetByPartialName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Oyunlar getirilirken bir hata oluştu"})
	}
	return c.JSON(http.StatusOK, result)
}

//objectID, err := primitive.ObjectIDFromHex(id)//aldığımız string tipindeki ıd değerini monodb id tipine dönüştürür
//err := c.Bind(&updatedGame); err != nil {//c.Bind http den gelen boudy yi gyani game nesnesinin json tipini &updategame in referansına atayabilirzse  err bil döner dmnemezse err hata mesajı döner
//query := c.QueryParam("field")//QueryParam() sorgu parametresine verilen değeri almak için kulanılr  field a verilen değeri alır bunu artandan azalana yada azalandan artana sıralamak için kulanırız
