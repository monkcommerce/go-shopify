package goshopify

// BcModifier ...
type BcModifier struct {
	ID           int    `json:"id"`
	ProductID    int    `json:"product_id"`
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	SortOrder    int    `json:"sort_order"`
	OptionValues []struct {
		ID        int         `json:"id"`
		OptionID  int         `json:"option_id"`
		Label     string      `json:"label"`
		SortOrder int         `json:"sort_order"`
		ValueData interface{} `json:"value_data"`
		IsDefault bool        `json:"is_default"`
		Adjusters struct {
			Price struct {
				Adjuster      string  `json:"adjuster"`
				AdjusterValue float64 `json:"adjuster_value"`
			} `json:"price"`
			Weight             interface{} `json:"weight"`
			ImageURL           string      `json:"image_url"`
			PurchasingDisabled struct {
				Status  bool   `json:"status"`
				Message string `json:"message"`
			} `json:"purchasing_disabled"`
		} `json:"adjusters"`
	} `json:"option_values"`
}

// BcOptionValues ...
type BcOptionValues struct {
	ID                int    `json:"id"`
	Label             string `json:"label"`
	OptionID          int    `json:"option_id"`
	OptionDisplayName string `json:"option_display_name"`
}
