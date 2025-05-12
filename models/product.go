package models

/*
'json:"....."'    ->json formatına dönüştüür
`bson:"alan_adi"` ->MongoDB veritabanına veri kaydederken ve veri okurken kullanılır
  `json:"-" bson:"password"`  // JSON'da gösterilmez
*/

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Platform, oyunun çalıştığı platformları temsil eder
type Platform struct {
	Name         string `json:"name" bson:"name"`                                     // Platform adı (PC, PS5, Xbox Series X, vb.)
	ReleaseDate  string `json:"release_date,omitempty" bson:"release_date,omitempty"` // Bu platformda çıkış tarihi
	Requirements string `json:"requirements,omitempty" bson:"requirements,omitempty"` // Platform gereksinimleri
}

// Requirements, belirli bir donanım gereksinim setini temsil eder
type Requirements struct {
	OS              string `json:"os,omitempty" bson:"os,omitempty"`                             // İşletim sistemi gereksinimleri
	Processor       string `json:"processor,omitempty" bson:"processor,omitempty"`               // İşlemci gereksinimleri
	Memory          string `json:"memory,omitempty" bson:"memory,omitempty"`                     // Bellek gereksinimleri
	Graphics        string `json:"graphics,omitempty" bson:"graphics,omitempty"`                 // Ekran kartı gereksinimleri
	DirectX         string `json:"directx,omitempty" bson:"directx,omitempty"`                   // DirectX gereksinimleri
	Storage         string `json:"storage,omitempty" bson:"storage,omitempty"`                   // Depolama gereksinimleri
	AdditionalNotes string `json:"additional_notes,omitempty" bson:"additional_notes,omitempty"` // Ek notlar
}

// SystemRequirements, oyunun sistem gereksinimlerini temsil eder
type SystemRequirements struct {
	Minimum     Requirements `json:"minimum,omitempty" bson:"minimum,omitempty"`         // Minimum gereksinimler
	Recommended Requirements `json:"recommended,omitempty" bson:"recommended,omitempty"` // Önerilen gereksinimler
}

// Genre, oyun türlerini temsil eder
type Genre struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                  // Tür ID'si
	Name        string             `json:"name" bson:"name"`                                   // Tür adı (Aksiyon, Macera, vb.)
	Description string             `json:"description,omitempty" bson:"description,omitempty"` // Tür açıklaması
}

// Publisher, oyun yayıncısını temsil eder
type Publisher struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                    // Yayıncı ID'si
	Name        string             `json:"name" bson:"name"`                                     // Yayıncı adı
	Country     string             `json:"country,omitempty" bson:"country,omitempty"`           // Ülke
	Website     string             `json:"website,omitempty" bson:"website,omitempty"`           // Web sitesi
	Description string             `json:"description,omitempty" bson:"description,omitempty"`   // Açıklama
	FoundedYear int                `json:"founded_year,omitempty" bson:"founded_year,omitempty"` // Kuruluş yılı
}

// Developer, oyun geliştiricisini temsil eder
type Developer struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                    // Geliştirici ID'si
	Name        string             `json:"name" bson:"name"`                                     // Geliştirici adı
	Country     string             `json:"country,omitempty" bson:"country,omitempty"`           // Ülke
	Website     string             `json:"website,omitempty" bson:"website,omitempty"`           // Web sitesi
	Description string             `json:"description,omitempty" bson:"description,omitempty"`   // Açıklama
	FoundedYear int                `json:"founded_year,omitempty" bson:"founded_year,omitempty"` // Kuruluş yılı
}

// Price, oyun fiyat bilgilerini temsil eder
type Price struct {
	Amount      float64   `json:"amount" bson:"amount"`                                   // Fiyat miktarı
	Currency    string    `json:"currency" bson:"currency"`                               // Para birimi (USD, EUR, TRY vb.)
	Discount    float64   `json:"discount,omitempty" bson:"discount,omitempty"`           // İndirim oranı (0-1 arası)
	OnSale      bool      `json:"on_sale" bson:"on_sale"`                                 // İndirimde mi?
	SaleEndDate time.Time `json:"sale_end_date,omitempty" bson:"sale_end_date,omitempty"` // İndirim bitiş tarihi
}

