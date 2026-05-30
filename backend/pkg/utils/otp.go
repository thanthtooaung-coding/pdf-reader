package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"
)

const otpDigits = 6

func GenerateOTP() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func HashOTP(otp, pepper string) string {
	mac := hmac.New(sha256.New, []byte(pepper))
	_, _ = mac.Write([]byte(otp))
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyOTP(otp, pepper, hash string) bool {
	expected := HashOTP(otp, pepper)
	return subtle.ConstantTimeCompare([]byte(expected), []byte(hash)) == 1
}
