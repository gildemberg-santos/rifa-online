// Package crypto fornece criptografia de campos sensíveis em repouso
// (application-level field encryption) e um índice cego (blind index) para
// permitir buscas por igualdade sobre campos criptografados.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

// encPrefix identifica valores criptografados, permitindo conviver com dados
// legados em texto puro durante a migração (Decrypt devolve-os como estão).
const encPrefix = "enc:v1:"

// Cipher criptografa/descriptografa valores (AES-256-GCM) e gera índices cegos
// (HMAC-SHA256). Se a chave de dados estiver vazia, a criptografia fica
// desabilitada e os valores trafegam em texto puro (útil em desenvolvimento).
type Cipher struct {
	aead     cipher.AEAD
	indexKey []byte
	enabled  bool
}

// New deriva chaves de 32 bytes (SHA-256) a partir dos segredos informados.
// dataKey vazio => criptografia desabilitada.
func New(dataKey, indexKey string) (*Cipher, error) {
	if dataKey == "" {
		return &Cipher{enabled: false}, nil
	}
	dk := sha256.Sum256([]byte(dataKey))
	block, err := aes.NewCipher(dk[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ik := sha256.Sum256([]byte(indexKey))
	return &Cipher{aead: aead, indexKey: ik[:], enabled: true}, nil
}

// Enabled informa se a criptografia está ativa.
func (c *Cipher) Enabled() bool { return c.enabled }

// IsEncrypted informa se o valor já está no formato criptografado.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encPrefix)
}

// Encrypt devolve um valor autodescritivo (prefixo + nonce + ciphertext em base64).
// Strings vazias e o modo desabilitado retornam o texto original.
func (c *Cipher) Encrypt(plaintext string) (string, error) {
	if !c.enabled || plaintext == "" {
		return plaintext, nil
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := c.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return encPrefix + base64.StdEncoding.EncodeToString(ct), nil
}

// Decrypt reverte Encrypt. Valores sem o prefixo são considerados texto puro
// (legado) e devolvidos como estão, viabilizando migração incremental.
func (c *Cipher) Decrypt(value string) (string, error) {
	if !strings.HasPrefix(value, encPrefix) {
		return value, nil
	}
	if !c.enabled {
		return "", errors.New("valor criptografado, mas criptografia desabilitada (chave ausente)")
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, encPrefix))
	if err != nil {
		return "", err
	}
	ns := c.aead.NonceSize()
	if len(raw) < ns {
		return "", errors.New("ciphertext inválido")
	}
	nonce, ct := raw[:ns], raw[ns:]
	pt, err := c.aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

// BlindIndex gera um identificador determinístico (HMAC-SHA256 em hex) para
// buscas por igualdade sobre campos criptografados (ex.: telefone). O mesmo
// valor sempre produz o mesmo índice, mas o índice é irreversível.
func (c *Cipher) BlindIndex(value string) string {
	if value == "" {
		return ""
	}
	mac := hmac.New(sha256.New, c.indexKey)
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}
