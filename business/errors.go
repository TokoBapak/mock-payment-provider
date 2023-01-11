package business

import "errors"

// ErrDuplicateOrderId should be returned if a duplicate order id was found
var ErrDuplicateOrderId = errors.New("duplicate order id")

// ErrMismatchedTransactionAmount should be returned if the calculation of ProductItems
// differs than the one defined on the TransactionAmount.
var ErrMismatchedTransactionAmount = errors.New("mismatched transaction amount")

// ErrTransactionNotFound should be returned when a transaction was not found
var ErrTransactionNotFound = errors.New("transaction not found")

<<<<<<< HEAD
// RequestValidationError should be returned when a request validation error occured
type RequestValidationError struct {
	Reason string
}

func (r RequestValidationError) Error() string {
	return r.Reason
}
=======
// ErrCannotModifyStatus should be returned for cases that overstep the status flow.
// For example, if the current status is settled, we can't request the status to be changed to
// canceled, and vice versa.
var ErrCannotModifyStatus = errors.New("cannot modify status")
>>>>>>> b7d07bdca84d9abc7f5ece0ab576c1db8aabe777
