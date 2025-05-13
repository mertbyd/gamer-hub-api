package repository

import (
	"api-steam/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepository arayüzü, ürün işlemleri için gereken metodları tanımlar
type ProductRepository interface {
	Insert(game models.Game) (bool, error)
	GetAll() ([]models.Game, error)
	Delete(id primitive.ObjectID) (bool, error) //Mongo db deki verimi primitive.ObjectID olduğu için primitive.ObjectID tipinde yolamam gerekiyor
	Update(id primitive.ObjectID, game models.Game) (bool, error)
	Patch(id primitive.ObjectID, updates map[string]interface{}) (bool, error)
	GetByID(id primitive.ObjectID) (models.Game, error)              // "*" eklendi
	GetAndSorted(sortField string, order int) ([]models.Game, error) //sortField string, order int   sortField=sıralamanın neye göre olcağı  order=+1 artana göre -1 azalana göre sıralalr
	GetByExactName(name string) ([]models.Game, error)
	GetByPartialName(name string) ([]models.Game, error)
	InsertMany(games []models.Game) (bool, error)
	GetByPriceRange(minPrice float64, maxPrice float64) ([]models.Game, error)
}

// ProductRepositoryDB, MongoDB işlemleri için collection(BAĞLANTI-DATABASE) ÇOK ALGILAYAMADIM
type ProductRepositoryDB struct {
	TodoCollection *mongo.Collection //mongo.Collection: MongoDB'deki bir koleksiyonu temsil eder (SQL'deki tabloya benzer)
}

// TodoCollection ile MongoDB koleksiyonuna erişim sağlayan repository nesnesini oluşturur
func NewProductRepository(dbClient *mongo.Collection) ProductRepository {
	return &ProductRepositoryDB{TodoCollection: dbClient}
}

// Veritabanına tek bir oyun ekler ve başarı durumunu döndürür
func (t *ProductRepositoryDB) Insert(game models.Game) (bool, error) { //t *ProductRepositoryDB bağlantı için reciver ettik
	// Gerekli alanları doldur
	game.ID = primitive.NewObjectID() //Mongo db nin kendi id si hariç bizim filtereememiz için benzersiz bir ıd atamada kulandık
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Fonksiyon bittiğinde context iptal edilir
	result, err := t.TodoCollection.InsertOne(ctx, game)
	if err != nil {
		log.Printf("Repository: Veritabanına oyun eklenirken hata: %v", err)
		return false, err
	}
	log.Printf("Repository: MongoDB'ye ekleme başarılı, ID: %v", result.InsertedID)
	return true, nil
}

// Veritabanına birden fazla oyun toplu olarak ekler ve başarı durumunu döndürür
func (t *ProductRepositoryDB) InsertMany(games []models.Game) (bool, error) { //t *ProductRepositoryDB bağlantı için reciver ettik
	var gamelist []interface{} //interface{} yapıyoruz ve yeni bir dizi oluşturuyoruz çünkü Insertmany interface{} istiyor
	for i := range games {
		games[i].ID = primitive.NewObjectID()
		games[i].CreatedAt = time.Now()
		games[i].UpdatedAt = time.Now()
		gamelist = append(gamelist, games[i])
	} //Mongo db nin kendi id si hariç bizim filtereememiz için benzersiz bir ıd atamada kulandık
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Fonksiyon bittiğinde context iptal edilir
	result, err := t.TodoCollection.InsertMany(ctx, gamelist)
	if err != nil {
		log.Printf("Repository: Toplu oyun eklenirken hata: %v", err)
		return false, err
	}
	log.Printf("Repository: MongoDB'ye ekleme başarılı, ID: %v", len(result.InsertedIDs))
	return true, nil
}

// Veritabanındaki tüm oyunları bir dizi olarak getirir
func (t *ProductRepositoryDB) GetAll() ([]models.Game, error) { //t *ProductRepositoryDB bağlantı için reciver ettik
	var game models.Game
	var games []models.Game
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()                                      // Fonksiyon bittiğinde context iptal edilir
	result, err := t.TodoCollection.Find(ctx, bson.M{}) //collection contextinden .find veri çekmek için kulanılır örnek dökümanı getrir
	if err != nil {
		log.Printf("Repository: Veritabanına oyun çekerken hata: %v", err)
		return nil, err
	}
	for result.Next(ctx) { //decode parça parça gelen veride gezinmek içn .Next() kulanılı pythondaki gibi44
		if err := result.Decode(&game); err != nil {
			log.Printf("Repository: Veritabanına oyun çekerken hata: %v", err)
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

// Belirtilen ID'ye sahip oyunu veritabanından siler
func (t *ProductRepositoryDB) Delete(id primitive.ObjectID) (bool, error) { //t *ProductRepositoryDB bağlantı için reciver ettik
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()                                                    // Fonksiyon bittiğinde context iptal edilir
	result, err := t.TodoCollection.DeleteOne(ctx, bson.M{"_id": id}) //colectionda bir nesne silmek için talep
	if err != nil || result.DeletedCount <= 0 {
		log.Printf("Repository: Veritabanına oyun çekerken hata: %v", err)
		return false, err
	}
	return true, nil
}

// Belirtilen ID'ye sahip oyunu tamamen günceller (PUT)
func (t *ProductRepositoryDB) Update(id primitive.ObjectID, game models.Game) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	game.ID = id                                                             //Güncelenecek objenin ıd si değişmemeli
	game.UpdatedAt = time.Now()                                              // Güncelleme zamanını güncelle
	result, err := t.TodoCollection.ReplaceOne(ctx, bson.M{"_id": id}, game) //ReplaceOne ile belgenin tamamını güncele
	if err != nil {
		log.Printf("Repository: Veritabanında oyun güncellenirken hata: %v", err)
		return false, err
	}
	if result.MatchedCount <= 0 { //gÜNCELENEN BELGE SAYISI KONTROLÜ
		log.Printf("Repository: Güncellenecek oyun bulunamadı, ID: %v", id)
		return false, nil
	}
	return true, nil
}

// Belirtilen ID'ye sahip oyunun sadece belirli alanlarını günceller (PATCH)
func (t *ProductRepositoryDB) Patch(id primitive.ObjectID, updates map[string]interface{}) (bool, error) {
	// Context tanımlama eklendi
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updates["updatedAt"] = time.Now()                                                          //güncelenme tarihini değişirmek için
	result, err := t.TodoCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates}) //patch işlemi için Updateone komutunu kulandık
	if err != nil {
		log.Printf("Repository: Veritabanında oyun güncellenirken hata: %v", err)
		return false, err
	}
	if result.MatchedCount <= 0 {
		log.Printf("Repository: Güncellenecek oyun bulunamadı, ID: %v", id)
		return false, nil
	}
	log.Printf("Repository: MongoDB'de kısmi güncelleme başarılı, ID: %v", id)
	return true, nil
}

// Belirtilen ID'ye göre tek bir oyun verisini getirir
func (t *ProductRepositoryDB) GetByID(id primitive.ObjectID) (models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var game models.Game
	err := t.TodoCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&game) //FindOne verilen id ye göre bulur ve decode ile tanımlanan game in referans adresine atar
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Repository: ID'si %v olan oyun bulunamadı", id)
			return models.Game{}, fmt.Errorf("oyun bulunamadı") // Boş game ve hata döndür
		}
		log.Printf("Repository: ID'si %v olan oyun getirilirken hata: %v", id, err)
		return models.Game{}, err // Boş game ve hata döndür
	}
	log.Printf("Repository: ID'si %v olan oyun başarıyla getirildi", id)
	return game, nil // Game ve nil hata döndür
}

