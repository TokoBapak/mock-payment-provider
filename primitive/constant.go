package primitive

import "regexp"

var EmailPattern = regexp.MustCompile(`^\w+@\w+(\.\w+)+$`)
var PhoneNumberPattern = regexp.MustCompile(`^\+\d+$`)

const PostalCodePattern = `^\d{5,10}$`
const CountryCodePattern = `^\d{0,5}$`
