package services

import (
	"indicar-api/internal/domain/entities"

	"gorm.io/gorm"
)

type CityService struct {
	db *gorm.DB
}

func NewCityService(db *gorm.DB) *CityService {
	return &CityService{
		db: db,
	}
}

type CityResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Country string `json:"country"`
}

func (s *CityService) ListAll() ([]CityResponse, error) {
	var cities []entities.City
	if err := s.db.Order("name ASC").Find(&cities).Error; err != nil {
		return nil, err
	}

	response := make([]CityResponse, len(cities))
	for i, city := range cities {
		response[i] = CityResponse{
			ID:      city.ID,
			Name:    city.Name,
			State:   city.StateCode,
			Country: s.getCountryName(city.CountryCode),
		}
	}

	return response, nil
}

func (s *CityService) getCountryName(countryCode string) string {
	countryMap := map[string]string{
		"BR": "Brasil",
		"US": "Estados Unidos",
		"AR": "Argentina",
		// Adicione mais países conforme necessário
	}

	if name, exists := countryMap[countryCode]; exists {
		return name
	}
	return countryCode // Retorna o código se não encontrar o nome
}
