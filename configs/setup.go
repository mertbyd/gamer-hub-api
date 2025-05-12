package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	/*Context  GO DA İŞLEMLERİN SINIRLARINI BELİRLEMEK ZAMAN AŞIMI GİBİ BİR SİSTEM UYGULAYIP SONLANDIRMAYA YARAR*/
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // İşlem bittiğinde context'i iptal et
	/*
	   Çalışması için bulunduğu fonksiyonu sonuna kadar erteler yani deffer ile çağırıln fonksiyonun sonuna kadar erteler
	*/
	// Tek adımda MongoDB bağlantısı kurma (önerilen yaklaşım)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatalf("MongoDB'ye bağlanılamadı: %v", err)
	}
	//*client.pig ile  MongoDB bağlantısının  çalışıp çalışmadığını kontrol ederiz*/
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping testi başarısız: %v", err)
	}
	log.Println("MongoDB bağlantısı başarıyla kuruldu")

	// Veritabanı ve koleksiyonu başlangıçta oluştur fonksiyonu yorum satırına alındı
	// initDatabase(client)

	return client
}

// *NewClient ile yapılandırdığımız  bağlantıyı Connect ile başlatık ctx değişkeni ile zaman aşımı ayarladık*/
var DB *mongo.Client = ConnectDB()

// GetCollection belirtilen koleksiyonu döndürür
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("GameApi").Collection(collectionName)
	// client.Database() veritabanına erişim sağlar
	// .Collection() belirtilen koleksiyona erişim sağlar (SQL'deki tablo benzeri yapı)
}

// ------------- YAPAY ZEKA -------------

// initDatabase veritabanı ve koleksiyonu başlangıçta oluşturur
/*
func initDatabase(client *mongo.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    // Veritabanını ve koleksiyonu oluştur
    db := client.Database("GameApi")
    // Koleksiyonun varlığını kontrol et ve gerekirse oluştur
    collections, err := db.ListCollectionNames(ctx, bson.M{})
    if err != nil {
        log.Printf("Koleksiyon listesi alınamadı: %v", err)
    }
    // "games" koleksiyonu var mı kontrol et
    collectionExists := false
    for _, collection := range collections {
        if collection == "games" {
            collectionExists = true
            break
        }
    }
    // Koleksiyon yoksa oluştur
    if !collectionExists {
        err = db.CreateCollection(ctx, "games")
        if err != nil {
            log.Printf("Koleksiyon oluşturma hatası: %v", err)
        } else {
            log.Println("'games' koleksiyonu başarıyla oluşturuldu")
        }
    } else {
        log.Println("'games' koleksiyonu zaten mevcut")
    }
}
*/
