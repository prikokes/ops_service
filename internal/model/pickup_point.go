package model

import (
	"errors"
	"time"
)

// City represents cities where pickup points can be located
type City string

const (
	// CityMoscow represents Moscow city
	CityMoscow City = "Москва"
	// CitySaintPetersburg represents Saint Petersburg city
	CitySaintPetersburg City = "Санкт-Петербург"
	// CityKazan represents Kazan city
	CityKazan City = "Казань"
)

// AllowedCities is a list of cities where pickup points can be created
var AllowedCities = []City{CityMoscow, CitySaintPetersburg, CityKazan}

// ErrCityNotAllowed is returned when trying to create a pickup point in a non-allowed city
var ErrCityNotAllowed = errors.New("pickup points can only be created in Moscow, Saint Petersburg or Kazan")

// PickupPoint represents a pickup point (ПВЗ)
type PickupPoint struct {
	ID           int       `json:"id" db:"id"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
	City         City      `json:"city" db:"city"`
}

// IsCityAllowed checks if the city is allowed for creating pickup points
func IsCityAllowed(city City) bool {
	for _, c := range AllowedCities {
		if c == city {
			return true
		}
	}
	return false
} 