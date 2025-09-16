package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DataSourceConnector handles connections to external APIs
type DataSourceConnector struct {
	httpClient *http.Client
	mapsAPIKey string
	weatherKey string
	emtAPIKey  string
}

// NewDataSourceConnector creates a new data source connector
func NewDataSourceConnector(mapsAPIKey, weatherKey, emtAPIKey string) *DataSourceConnector {
	return &DataSourceConnector{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		mapsAPIKey: mapsAPIKey,
		weatherKey: weatherKey,
		emtAPIKey:  emtAPIKey,
	}
}

// Google Places API Response structures
type PlacesResponse struct {
	Results []PlaceResult `json:"results"`
	Status  string        `json:"status"`
}

type PlaceResult struct {
	PlaceID      string             `json:"place_id"`
	Name         string             `json:"name"`
	Types        []string           `json:"types"`
	Rating       float64            `json:"rating"`
	PriceLevel   int                `json:"price_level"`
	Geometry     PlaceGeometry      `json:"geometry"`
	OpeningHours *PlaceOpeningHours `json:"opening_hours,omitempty"`
	Photos       []PlacePhoto       `json:"photos,omitempty"`
	Vicinity     string             `json:"vicinity"`
}

type PlaceGeometry struct {
	Location PlaceLocation `json:"location"`
}

type PlaceLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type PlaceOpeningHours struct {
	OpenNow     bool     `json:"open_now"`
	WeekdayText []string `json:"weekday_text"`
}

type PlacePhoto struct {
	PhotoReference string `json:"photo_reference"`
	Height         int    `json:"height"`
	Width          int    `json:"width"`
}

// Weather API Response structures
type WeatherResponse struct {
	Current WeatherCurrent `json:"current"`
	Daily   []WeatherDaily `json:"daily"`
}

