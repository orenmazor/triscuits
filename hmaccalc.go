package main

import "crypto/hmac"
import "crypto/sha256"
import "os"
import "encoding/base64"

func DecodeHMAC(message string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil
	}
	return decoded
}

func CalculateHMAC(message string) string {
	key := os.Getenv("TRISCUITS_HMAC")
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(expectedMAC)
}

func CheckHMAC(message string, messageMAC string) bool {
	extracted_mac := DecodeHMAC(messageMAC)
	calculated_hmac := DecodeHMAC(CalculateHMAC(message))
	return hmac.Equal(calculated_hmac, extracted_mac)
}
