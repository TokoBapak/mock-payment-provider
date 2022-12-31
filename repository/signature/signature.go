package signature

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func Generate(orderId string, statusCode int, transactionAmount int64, serverKey string) string {
	hash := sha512.New()
	hash.Write([]byte(fmt.Sprintf("%s%d%d%s", orderId, statusCode, transactionAmount, serverKey)))

	return hex.EncodeToString(hash.Sum(nil))
}
