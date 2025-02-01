package goshopify

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const recurringApplicationChargesBasePath = "recurring_application_charges"

// RecurringApplicationChargeService is an interface for interacting with the
// RecurringApplicationCharge endpoints of the Shopify API.
// See https://help.shopify.com/api/reference/billing/recurringapplicationcharge
type RecurringApplicationChargeService interface {
	Create(RecurringApplicationCharge) (*RecurringApplicationCharge, error)
	Get(int64, interface{}) (*RecurringApplicationCharge, error)
	List(interface{}) ([]RecurringApplicationCharge, error)
	Activate(RecurringApplicationCharge) (*RecurringApplicationCharge, error)
	Delete(int64) error
	Update(int64, int64) (*RecurringApplicationCharge, error)
}

// RecurringApplicationChargeServiceOp handles communication with the
// RecurringApplicationCharge related methods of the Shopify API.
type RecurringApplicationChargeServiceOp struct {
	client *Client
}

// RecurringApplicationCharge represents a Shopify RecurringApplicationCharge.
type RecurringApplicationCharge struct {
	ConfirmationURL  string                   `json:"confirmation_url"`
	ID               int64                    `json:"id"`
	Name             string                   `json:"name"`
	ReturnURL        string                   `json:"return_url"`
	CurrencyCode     string                   `json:"currencyCode"`
	Price            *decimal.Decimal         `json:"price"`
	Interval         string                   `json:"interval"`
	RecurringPricing *AppPlanRecurringPricing `json:"recurringPricing"`
	CappedAmount     *decimal.Decimal         `json:"capped_amount"`
	Terms            string                   `json:"terms"`
	UsagePricing     *AppPlanUsagePricing     `json:"usagePricing"`
	Status           string                   `json:"status"`
	Test             *bool                    `json:"test"`
	UpdatedAt        *time.Time               `json:"updated_at"`
}

// AppPlanRecurringPricing ... The pricing model input can be either appRecurringPricingDetails or appUsagePricingDetails.
type AppPlanRecurringPricing struct {
	AppRecurringPricingDetails struct {
		Interval string     `json:"interval"`
		Price    MoneyInput `json:"price"`
	} `json:"appRecurringPricingDetails"`
}

// AppPlanUsagePricing ... The pricing model input can be either appRecurringPricingDetails or appUsagePricingDetails.
type AppPlanUsagePricing struct {
	AppUsagePricingDetails struct {
		CappedAmount MoneyInput `json:"cappedAmount"`
		Terms        string     `json:"terms"`
	} `json:"appUsagePricingDetails"`
}

// MoneyInput ...
type MoneyInput struct {
	Amount       *decimal.Decimal `json:"amount"`
	CurrencyCode string           `json:"currencyCode"`
}

func parse(dest **time.Time, data *string) error {
	if data == nil {
		return nil
	}
	// This is what API doc says: "2013-06-27T08:48:27-04:00"
	format := time.RFC3339
	if len(*data) == 10 {
		// This is how the date looks.
		format = "2006-01-02"
	}
	t, err := time.Parse(format, *data)
	if err != nil {
		return err
	}
	*dest = &t
	return nil
}

// UnmarshalJSON ...
func (r *RecurringApplicationCharge) UnmarshalJSON(data []byte) error {
	// This is a workaround for the API returning incomplete results:
	// https://ecommerce.shopify.com/c/shopify-apis-and-technology/t/-523203
	// For a longer explanation of the hack check:
	// http://choly.ca/post/go-json-marshalling/
	type alias RecurringApplicationCharge
	aux := &struct {
		ActivatedOn *string `json:"activated_on"`
		BillingOn   *string `json:"billing_on"`
		CancelledOn *string `json:"cancelled_on"`
		CreatedAt   *string `json:"created_at"`
		TrialEndsOn *string `json:"trial_ends_on"`
		UpdatedAt   *string `json:"updated_at"`
		*alias
	}{alias: (*alias)(r)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if err := parse(&r.UpdatedAt, aux.UpdatedAt); err != nil {
		return err
	}
	return nil
}

// RecurringApplicationChargeResource represents the result from the
// admin/recurring_application_charges{/X{/activate.json}.json}.json endpoints.
type RecurringApplicationChargeResource struct {
	Charge *RecurringApplicationCharge `json:"recurring_application_charge"`
}

// RecurringApplicationChargesResource represents the result from the
// admin/recurring_application_charges.json endpoint.
type RecurringApplicationChargesResource struct {
	Charges []RecurringApplicationCharge `json:"recurring_application_charges"`
}

// Create creates new recurring application charge.
func (r *RecurringApplicationChargeServiceOp) Create(charge RecurringApplicationCharge) (
	*RecurringApplicationCharge, error) {

	path := fmt.Sprintf("%s.json", recurringApplicationChargesBasePath)
	wrappedData := RecurringApplicationChargeResource{Charge: &charge}
	resource := &RecurringApplicationChargeResource{}
	err := r.client.Post(path, wrappedData, resource)
	return resource.Charge, err
}

// Get gets individual recurring application charge.
func (r *RecurringApplicationChargeServiceOp) Get(chargeID int64, options interface{}) (
	*RecurringApplicationCharge, error) {

	path := fmt.Sprintf("%s/%d.json", recurringApplicationChargesBasePath, chargeID)
	resource := &RecurringApplicationChargeResource{}
	err := r.client.Get(path, resource, options)
	return resource.Charge, err
}

// List gets all recurring application charges.
func (r *RecurringApplicationChargeServiceOp) List(options interface{}) (
	[]RecurringApplicationCharge, error) {

	path := fmt.Sprintf("%s.json", recurringApplicationChargesBasePath)
	resource := &RecurringApplicationChargesResource{}
	err := r.client.Get(path, resource, options)
	return resource.Charges, err
}

// Activate activates recurring application charge.
func (r *RecurringApplicationChargeServiceOp) Activate(charge RecurringApplicationCharge) (
	*RecurringApplicationCharge, error) {

	path := fmt.Sprintf("%s/%d/activate.json", recurringApplicationChargesBasePath, charge.ID)
	wrappedData := RecurringApplicationChargeResource{Charge: &charge}
	resource := &RecurringApplicationChargeResource{}
	err := r.client.Post(path, wrappedData, resource)
	return resource.Charge, err
}

// Delete deletes recurring application charge.
func (r *RecurringApplicationChargeServiceOp) Delete(chargeID int64) error {
	return r.client.Delete(fmt.Sprintf("%s/%d.json", recurringApplicationChargesBasePath, chargeID))
}

// Update updates recurring application charge.
func (r *RecurringApplicationChargeServiceOp) Update(chargeID, newCappedAmount int64) (
	*RecurringApplicationCharge, error) {

	path := fmt.Sprintf("%s/%d/customize.json?recurring_application_charge[capped_amount]=%d",
		recurringApplicationChargesBasePath, chargeID, newCappedAmount)
	resource := &RecurringApplicationChargeResource{}
	err := r.client.Put(path, nil, resource)
	return resource.Charge, err
}