// Media, oyun medya içeriklerini temsil eder
type Media struct {
	CoverImage   string   `json:"cover_image" bson:"cover_image"`                         // Kapak resmi
	Screenshots  []string `json:"screenshots,omitempty" bson:"screenshots,omitempty"`     // Ekran görüntüleri
	Videos       []string `json:"videos,omitempty" bson:"videos,omitempty"`               // Video URL'leri
	Trailers     []string `json:"trailers,omitempty" bson:"trailers,omitempty"`           // Fragman URL'leri
	ThumbnailURL string   `json:"thumbnail_url,omitempty" bson:"thumbnail_url,omitempty"` // Küçük resim URL'si
	BannerImage  string   `json:"banner_image,omitempty" bson:"banner_image,omitempty"`   // Banner resmi
}

// Rating, oyun derecelendirmelerini temsil eder
type Rating struct {
	AverageScore       float64 `json:"average_score,omitempty" bson:"average_score,omitempty"`             // Ortalama puan (0-10 arası)
	TotalReviews       int     `json:"total_reviews,omitempty" bson:"total_reviews,omitempty"`             // Toplam değerlendirme sayısı
	PositivePercentage int     `json:"positive_percentage,omitempty" bson:"positive_percentage,omitempty"` // Olumlu değerlendirme yüzdesi
	ESRB               string  `json:"esrb,omitempty" bson:"esrb,omitempty"`                               // ESRB derecesi (E, T, M, vb.)
	PEGI               string  `json:"pegi,omitempty" bson:"pegi,omitempty"`                               // PEGI derecesi (3, 7, 12, 16, 18)
}

// Game, API'deki ana oyun verisini temsil eder
type Game struct {
	ID               primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`                                  // Benzersiz tanımlayıcı
	Title            string               `json:"title" bson:"title"`                                                 // Oyun adı
	Description      string               `json:"description,omitempty" bson:"description,omitempty"`                 // Açıklama
	ShortDescription string               `json:"short_description,omitempty" bson:"short_description,omitempty"`     // Kısa açıklama
	ReleaseDate      time.Time            `json:"release_date" bson:"release_date"`                                   // Yayın tarihi
	Developers       []Developer          `json:"developers,omitempty" bson:"developers,omitempty"`                   // Geliştiriciler
	Publishers       []Publisher          `json:"publishers,omitempty" bson:"publishers,omitempty"`                   // Yayıncılar
	Genres           []Genre              `json:"genres,omitempty" bson:"genres,omitempty"`                           // Türler
	Tags             []string             `json:"tags,omitempty" bson:"tags,omitempty"`                               // Etiketler
	Platforms        []Platform           `json:"platforms,omitempty" bson:"platforms,omitempty"`                     // Platformlar
	SystemReqs       SystemRequirements   `json:"system_requirements,omitempty" bson:"system_requirements,omitempty"` // Sistem gereksinimleri
	Price            Price                `json:"price" bson:"price"`                                                 // Fiyat bilgileri
	Media            Media                `json:"media" bson:"media"`                                                 // Medya içerikleri
	Rating           Rating               `json:"rating,omitempty" bson:"rating,omitempty"`                           // Değerlendirme bilgileri
	Features         []string             `json:"features,omitempty" bson:"features,omitempty"`                       // Özellikler (çok oyunculu, bulut kaydetme, vb.)
	Languages        []string             `json:"languages,omitempty" bson:"languages,omitempty"`                     // Desteklenen diller
	IsEarlyAccess    bool                 `json:"is_early_access" bson:"is_early_access"`                             // Erken erişimde mi?
	IsMultiplayer    bool                 `json:"is_multiplayer" bson:"is_multiplayer"`                               // Çok oyunculu mu?
	TotalPlayTime    int                  `json:"total_playtime,omitempty" bson:"total_playtime,omitempty"`           // Ortalama oynanış süresi (dakika)
	SimilarGames     []primitive.ObjectID `json:"similar_games,omitempty" bson:"similar_games,omitempty"`             // Benzer oyunların ID'leri
	CreatedAt        time.Time            `json:"created_at" bson:"created_at"`                                       // Veritabanına eklenme tarihi
	UpdatedAt        time.Time            `json:"updated_at" bson:"updated_at"`                                       // Son güncelleme tarihi
	Status           string               `json:"status" bson:"status"`                                               // Oyunun durumu (active, coming_soon, removed, vb.)
}
