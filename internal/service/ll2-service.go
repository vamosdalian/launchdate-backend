package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/util"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LL2Service struct {
	mongoClient        *db.MongoDB
	LL2URLPrefix       string
	LL2RequestInterval int
}

func NewLL2Service(conf *config.Config, db *db.MongoDB) *LL2Service {
	return &LL2Service{
		mongoClient:        db,
		LL2URLPrefix:       conf.LL2URLPrefix,
		LL2RequestInterval: conf.LL2RequestInterval,
	}
}

// if async is true, function runs in background
// otherwise, it runs synchronously
func (s *LL2Service) UpdateLaunches(async bool) error {
	if async {
		go func() {
			err := s.updateLaunchesAsync()
			if err != nil {
				logrus.Errorf("%s", err)
			}
		}()
		return nil
	}
	return s.updateLaunchesAsync()
}

func (s *LL2Service) updateLaunchesAsync() error {
	count := 1
	offset := 0
	rl := util.NewRateLimit(time.Duration(s.LL2RequestInterval) * time.Second)
	logrus.Info("Starting LL2 launches update...")
	for {
		if offset >= count {
			break
		}

		rl.Wait()
		launches, err := s.LoadLaunches(100, offset)
		if err != nil {
			return err
		}
		count = launches.Count
		logrus.Infof("Fetched %d/%d launches from LL2", offset+len(launches.Results), count)

		for _, launch := range launches.Results {
			filter := map[string]any{
				"id": launch.ID,
			}
			update := map[string]any{
				"$set": launch,
			}
			opts := options.Update().SetUpsert(true)
			_, err := s.mongoClient.Collection("ll2_launch").UpdateOne(context.Background(), filter, update, opts)
			if err != nil {
				return err
			}
		}
		offset += len(launches.Results)
	}
	return nil
}

func (s *LL2Service) LoadLaunches(limit, offset int) (*models.LL2Response, error) {
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

	return launches, nil
}
