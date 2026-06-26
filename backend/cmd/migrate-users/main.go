// Comando de migração: criptografa name e phone dos usuários existentes.
// Idempotente — registros já criptografados são ignorados.
//
// Uso:
//
//	cd backend && go run ./cmd/migrate-users
//
// Requer DATA_ENCRYPTION_KEY e BLIND_INDEX_KEY no ambiente.
package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/crypto"
)

func main() {
	cfg := config.Load()

	cipher, err := crypto.New(cfg.DataEncryptionKey, cfg.BlindIndexKey)
	if err != nil {
		log.Fatalf("falha ao iniciar cipher: %v", err)
	}
	if !cipher.Enabled() {
		log.Fatal("DATA_ENCRYPTION_KEY ausente: defina a chave antes de migrar")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("falha ao conectar no mongo: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.MongoDBName)
	coll := db.Collection("users")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("erro buscando usuários: %v", err)
	}
	defer cursor.Close(ctx)

	migrated := 0
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Fatalf("erro decodificando: %v", err)
		}

		name, _ := doc["name"].(string)
		phone, _ := doc["phone"].(string)

		if crypto.IsEncrypted(name) || crypto.IsEncrypted(phone) {
			continue
		}

		set := bson.M{}
		if name != "" {
			enc, err := cipher.Encrypt(name)
			if err != nil {
				log.Fatalf("erro criptografando name: %v", err)
			}
			set["name"] = enc
		}
		if phone != "" {
			enc, err := cipher.Encrypt(phone)
			if err != nil {
				log.Fatalf("erro criptografando phone: %v", err)
			}
			set["phone"] = enc
		}
		if len(set) == 0 {
			continue
		}

		if _, err := coll.UpdateByID(ctx, doc["_id"], bson.M{"$set": set}); err != nil {
			log.Fatalf("erro atualizando usuário: %v", err)
		}
		migrated++
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("erro no cursor: %v", err)
	}

	log.Printf("usuários migrados: %d", migrated)
	log.Println("migração de usuários concluída")
}
