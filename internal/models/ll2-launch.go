package models

type LL2Response struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []LL2Launch `json:"results"`
}

type LL2Launch struct {
	ID             string             `json:"id"`
	URL            string             `json:"url"`
	Slug           string             `json:"slug"`
	Name           string             `json:"name"`
	Status         LL2Status          `json:"status"`
	LastUpdated    string             `json:"last_updated"`
	Net            string             `json:"net"`
	WindowEnd      string             `json:"window_end"`
	WindowStart    string             `json:"window_start"`
	Image          string             `json:"image"`
	Infographic    string             `json:"infographic"`
	Mission        *LL2Mission        `json:"mission"`
	Pad            *LL2Pad            `json:"pad"`
	LaunchProvider *LL2LaunchProvider `json:"launch_service_provider"`
}

type LL2Status struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type LL2Mission struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Orbit       *LL2Orbit `json:"orbit"`
}

type LL2Orbit struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`
}

type LL2Pad struct {
	ID        int          `json:"id"`
	URL       string       `json:"url"`
	Name      string       `json:"name"`
	Latitude  string       `json:"latitude"`
	Longitude string       `json:"longitude"`
	Location  *LL2Location `json:"location"`
}

type LL2Location struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
}

type LL2LaunchProvider struct {
	ID   int    `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
}