type WeatherCurrent struct {
	Temp      float64 `json:"temp"`
	Humidity  int     `json:"humidity"`
	WindSpeed float64 `json:"wind_speed"`
	Weather   []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

type WeatherDaily struct {
	Dt   int64 `json:"dt"`
	Temp struct {
		Day float64 `json:"day"`
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"temp"`
	Humidity int `json:"humidity"`
	Weather  []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

// FetchAttractions retrieves attractions from Google Places API
func (dsc *DataSourceConnector) FetchAttractions(ctx context.Context, destination string, interests []string) ([]Attraction, error) {
	if dsc.mapsAPIKey == "" {
		log.Println("Maps API key not configured, returning mock attractions")
		return dsc.getMockAttractions(destination), nil
	}

	var allAttractions []Attraction

	// Map interests to place types
	placeTypes := dsc.mapInterestsToPlaceTypes(interests)

	for _, placeType := range placeTypes {
		attractions, err := dsc.fetchAttractionsByType(ctx, destination, placeType)
		if err != nil {
			log.Printf("Error fetching attractions for type %s: %v", placeType, err)
			continue
		}
		allAttractions = append(allAttractions, attractions...)
	}

	// If no attractions found via API, return mock data
	if len(allAttractions) == 0 {
		return dsc.getMockAttractions(destination), nil
	}

	return allAttractions, nil
}

func (dsc *DataSourceConnector) fetchAttractionsByType(ctx context.Context, destination, placeType string) ([]Attraction, error) {
	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"

	params := url.Values{}
	params.Add("query", fmt.Sprintf("%s %s", placeType, destination))
	params.Add("key", dsc.mapsAPIKey)
	params.Add("type", placeType)

	resp, err := dsc.httpClient.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attractions: %v", err)
	}
	defer resp.Body.Close()

	var placesResp PlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&placesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if placesResp.Status != "OK" {
		return nil, fmt.Errorf("places API error: %s", placesResp.Status)
	}

	var attractions []Attraction
	for _, place := range placesResp.Results {
		attraction := Attraction{
			ID:   place.PlaceID,
			Name: place.Name,
			Type: dsc.mapPlaceTypeToCategory(place.Types),
			Location: Location{
				Latitude:  place.Geometry.Location.Lat,
				Longitude: place.Geometry.Location.Lng,
				Address:   place.Vicinity,
			},
			Rating:     place.Rating,
			PriceLevel: place.PriceLevel,
			Available:  true,
			Tags:       place.Types,
		}

		if place.OpeningHours != nil {
			attraction.OpeningHours = place.OpeningHours.WeekdayText
		}

		attractions = append(attractions, attraction)
	}

	return attractions, nil
}

// FetchWeather retrieves weather forecast
func (dsc *DataSourceConnector) FetchWeather(ctx context.Context, latitude, longitude float64) (*WeatherForecast, error) {
	if dsc.weatherKey == "" {
		log.Println("Weather API key not configured, returning mock weather")
		return dsc.getMockWeather(), nil
	}

	// Using OpenWeatherMap One Call API
	baseURL := "https://api.openweathermap.org/data/3.0/onecall"

	params := url.Values{}
	params.Add("lat", strconv.FormatFloat(latitude, 'f', 6, 64))
	params.Add("lon", strconv.FormatFloat(longitude, 'f', 6, 64))
	params.Add("appid", dsc.weatherKey)
	params.Add("units", "metric")
	params.Add("exclude", "minutely,hourly,alerts")

	resp, err := dsc.httpClient.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		log.Printf("Weather API error: %v", err)
		return dsc.getMockWeather(), nil
	}
	defer resp.Body.Close()

	var weatherResp WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		log.Printf("Weather decode error: %v", err)
		return dsc.getMockWeather(), nil
	}

	forecast := &WeatherForecast{
		Current: WeatherCondition{
			Date:        time.Now(),
			Temperature: weatherResp.Current.Temp,
			Humidity:    weatherResp.Current.Humidity,
			WindSpeed:   weatherResp.Current.WindSpeed,
		},
	}

	if len(weatherResp.Current.Weather) > 0 {
		forecast.Current.Description = weatherResp.Current.Weather[0].Description
		forecast.Current.Icon = weatherResp.Current.Weather[0].Icon
	}

	// Add daily forecast
	for _, daily := range weatherResp.Daily {
		condition := WeatherCondition{
			Date:        time.Unix(daily.Dt, 0),
			Temperature: daily.Temp.Day,
			Humidity:    daily.Humidity,
		}

		if len(daily.Weather) > 0 {
			condition.Description = daily.Weather[0].Description
			condition.Icon = daily.Weather[0].Icon
		}

		forecast.Forecast = append(forecast.Forecast, condition)
	}

	return forecast, nil
}

// FetchHotels retrieves hotel options
func (dsc *DataSourceConnector) FetchHotels(ctx context.Context, destination string, checkIn, checkOut time.Time, budget float64) ([]Hotel, error) {
	if dsc.mapsAPIKey == "" {
		log.Println("Maps API key not configured, returning mock hotels")
		return dsc.getMockHotels(destination), nil
	}

	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"

	params := url.Values{}
	params.Add("query", fmt.Sprintf("hotels in %s", destination))
	params.Add("key", dsc.mapsAPIKey)
	params.Add("type", "lodging")

	resp, err := dsc.httpClient.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		log.Printf("Hotels API error: %v", err)
		return dsc.getMockHotels(destination), nil
	}
	defer resp.Body.Close()

	var placesResp PlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&placesResp); err != nil {
		log.Printf("Hotels decode error: %v", err)
		return dsc.getMockHotels(destination), nil
	}

	var hotels []Hotel
	for _, place := range placesResp.Results {
		// Estimate price based on price_level and budget
		pricePerNight := dsc.estimateHotelPrice(place.PriceLevel, budget)

		hotel := Hotel{
			ID:   place.PlaceID,
			Name: place.Name,
			Location: Location{
				Latitude:  place.Geometry.Location.Lat,
				Longitude: place.Geometry.Location.Lng,
				Address:   place.Vicinity,
			},
			Rating:        place.Rating,
			PricePerNight: pricePerNight,
			Available:     true,
			Amenities:     []string{"WiFi", "Air Conditioning"}, // Default amenities
		}

		hotels = append(hotels, hotel)
	}

	if len(hotels) == 0 {
		return dsc.getMockHotels(destination), nil
	}

	return hotels, nil
}

// Helper methods

func (dsc *DataSourceConnector) mapInterestsToPlaceTypes(interests []string) []string {
	placeTypes := []string{"tourist_attraction"} // Default

	for _, interest := range interests {
		switch interest {
		case "culture", "history":
			placeTypes = append(placeTypes, "museum", "art_gallery")
		case "food", "dining":
			placeTypes = append(placeTypes, "restaurant", "cafe")
		case "nature", "outdoor":
			placeTypes = append(placeTypes, "park", "natural_feature")
		case "adventure", "activity":
			placeTypes = append(placeTypes, "amusement_park", "zoo")
		case "shopping":
			placeTypes = append(placeTypes, "shopping_mall", "store")
		case "nightlife":
			placeTypes = append(placeTypes, "night_club", "bar")
		}
	}

	return placeTypes
}

func (dsc *DataSourceConnector) mapPlaceTypeToCategory(types []string) string {
	for _, placeType := range types {
		switch placeType {
		case "museum", "art_gallery":
			return "culture"
		case "restaurant", "cafe", "food":
			return "dining"
		case "park", "natural_feature":
			return "nature"
		case "amusement_park", "zoo":
			return "entertainment"
		case "shopping_mall", "store":
			return "shopping"
		case "tourist_attraction":
			return "attraction"
		}
	}
	return "attraction"
}

func (dsc *DataSourceConnector) estimateHotelPrice(priceLevel int, budget float64) float64 {
	// Estimate based on price level (0-4) and user budget
	basePrices := []float64{50, 100, 150, 250, 400} // Price levels 0-4

	if priceLevel >= 0 && priceLevel < len(basePrices) {
		basePrice := basePrices[priceLevel]

		// Adjust based on budget
		if budget > 0 {
			budgetPerNight := budget / 7 // Assume 7 nights
			if budgetPerNight < basePrice {
				return budgetPerNight
			}
		}

		return basePrice
	}

	return 100 // Default price
}

// Mock data methods

func (dsc *DataSourceConnector) getMockAttractions(destination string) []Attraction {
	return []Attraction{
		{
			ID:          "mock_1",
			Name:        fmt.Sprintf("%s Historic Center", destination),
			Type:        "culture",
			Rating:      4.5,
			PriceLevel:  1,
			Description: "Historic downtown area with traditional architecture",
			Available:   true,
			Tags:        []string{"historic", "walking", "culture"},
		},
		{
			ID:          "mock_2",
			Name:        fmt.Sprintf("%s Art Museum", destination),
			Type:        "museum",
			Rating:      4.3,
			PriceLevel:  2,
			Description: "Local art and cultural exhibits",
			Available:   true,
			Tags:        []string{"art", "culture", "indoor"},
		},
	}
}

func (dsc *DataSourceConnector) getMockHotels(destination string) []Hotel {
	return []Hotel{
		{
			ID:            "hotel_1",
			Name:          fmt.Sprintf("Grand %s Hotel", destination),
			Rating:        4.2,
			PricePerNight: 120,
			Available:     true,
			Amenities:     []string{"WiFi", "Pool", "Restaurant", "Gym"},
		},
		{
			ID:            "hotel_2",
			Name:          fmt.Sprintf("%s Budget Inn", destination),
			Rating:        3.8,
			PricePerNight: 80,
			Available:     true,
			Amenities:     []string{"WiFi", "Parking"},
		},
	}
}

func (dsc *DataSourceConnector) getMockWeather() *WeatherForecast {
	now := time.Now()
	return &WeatherForecast{
		Current: WeatherCondition{
			Date:        now,
			Temperature: 22,
			Description: "Partly cloudy",
			Humidity:    65,
			WindSpeed:   10,
			Icon:        "02d",
		},
		Forecast: []WeatherCondition{
			{
				Date:        now.AddDate(0, 0, 1),
				Temperature: 24,
				Description: "Sunny",
				Humidity:    60,
				Icon:        "01d",
			},
			{
				Date:        now.AddDate(0, 0, 2),
				Temperature: 20,
				Description: "Light rain",
				Humidity:    80,
				Icon:        "10d",
			},
		},
	}
}
