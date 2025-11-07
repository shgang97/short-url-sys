package geoip

import (
	"fmt"

	"github.com/ip2location/ip2location-go"
)

type Service interface {
	GetGeoInfo(ip string) (*GeoInfo, error)
}

type GeoInfo struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

type DefaultGeoIPService struct {
	db *ip2location.DB
}

func New(dbPath string) (*DefaultGeoIPService, error) {
	db, err := ip2location.OpenDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GeoIP database: %w", err)
	}
	return &DefaultGeoIPService{db: db}, nil
}

func (s *DefaultGeoIPService) GetGeoInfo(ip string) (*GeoInfo, error) {
	record, err := s.db.Get_all(ip)
	if err != nil {
		return nil, err
	}

	return &GeoInfo{
		Country: record.Country_short,
		Region:  record.Region,
		City:    record.City,
	}, nil
}
