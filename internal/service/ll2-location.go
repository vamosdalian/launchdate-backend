package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/util"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *LL2Service) GetLocationsFromApi(limit, offset int) (*models.LL2LocationResponse, error) {
	var locations *models.LL2LocationResponse
	err := s.LoadDataFromAPI("locations", limit, offset, &locations)
	return locations, err
}

func (s *LL2Service) GetLocationsFromDB(limit, offset int) ([]models.LL2LocationSerializerWithPads, error) {
	var locations []models.LL2LocationSerializerWithPads
	collection := s.mongoClient.Collection("ll2_location")
	ctx := context.Background()
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(map[string]int{"id": 1})
	cursor, err := collection.Find(ctx, struct{}{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var location models.LL2LocationSerializerWithPads
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func (s *LL2Service) UpdateLocationsAsync(async bool) error {
	if async {
		go func() {
			err := s.updateLocations()
			if err != nil {
				logrus.Errorf("Failed to update LL2 locations: %s", err)
			}
		}()
		return nil
	}
	return s.updateLocations()
}

func (s *LL2Service) updateLocations() error {
	limit := 10
	offset := 0
	rl := util.NewRateLimit(time.Duration(s.LL2RequestInterval) * time.Second)
	logrus.Info("Starting LL2 location update...")
	for {
		rl.Wait()
		locations, err := s.GetLocationsFromApi(limit, offset)
		if err != nil {
			return err
		}
		if len(locations.Results) == 0 {
			break
		}
		logrus.Infof("Fetched %d/%d locations from LL2", offset+len(locations.Results), locations.Count)
		for _, location := range locations.Results {
			filter := map[string]any{
				"id": location.ID,
			}
			update := map[string]any{
				"$set": location,
			}
			opts := options.Update().SetUpsert(true)
			collection := s.mongoClient.Collection("ll2_location")
			_, err := collection.UpdateOne(context.Background(), filter, update, opts)
			if err != nil {
				return err
			}
		}
		offset += len(locations.Results)
	}
	return nil
}

func (s *LL2Service) LoadPadsFromAPI(limit, offset int) (*models.LL2PadResponse, error) {
	var pads *models.LL2PadResponse
	err := s.LoadDataFromAPI("pads", limit, offset, &pads)
	return pads, err
}

func (s *LL2Service) GetPadsFromDB(limit, offset int) ([]models.LL2Pad, error) {
	var pads []models.LL2Pad
	collection := s.mongoClient.Collection("ll2_pad")
	ctx := context.Background()
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(map[string]int{"id": 1})
	cursor, err := collection.Find(ctx, struct{}{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var pad models.LL2Pad
		if err := cursor.Decode(&pad); err != nil {
			return nil, err
		}
		pads = append(pads, pad)
	}
	return pads, nil
}

func (s *LL2Service) UpdatePadsAsync(async bool) error {
	if async {
		go func() {
			err := s.updatePads()
			if err != nil {
				logrus.Errorf("Failed to update LL2 pads: %s", err)
			}
		}()
		return nil
	}
	return s.updatePads()
}

func (s *LL2Service) updatePads() error {
	limit := 10
	offset := 0
	rl := util.NewRateLimit(time.Duration(s.LL2RequestInterval) * time.Second)
	logrus.Info("Starting LL2 pad update...")
	for {
		rl.Wait()
		pads, err := s.LoadPadsFromAPI(limit, offset)
		if err != nil {
			return err
		}
		if len(pads.Results) == 0 {
			break
		}
		logrus.Infof("Fetched %d/%d pads from LL2", offset+len(pads.Results), pads.Count)
		for _, pad := range pads.Results {
			filter := map[string]any{
				"id": pad.Id,
			}
			update := map[string]any{
				"$set": pad,
			}
			opts := options.Update().SetUpsert(true)
			collection := s.mongoClient.Collection("ll2_pad")
			_, err := collection.UpdateOne(context.Background(), filter, update, opts)
			if err != nil {
				return err
			}
		}
		offset += len(pads.Results)
	}
	return nil
}
