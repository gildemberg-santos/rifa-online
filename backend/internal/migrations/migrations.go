package migrations

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	Version     int
	Description string
	Up          func(ctx context.Context, db *mongo.Database) error
}

var registry []Migration

func register(m Migration) {
	registry = append(registry, m)
}

func Run(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection("_migrations")

	if err := ensureCollection(ctx, coll); err != nil {
		return fmt.Errorf("migrations: ensure collection: %w", err)
	}

	applied := getApplied(ctx, coll)
	appliedMap := make(map[int]bool)
	for _, v := range applied {
		appliedMap[v] = true
	}

	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Version < registry[j].Version
	})

	for _, m := range registry {
		if appliedMap[m.Version] {
			slog.Debug("migrations: already applied", "version", m.Version, "description", m.Description)
			continue
		}

		slog.Info("migrations: applying", "version", m.Version, "description", m.Description)

		if err := m.Up(ctx, db); err != nil {
			return fmt.Errorf("migrations: version %d (%s): %w", m.Version, m.Description, err)
		}

		if _, err := coll.InsertOne(ctx, bson.M{
			"version":     m.Version,
			"description": m.Description,
			"appliedAt":   time.Now(),
		}); err != nil {
			return fmt.Errorf("migrations: failed to record version %d: %w", m.Version, err)
		}

		slog.Info("migrations: applied", "version", m.Version, "description", m.Description)
	}

	slog.Info("migrations: all up to date", "total", len(registry))
	return nil
}

func ensureCollection(ctx context.Context, coll *mongo.Collection) error {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "version", Value: 1}},
	})
	return err
}

func getApplied(ctx context.Context, coll *mongo.Collection) []int {
	cursor, err := coll.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil
	}
	var results []struct {
		Version int `bson:"version"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil
	}
	out := make([]int, len(results))
	for i, r := range results {
		out[i] = r.Version
	}
	return out
}
