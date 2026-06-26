package crypto

import "testing"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	c, err := New("uma-chave-de-dados-bem-secreta", "uma-chave-de-indice")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	plain := "(11) 99999-8888"
	enc, err := c.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc == plain {
		t.Fatal("esperava ciphertext diferente do texto puro")
	}

	dec, err := c.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec != plain {
		t.Fatalf("round-trip falhou: %q != %q", dec, plain)
	}
}

func TestEncryptIsNonDeterministic(t *testing.T) {
	c, _ := New("k", "i")
	a, _ := c.Encrypt("mesmo-valor")
	b, _ := c.Encrypt("mesmo-valor")
	if a == b {
		t.Fatal("ciphertexts deveriam diferir (nonce aleatório)")
	}
}

func TestDecryptPassesThroughPlaintext(t *testing.T) {
	c, _ := New("k", "i")
	// valor legado (sem prefixo) deve ser devolvido como está
	got, err := c.Decrypt("texto-legado")
	if err != nil {
		t.Fatalf("Decrypt legado: %v", err)
	}
	if got != "texto-legado" {
		t.Fatalf("esperava passthrough, obtive %q", got)
	}
}

func TestDisabledCipherIsPassthrough(t *testing.T) {
	c, _ := New("", "")
	if c.Enabled() {
		t.Fatal("esperava criptografia desabilitada")
	}
	enc, _ := c.Encrypt("oi")
	if enc != "oi" {
		t.Fatalf("esperava passthrough, obtive %q", enc)
	}
}

func TestBlindIndexDeterministicAndIrreversible(t *testing.T) {
	c, _ := New("k", "indice-secreto")
	i1 := c.BlindIndex("11999998888")
	i2 := c.BlindIndex("11999998888")
	if i1 != i2 {
		t.Fatal("índice cego deveria ser determinístico")
	}
	if i1 == "11999998888" || len(i1) != 64 {
		t.Fatalf("índice cego inesperado: %q", i1)
	}
	if c.BlindIndex("11999998888") == c.BlindIndex("11999990000") {
		t.Fatal("valores diferentes deveriam gerar índices diferentes")
	}
}

func TestBlindIndexDependsOnKey(t *testing.T) {
	a, _ := New("k", "chave-A")
	b, _ := New("k", "chave-B")
	if a.BlindIndex("x") == b.BlindIndex("x") {
		t.Fatal("índices com chaves diferentes não deveriam coincidir")
	}
}
