package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type LL2Service struct {
	mongoClient  *db.MongoDB
	LL2URLPrefix string
}

func NewLL2Service(conf *config.Config, db *db.MongoDB) *LL2Service {
	return &LL2Service{
		mongoClient:  db,
		LL2URLPrefix: conf.LL2URLPrefix,
	}
}

func (s *LL2Service) LoadLaunches(limit, offset int) ([]*models.LL2LaunchDetailed, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url := fmt.Sprintf("%s/2.3.0/launches?limit=%d&offset=%d&mode=detailed", s.LL2URLPrefix, limit, offset)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var launches *models.LL2Response
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &launches)
	if err != nil {
		return nil, err
	}

	return launches.Results, nil
}
