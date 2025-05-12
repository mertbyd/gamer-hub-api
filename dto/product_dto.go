package dto

type GameDTO struct { //servis işlemleri yapışdığında bool nesnesi döncek bunu  Status olarak alıcaz ve servis dosyamızda kontrol edicez)
	Status bool `json:"status,omitempty"`
}

/*
Burda kurduumuz dto aslında postmandan gerçekleştirdiğimiiz işlemlerde bize true fase dönerek çalışıp çalışmadığını anlicaz
*/
