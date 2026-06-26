// Comando de migração: criptografa, em repouso, os dados pessoais de
// participantes já existentes (tickets e payments) e popula o índice cego do
// telefone. É idempotente — registros já criptografados são ignorados.
//
// Uso:
//
//	cd backend && go run ./cmd/migrate-encrypt
//
// Requer DATA_ENCRYPTION_KEY (e BLIND_INDEX_KEY) definidos no ambiente, com os
// MESMOS valores usados pelo servidor.
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

	n1, err := migrateCollection(ctx, db.Collection("tickets"), cipher, false)
	if err != nil {
		log.Fatalf("erro migrando tickets: %v", err)
	}
	log.Printf("tickets migrados: %d", n1)

	n2, err := migrateCollection(ctx, db.Collection("payments"), cipher, true)
	if err != nil {
		log.Fatalf("erro migrando payments: %v", err)
	}
	log.Printf("payments migrados: %d", n2)

	log.Println("migração concluída")
}

// migrateCollection criptografa buyerName/buyerPhone (e buyerEmail, quando
// withEmail) dos documentos ainda em texto puro e grava o índice cego do telefone.
func migrateCollection(ctx context.Context, coll *mongo.Collection, cipher *crypto.Cipher, withEmail bool) (int, error) {
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	migrated := 0
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return migrated, err
		}

		name, _ := doc["buyerName"].(string)
		phone, _ := doc["buyerPhone"].(string)
		email, _ := doc["buyerEmail"].(string)

		// Já criptografado em qualquer campo => registro já migrado.
		if crypto.IsEncrypted(name) || crypto.IsEncrypted(phone) || crypto.IsEncrypted(email) {
			continue
		}

		set := bson.M{}
		if phone != "" {
			set["buyerPhoneIndex"] = cipher.BlindIndex(phone)
			enc, err := cipher.Encrypt(phone)
			if err != nil {
				return migrated, err
			}
			set["buyerPhone"] = enc
		}
		if name != "" {
			enc, err := cipher.Encrypt(name)
			if err != nil {
				return migrated, err
			}
			set["buyerName"] = enc
		}
		if withEmail && email != "" {
			enc, err := cipher.Encrypt(email)
			if err != nil {
				return migrated, err
			}
			set["buyerEmail"] = enc
		}
		if len(set) == 0 {
			continue
		}

		if _, err := coll.UpdateByID(ctx, doc["_id"], bson.M{"$set": set}); err != nil {
			return migrated, err
		}
		migrated++
	}
	return migrated, cursor.Err()
}
