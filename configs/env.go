package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	/*çalıştığı dizindeki .env dosyasını arar içindeki ortam dğişkenlerini mevcut ortamda kulanılablir hale getirir*/
	/* .env dosyası
	   uygulamanın yapılandırma ayarlarını ve hassas bilgilerinin saklandığı dosyadır
	*/
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	mongoURI := os.Getenv("MONGOURI")
	if mongoURI == "" {
		log.Println("MONGOURI environment variable is not set, using default connection string")
		mongoURI = "mongodb://localhost:27017/GameApi"
	}

	fmt.Println(mongoURI)
	return mongoURI
}
