package services

import (
	"api-steam/dto"
	"api-steam/models"
	"api-steam/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductService ürün servisi için arayüz tanımlar bunuda repostroy katmanından verialarak yapar  ProductRepository den çekerek işlemi servies->Handeler a taşımak için kulanırız katmanına taşır
type ProductService interface {
	ProductInsert(product models.Game) (*dto.GameDTO, error)                          //veri eklemk
	ProductGetAll() ([]models.Game, error)                                            //Tüm verileri getirme
	ProductDelete(id primitive.ObjectID) (bool, error)                                //İD ye göre veri silme
	ProductUptade(id primitive.ObjectID, game models.Game) (bool, error)              //Veriyi komple günceleme
	ProductPatch(id primitive.ObjectID, updates map[string]interface{}) (bool, error) //Verilen bütünlüğü kadar günceleme
	ProductGetByID(id primitive.ObjectID) (models.Game, error)                        //Id ye göre arama
	ProductGetSorted(sortField string, order int) ([]models.Game, error)              //fiyata göre sıralama
	ProductGetByExactName(name string) ([]models.Game, error)                         //Tam isme göre arama
	ProductGetByPartialName(name string) ([]models.Game, error)                       //Kısmi isme göre arama
}

// DefaultProductService Repistory katmanında tanımladığımız fonksiyonları kulanmak için nesne türetme benzeri bir işlem
type DefaultProductService struct {
	Repo repository.ProductRepository
}

// ProductInsert ürün eklemek için servis işlemini gerçekleştirir
func (s *DefaultProductService) ProductInsert(product models.Game) (*dto.GameDTO, error) {
	var res dto.GameDTO
	result, err := s.Repo.Insert(product)
	if err != nil || !result {
		res.Status = false
		return &res, err
	}
	res = dto.GameDTO{Status: result}
	return &res, nil
}

// ürün listesini
func (s *DefaultProductService) ProductGetAll() ([]models.Game, error) { //servis katmanında repostory de tanımladığımız GetAll ukulanan fonksiyon
	result, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	return result, err
}

// ürün silme
func (s *DefaultProductService) ProductDelete(id primitive.ObjectID) (bool, error) {
	result, err := s.Repo.Delete(id)
	if err != nil || result == false {
		return false, err
	}
	return true, nil
}

// id ye göre ürünü kmple günceler
func (s *DefaultProductService) ProductUptade(id primitive.ObjectID, game models.Game) (bool, error) {
	result, err := s.Repo.Update(id, game)
	if err != nil || result == false {
		return false, err
	}
	return true, nil

}

// ProductPatch, bir ürünün belirli alanlarını günceller
func (s *DefaultProductService) ProductPatch(id primitive.ObjectID, updates map[string]interface{}) (bool, error) {
	// Repository katmanındaki Patch metodunu çağır
	result, err := s.Repo.Patch(id, updates)
	if err != nil {
		return false, err
	}

	return result, nil
}
func (s *DefaultProductService) ProductGetByID(id primitive.ObjectID) (models.Game, error) {
	result, err := s.Repo.GetByID(id)
	if err != nil {
		return models.Game{}, err //boş game ve hata döner
	}

	return result, nil
}

func (s *DefaultProductService) ProductGetSorted(sortField string, order int) ([]models.Game, error) { //Fiyata artan ve azalana göre sıralama
	result, err := s.Repo.GetAndSorted(sortField, order)
	if err != nil {
		return nil, err //boş games ve hata döner
	}

	return result, nil
}

// tam isme göre filtereleme
func (s *DefaultProductService) ProductGetByExactName(name string) ([]models.Game, error) {
	result, err := s.Repo.GetByExactName(name)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *DefaultProductService) ProductGetByPartialName(name string) ([]models.Game, error) {
	result, err := s.Repo.GetByPartialName(name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// NewProductService  servis katmanındakş funclarımı kulanabilmek içinb bir nesne türetme işlemi gibi
func NewProductService(repo repository.ProductRepository) ProductService {
	return &DefaultProductService{Repo: repo}
}
