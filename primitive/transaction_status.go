package primitive

type TransactionStatus uint8

const (
	// TransactionStatusUnspecified sets the zero value. If this is ever read, it means
	// something is wrong with the code.
	TransactionStatusUnspecified TransactionStatus = iota
	// TransactionStatusPending states that transaction is in progress, and it's pending for verification
	// and/or settlement from the bank.
	TransactionStatusPending
	// TransactionStatusDenied tells that the transaction has been denied by the payment provider.
	TransactionStatusDenied
	// TransactionStatusSettled tells that the transaction has been settled and successful. No further steps are needed.
	TransactionStatusSettled
	// TransactionStatusExpired tells that the transaction has exceeds the time limit that the user is allowed to pay.
	TransactionStatusExpired
	// TransactionStatusCanceled indicates that the transaction is canceled.
	TransactionStatusCanceled
)

func (t TransactionStatus) String() string {
	switch t {
	case TransactionStatusPending:
		return "pending"
	case TransactionStatusDenied:
		return "denied"
	case TransactionStatusSettled:
		return "settled"
	case TransactionStatusExpired:
		return "expired"
	case TransactionStatusCanceled:
		return "canceled"
	case TransactionStatusUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}
