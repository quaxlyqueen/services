package services

import (
	//"context"
	"crypto/aes"
	"crypto/cipher"
	//"sync"

	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	//"encoding/json"
	"errors"
	"fmt"

	"io"
	"log"
	"net/http"

	//"os"
	//"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// TODO: Dynamically get key? Need to research best way to store private keys between two devices.
var key = []byte("passphrasewhichneedstobe32bytes!")

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	jwtToken, err := token.SignedString(key)
	if err != nil {
		return "Error creating JWT.", err
	}
	return jwtToken, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func hash(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func encrypt() {
	var newCommunication Communication

	text := []byte(newCommunication.Communication)

	c, err := aes.NewCipher(key)
	if err != nil {
		return
		//gc.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
		//gc.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
		//gc.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	}

	var b []byte = gcm.Seal(nonce, nonce, text, nil)
	hex, err := convertBytesToHex(b)
	if err != nil {
		fmt.Println(err)
	}
	var test string = "{data: " + hex + ", hash: " + hash(string(text)) + "}"
	log.Println("encrypted string:" + test)
	//gc.IndentedJSON(http.StatusCreated, gin.H{"message": test})
}

func convertBytesToHex(b []byte) (string, error) {
	// Handle nil pointer case
	if b == nil {
		return "", errors.New("nil pointer provided for hex string")
	}

	// Split the hex string by spaces
	var h string = hex.EncodeToString(b)
	var builder strings.Builder
	for i := 0; i < len(h); i += 2 {
		end := i + 2
		if end > len(h) {
			end = len(h)
		}
		chunk := h[i:end]
		builder.WriteString(chunk)
		if end < len(h) {
			builder.WriteString(" ")
		}
	}

	// Return the slice of uint8 and any errors encountered
	return builder.String(), nil
}

func convertHexToBytes(hexString *string) ([]uint8, error) {
	// Handle nil pointer case
	if hexString == nil {
		return nil, errors.New("nil pointer provided for hex string")
	}

	// Split the hex string by spaces
	hexBytes := strings.Fields(*hexString)

	// Initialize an empty slice for uint8
	data := make([]uint8, len(hexBytes))

	// Iterate and convert each hex byte
	for i, hexByte := range hexBytes {
		// Convert each hex string to a uint8 value (handling errors)
		value, err := strconv.ParseUint(hexByte, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("error parsing hex byte '%s': %w", hexByte, err)
		}

		// Assign the converted value to the slice
		data[i] = uint8(value)
	}

	// Return the slice of uint8 and any errors encountered
	return data, nil
}

func decrypt(input *Communication, output *Chat) bool {
	// TODO: Add testing flag for easier manipulation.
	//ciphertext, err := ioutil.ReadFile("myfile")
	ciphertext, err := convertHexToBytes(&input.Communication)

	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		fmt.Println(err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	s := string(plaintext)
	if validateHash(s, input.Hash) {
		output.Content = s
		return true
	}
	input.Communication = "DATA CORRUPTED OR TAMPERED"
	return false
}

// Validate there's no tampering with SHA-1 sum. The decrypted hash and the transmitted hash should be identical.
func validateHash(decrypted string, hash string) bool {
	// Calculate the SHA256 sum of the decrypted request
	hasher := sha1.New()
	hasher.Write([]byte(decrypted))
	decryptedHashString := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// TODO: Add error handling if hash doesn't match.
	return decryptedHashString == hash
}

// Serve the authentication and encryption layer to a provided local port.
// Authentication takes place solely on the backend.
func serveAuthentication(
	router *mux.Router,
	port int,
) {
	// TODO: Add error handling
	authR := router.Host("http://localhost").Subrouter()
	authSrv := &http.Server{
		Addr:         "0.0.0.0:" + string(port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      authR,
	}

	go func() {
		if err := authSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	log.Println("Auth server is running on port " + string(port))
}