// Oyunları belirtilen alana göre artana veya azalana sıralayarak getirir
func (t *ProductRepositoryDB) GetAndSorted(sortField string, order int) ([]models.Game, error) { //sortField string, order int   sortField=sıralamanın neye göre olcağı  order=+1 artana göre -1 azalana göre sıralalr
	var game models.Game
	var games []models.Game
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: sortField, Value: order}}) //options.Find() sorgu yaparken sıralama yapabileceğimiz ekseçenekler sunar SetSort=parametreye göre sıralama  sortField string, order int   sortField=sıralamanın neye göre olcağı  order=+1 artana göre -1 azalana göre sıralalr
	result, err := t.TodoCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Printf("Repository: Verileri sıralarken hata: %v", err)
		return nil, err
	}
	for result.Next(ctx) { //result ile next ile nesnelerde
		if err := result.Decode(&game); err != nil { //referansına atama
			log.Printf("Repository: Veritabanına oyun çekerken hata: %v", err)
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

// Tam olarak eşleşen isme sahip oyunları getirir
func (t *ProductRepositoryDB) GetByExactName(name string) ([]models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"title": name} //MongoDB sorguları oluşturmak için kullanılan bir filtre nesnesidir
	result, err := t.TodoCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Repository: Tam isim ile sorgu sırasında hata: %v", err)
		return nil, err
	}
	var games []models.Game
	if err = result.All(ctx, &games); err != nil {
		log.Printf("Repository: Sonuçları okurken hata: %v", err)
		return nil, err
	}
	return games, nil
}

