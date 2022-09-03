package aes

import (
	"bytes"
	"testing"
)

const (
	Key = "thisis32bitlongpassphraseimusing"
)

var (
	Texts = []string{
		"This is a secret",
		"abc",
		"666",
		"1",
		"myzzz",
		"",
		"zcxzcxdasas",
		"zcxczxczx2131",
	}
)

func TestNormal(t *testing.T) {
	for _, text := range Texts {
		c := New([]byte(Key))
		enc := c.Encrypt([]byte(text))
		dec := c.Decrypt(enc)
		if !bytes.Equal(dec, []byte(text)) {
			t.Errorf("wand %s, but %s", text, dec)
		}
	}
}

func TestBase64(t *testing.T) {
	for _, text := range Texts {
		c := New([]byte(Key))
		enc := c.EncryptToBase64(text)
		dec, err := c.DecryptFromBase64(enc)
		if err != nil {
			t.Errorf("err %v", err)
		}

		if dec != text {
			t.Errorf("wand %s, but %s", text, dec)
		}
	}
}
