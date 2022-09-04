package query

import (
	"bytes"
	"testing"
)

const (
	Key       = "thisis32bitlongpassphraseimusing"
	SliceSize = 4
)

var (
	Texts = []string{
		// "This is a secret",
		// "abcz",
		// "666dd",
		// "14512512",
		// "myzzz",
		// "dasdasz",
		// "zcxzcxdasas",
		// "zcxczxczx2131",
		"012345678901",
	}
)

func TestNormal(t *testing.T) {
	c := New([]byte(Key), SliceSize)
	for _, text := range Texts {
		enc, err := c.Encrypt([]byte(text))
		if err != nil {
			t.Errorf("err %v", err)
		}
		dec, err := c.Decrypt(enc)
		if err != nil {
			t.Errorf("err %v", err)
		}
		if !bytes.Equal(dec, []byte(text)) {
			t.Errorf("wand %s, but %s", text, dec)
		}
	}
}

func TestBase64(t *testing.T) {
	c := New([]byte(Key), SliceSize)
	for _, text := range Texts {
		enc, err := c.EncryptToBase64(text)
		if err != nil {
			t.Errorf("err %v", err)
		}
		dec, err := c.DecryptFromBase64(enc)
		if err != nil {
			t.Errorf("err %v", err)
		}
		if dec != text {
			t.Errorf("wand %s, but %s", text, dec)
		}
	}
}
