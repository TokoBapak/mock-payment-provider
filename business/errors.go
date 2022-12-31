package business

import "errors"

// ErrDuplicateOrderId should be returned if a duplicate order id was found
var ErrDuplicateOrderId = errors.New("duplicate order id")

// ErrMismatchedTransactionAmount should be returned if the calculation of ProductItems
// differs than the one defined on the TransactionAmount.
var ErrMismatchedTransactionAmount = errors.New("mismatched transaction amount")

// ErrTransactionNotFound should be returned when a transaction was not found
var ErrTransactionNotFound = errors.New("transaction not found")
