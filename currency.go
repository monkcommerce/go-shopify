package goshopify

// CurrencyService is an interface for interfacing with the shop endpoint of the
// Shopify API.
// See: https://help.shopify.com/api/reference/shop
type CurrencyService interface {
	Get(options interface{}) ([]Currency, error)
}

// CurrencyServiceOp handles communication with the shop related methods of the
// Shopify API.
type CurrencyServiceOp struct {
	client *Client
}

// Currency ... represents a currency for the shop
type Currency struct {
	Currency string `json:"currency"`
	Enabled  bool   `json:"enabled"`
}

// CurrencyResource ... Represents the result from the admin/shop.json endpoint
type CurrencyResource struct {
	Currencies []Currency `json:"currencies"`
}

// Get shop
func (s *CurrencyServiceOp) Get(options interface{}) ([]Currency, error) {
	resource := new(CurrencyResource)
	err := s.client.Get("currencies.json", resource, options)
	return resource.Currencies, err
}
