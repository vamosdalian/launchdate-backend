package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/util"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *LL2Service) LoadLaunchersFromAPI(limit, offset int) (*models.LL2LauncherResponse, error) {
	var launches *models.LL2LauncherResponse
	err := s.LoadDataFromAPI("launcher_configurations", limit, offset, &launches)
	return launches, err
}

func (s *LL2Service) GetLaunchersFromDB(limit, offset int) ([]models.LL2LauncherConfigNormal, error) {
	var launchers []models.LL2LauncherConfigNormal
	collection := s.mongoClient.Collection("ll2_launcher")
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
		var launcher models.LL2LauncherConfigNormal
		if err := cursor.Decode(&launcher); err != nil {
			return nil, err
		}
		launchers = append(launchers, launcher)
	}
	return launchers, nil
}

func (s *LL2Service) UpdateLaunchersAsync(async bool) error {
	if async {
		go func() {
			err := s.updateLaunchers()
			if err != nil {
				logrus.Errorf("Failed to update LL2 launchers: %s", err)
			}
		}()
		return nil
	}
	return s.updateLaunchers()
}

func (s *LL2Service) updateLaunchers() error {
	limit := 10
	offset := 0
	rl := util.NewRateLimit(time.Duration(s.LL2RequestInterval) * time.Second)
	logrus.Info("Starting LL2 launcher update...")
	for {
		rl.Wait()
		launchers, err := s.LoadLaunchersFromAPI(limit, offset)
		if err != nil {
			return err
		}
		if len(launchers.Results) == 0 {
			break
		}
		logrus.Infof("Fetched %d/%d launches from LL2", offset+len(launchers.Results), launchers.Count)
		for _, launcher := range launchers.Results {
			filter := map[string]any{
				"id": launcher.ID,
			}
			update := map[string]any{
				"$set": launcher,
			}
			opts := options.Update().SetUpsert(true)
			_, err := s.mongoClient.Collection("ll2_launcher").UpdateOne(context.Background(), filter, update, opts)
			if err != nil {
				return err
			}
		}
		offset += len(launchers.Results)
	}
	return nil
}

func (s *LL2Service) LoadLauncherFamiliesFromAPI(limit, offset int) (*models.LL2LauncherFamilyResponse, error) {
	var families *models.LL2LauncherFamilyResponse
	err := s.LoadDataFromAPI("launcher_configuration_families", limit, offset, &families)
	return families, err
}

func (s *LL2Service) GetLauncherFamiliesFromDB(limit, offset int) ([]models.LL2LauncherConfigFamilyDetailed, error) {
	var families []models.LL2LauncherConfigFamilyDetailed
	collection := s.mongoClient.Collection("ll2_launcher_family")
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
		var family models.LL2LauncherConfigFamilyDetailed
		if err := cursor.Decode(&family); err != nil {
			return nil, err
		}
		families = append(families, family)
	}
	return families, nil
}

func (s *LL2Service) UpdateLauncherFamiliesAsync(async bool) error {
	if async {
		go func() {
			err := s.updateLauncherFamilies()
			if err != nil {
				logrus.Errorf("Failed to update LL2 launcher families: %s", err)
			}
		}()
		return nil
	}
	return s.updateLauncherFamilies()
}

func (s *LL2Service) updateLauncherFamilies() error {
	limit := 10
	offset := 0
	rl := util.NewRateLimit(time.Duration(s.LL2RequestInterval) * time.Second)
	logrus.Info("Starting LL2 launcher family update...")
	for {
		rl.Wait()
		families, err := s.LoadLauncherFamiliesFromAPI(limit, offset)
		if err != nil {
			return err
		}
		if len(families.Results) == 0 {
			break
		}
		logrus.Infof("Fetched %d/%d launcher families from LL2", offset+len(families.Results), families.Count)
		for _, family := range families.Results {
			filter := map[string]any{
				"id": family.ID,
			}
			update := map[string]any{
				"$set": family,
			}
			opts := options.Update().SetUpsert(true)
			_, err := s.mongoClient.Collection("ll2_launcher_family").UpdateOne(context.Background(), filter, update, opts)
			if err != nil {
				return err
			}
		}
		offset += len(families.Results)
	}
	return nil
}
