package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
)

func ValidateSignature(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read the body
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}

			// Restore the body so the next handler can read it
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			gotSig := r.Header.Get("X-Vexora-Signature")
			if !verifySignature(secret, bodyBytes, gotSig) {
				log.Println("⛔ Security Alert: Invalid Signature")
				http.Error(w, "Unauthorized", 401)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func verifySignature(secret string, body []byte, gotSig string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(gotSig), []byte(expectedSig))
}
