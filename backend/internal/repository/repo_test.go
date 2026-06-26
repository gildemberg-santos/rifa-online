package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/user/rifa-online/internal/crypto"
)

var testDB *mongo.Database
var testCipher, _ = crypto.New("test-data-key", "test-index-key")

func TestMain(m *testing.M) {
	uri := os.Getenv("MONGODB_TEST_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to test mongodb: %v", err)
	}

	testDB = client.Database("rifaonline_test")

	code := m.Run()

	testDB.Drop(context.Background())
	client.Disconnect(context.Background())
	os.Exit(code)
}