// İsmin bir kısmıyla eşleşen oyunları getirir (regex kullanarak)
func (t *ProductRepositoryDB) GetByPartialName(name string) ([]models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var game models.Game
	var games []models.Game
	//YAPAY ZEKA
	// Kısmi eşleşme için regex kullan
	regexPattern := primitive.Regex{
		Pattern: name,
		Options: "i", // i = case insensitive (büyük/küçük harf duyarsız)
	}
	//YAPAY ZEKA
	filter := bson.M{"title": regexPattern}
	result, err := t.TodoCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Repository: Kısmi isim ile sorgu sırasında hata: %v", err)
		return nil, err
	}
	for result.Next(ctx) { //result a gelen nesnelerde next ile gezindik
		if err := result.Decode(&game); err != nil { //Decode ile game değişkenin referansına atadık
			log.Printf("Repository: Veritabanına oyun çekerken hata: %v", err)
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

// fiyat aralığına göre filtreleme
func (t *ProductRepositoryDB) GetByPriceRange(minPrice float64, maxPrice float64) ([]models.Game, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var games []models.Game
	var game models.Game
	//Yapay Zeka
	filter := bson.M{
		"price.amount": bson.M{
			"$gte": minPrice,
			"$lte": maxPrice,
		},
	}
	//Yappay Zeka
	result, err := t.TodoCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Repository: Fiyat aralığına göre sorgulama hatası: %v", err)
		return nil, err
	}
	for result.Next(ctx) {
		if result.Decode(&game); err != nil {
			log.Printf("Repository: Oyun verisi çözümlenirken hata: %v", err)
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

//InsertOne() mongodb de 1 tane veri eklemek için
//Find() veri çekmek için
//DeleteOne() veri silmek için
//ReplaceOne Komple veriyi güncelemek için ReplaceOne
//UpdateOne Patch işleminde tek veri güncelmek için
//FindOne() 1 tane veri çekmek için
//Find().SetSort(bson.D{{Key: sortField, Value: order}}) find ile gelen verileri  sortField sıralanack parametre Value sıralama tipi
/*Context (Bağlam) ve ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second):
Context, Go'da işlemleri kontrol etmek, iptal etmek veya zaman aşımına uğratmak için kullanılan bir yapıdır. Özellikle:
context.Background(): Boş bir ana context oluşturur
context.WithTimeout(): Belirtilen süre sonunda otomatik iptal olacak bir context oluşturur
defer cancel(): İşlem tamamlandığında veya hata oluştuğunda context'i iptal eder
Bu yapıyı MongoDB işlemlerinde kullanıyorsunuz çünkü:
Veritabanı işlemleri sonsuza kadar asılı kalmasın diye (10 saniye zaman aşımı)
İşlem tamamlandığında kaynakların serbest bırakılması için
Uzun süren işlemlerin kontrollü bir şekilde sonlandırılması için
TodoCollection *mongo.Collection ve Client:
Bu yapı MongoDB ile etkileşim kurmanın temelini oluşturur:
mongo.Collection: MongoDB'deki bir koleksiyonu temsil eder (SQL'deki tabloya benzer)
TodoCollection: Repository struct'ınızın içinde MongoDB koleksiyonuna bir referans tutmak için kullanılır
dbClient: NewProductRepository fonksiyonuna geçirilen parametre, ana uygulamanızda oluşturulan MongoDB bağlantısından gelen koleksiyon referansıdır
main.go dosyanızda muhtemelen şöyle bir kod vardır:
go
Kopyala
1
2
client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
dbClient := client.Database("steamdb").Collection("games")
Bu yapıyı kullanmanızın nedeni:
Bağımlılık enjeksiyonu (dependency injection) prensibi: Repository doğrudan veritabanına bağlanmak yerine, dışarıdan bir bağlantı alır
Test edilebilirlik: Gerçek veritabanı yerine sahte (mock) bir koleksiyon kullanarak testler yazabilirsiniz
Esneklik: Aynı repository yapısını farklı koleksiyonlarla kullanabilirsiniz
Bu yapı, Clean Architecture prensiplerine uygun olarak yazılmış ve katmanlar arası bağımlılıkları düzgün bir şekilde yönetmenizi sağlıyor.
*/

/*
claude-3.7-sonnet-latest
Today at 00:37
Go dilinde * ve & operatörleri bellek yönetimi ve pointer manipülasyonu için kullanılır. Bu temel kavramları açıklayayım:
& (Adres Operatörü)
& işareti, bir değişkenin bellek adresini alır. Yani, değişkenin değerini değil, o değişkenin bellekte nerede saklandığını verir.
go
Kaydet
Kopyala
1
2
var x int = 10
var p *int = &x  // p, x'in bellek adresini tutar
Burada &x, x değişkeninin bellek adresini te
*/
