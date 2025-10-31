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

const LL2COLLECTION = "ll2_launch"

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
		launches, err := s.LoadLaunches(10, offset)
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
			_, err := s.mongoClient.Collection(LL2COLLECTION).UpdateOne(context.Background(), filter, update, opts)
			if err != nil {
				return err
			}
		}
		offset += len(launches.Results)
	}
	return nil
}

func (s *LL2Service) LoadLaunches(limit, offset int) (*models.LL2Response, error) {
	var launches *models.LL2Response
	err := s.LoadDataFromAPI("launches", limit, offset, launches)

	return launches, err
}

func (s *LL2Service) GetLaunchesFromDB(limit, offset int) ([]models.LL2LaunchNormal, error) {
	collection := s.mongoClient.Collection(LL2COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(map[string]int{"net": 1}) // Sort by net ascending

	cursor, err := collection.Find(ctx, map[string]any{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var launches []models.LL2LaunchNormal
	for cursor.Next(ctx) {
		var launch models.LL2LaunchNormal
		if err := cursor.Decode(&launch); err != nil {
			return nil, err
		}
		launches = append(launches, launch)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return launches, nil
}

func (s *LL2Service) LoadAngecyFromAPI(limit, offset int) (*models.LL2AngecyResponse, error) {
	var launches *models.LL2AngecyResponse
	err := s.LoadDataFromAPI("agencies", limit, offset, &launches)
	return launches, err
}

func (s *LL2Service) LoadDataFromAPI(endpoint string, limit, offset int, payload any) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url := fmt.Sprintf("%s/2.3.0/%s?limit=%d&offset=%d&mode=detailed", s.LL2URLPrefix, endpoint, limit, offset)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LL2 API returned status code %d, url:%s", resp.StatusCode, url)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *LL2Service) GetAngecyFromDB(limit, offset int) ([]models.LL2AgencyDetailed, error) {
	collection := s.mongoClient.Collection("ll2_agency")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(map[string]int{"id": 1}) // Sort by id ascending

	cursor, err := collection.Find(ctx, map[string]any{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var agencies []models.LL2AgencyDetailed
	for cursor.Next(ctx) {
		var agency models.LL2AgencyDetailed
		if err := cursor.Decode(&agency); err != nil {
			return nil, err
		}
		agencies = append(agencies, agency)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return agencies, nil
}

func (s *LL2Service) UpdateAngecy(async bool) error {
	if async {
		go func() {
			err := s.updateAngecyAsync()
			if err != nil {
				logrus.Errorf("%s", err)
			}
		}()
		return nil
	}
	return s.updateAngecyAsync()
}

func (s *LL2Service) updateAngecyAsync() error {
	count := 10
	offset := 0
	rl := util.NewRateLimit(time.Duration(s.LL2RequestInterval) * time.Second)
	logrus.Info("Starting LL2 angecy update...")
	for {
		if offset >= count {
			break
		}

		rl.Wait()
		agencies, err := s.LoadAngecyFromAPI(10, offset)
		if err != nil {
			return err
		}
		count = agencies.Count
		logrus.Infof("Fetched %d/%d angecies from LL2", offset+len(agencies.Results), count)

		for _, agency := range agencies.Results {
			filter := map[string]any{
				"id": agency.ID,
			}
			update := map[string]any{
				"$set": agency,
			}
			opts := options.Update().SetUpsert(true)
			_, err := s.mongoClient.Collection("ll2_agency").UpdateOne(context.Background(), filter, update, opts)
			if err != nil {
				return err
			}
		}
		offset += len(agencies.Results)
	}
	return nil
}
