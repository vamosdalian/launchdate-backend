package service

import "github.com/vamosdalian/launchdate-backend/internal/db"

type LL2Service struct {
	mongoClient *db.MongoDB
}

func NewLL2Service(db *db.MongoDB) *LL2Service {
	return &LL2Service{
		mongoClient: db,
	}
}

func (s *LL2Service) Ping() error {
	return s.mongoClient.Client.Ping(nil, nil)
}
