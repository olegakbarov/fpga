package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/julienschmidt/httprouter"
)

var SECRET string = os.Getenv("SECRET")

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func GenerateToken() ([]byte, error) {
	t := time.Now().Format("2006-01-02 15:04:05")

	return encrypt([]byte(t), []byte(SECRET))
}

func GetToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// user id
	// password candidate
	var rec RawUser

	err := db.QueryRow("SELECT * FROM users WHERE id=$1 ORDER BY id", id).Scan(

	if err != nil {
		log.Fatal("Error quering the db- " + err.Error())
		w.WriteHeader(500)
		return
	}

	// hash, err := HashPassword(password)
	// match := CheckPasswordHash(password, hash)
	token, err := GenerateToken()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(401)
	}

	res := Envelope{
		Result: "OK",
		Data:   token,
	}

	data, err := json.Marshal(res)

	if err != nil {
		log.Fatal("Failed marshaling json" + err.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

func CheckAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var authHeader string
	if _, ok := r.Header["Authentication"]; ok {
		authHeader = r.Header.Get("Authentication")
	}
}
