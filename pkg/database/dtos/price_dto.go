package dtos

type Price struct {
	ID           uint   `json:"id"`
	Fractional   int    `json:"fractional"`
	Decimal      string `json:"decimal"`
	Tax          int    `json:"tax"`
	TaxDecimal   string `json:"tax_decimal"`
	DisplayValue string `json:"display_value"`
	FullDisplay  string `json:"full_display"`
	Currency     string `json:"currency"`
}
