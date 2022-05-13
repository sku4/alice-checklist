package channel

import "github.com/sku4/alice-checklist/models/alice"

type Response struct {
	alice.Response
	Error error
}
