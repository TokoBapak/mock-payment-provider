package virtual_account

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func generateVirtualAccountNumber() string {
	// Virtual account is a set number with length of 15.
	// Why 15? Because real virtual account number ranges from 11-12 in length,
	// we want to avoid people paying to the actual bank.
	// To make the matter unique, we will use some combination of current date

	var builder strings.Builder

	currentTime := time.Now()
	builder.WriteString(currentTime.Format("06"))   // year
	builder.WriteString(strconv.Itoa(rand.Intn(9))) // random single digit
	builder.WriteString(currentTime.Format("02"))   // day
	builder.WriteString(strconv.Itoa(rand.Intn(9))) // random single digit
	builder.WriteString(currentTime.Format("01"))   // month
	builder.WriteString(strconv.Itoa(rand.Intn(9))) // random single digit
	builder.WriteString(currentTime.Format("15"))   // minute
	builder.WriteString(strconv.Itoa(rand.Intn(9))) // random single digit
	builder.WriteString(currentTime.Format("04"))   // second

	return builder.String()
}
