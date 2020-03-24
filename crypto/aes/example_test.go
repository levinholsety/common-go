package aes_test

import (
	"crypto/rand"
	"io"
	"log"
	"os"

	"github.com/levinholsety/common-go/crypto/aes"
)

func ExampleNewEncryptionWriter() {
	src, err := os.Open("src.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer src.Close()
	dst, err := os.Create("dst.bin")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dst.Close()
	key := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, key); err != nil {
		log.Fatal(err)
		return
	}
	iv := make([]byte, 16)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
		return
	}
	w, err := aes.NewEncryptionWriter(dst, key, iv)
	if err != nil {
		log.Fatal(err)
		return
	}
	if _, err = io.Copy(w, src); err != nil {
		log.Fatal(err)
		return
	}
	if err = w.Close(); err != nil {
		log.Fatal(err)
		return
	}
}

var key, iv []byte

func ExampleNewDecryptionReader() {
	src, err := os.Open("src.bin")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer src.Close()
	dst, err := os.Open("dst.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dst.Close()
	r, err := aes.NewDecryptionReader(src, key, iv)
	if err != nil {
		log.Fatal(err)
		return
	}
	if _, err = io.Copy(dst, r); err != nil {
		log.Fatal(err)
		return
	}
}
