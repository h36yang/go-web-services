package product

// Product model class
type Product struct {
	ProductID      int     `json:"productId"`
	Manufacturer   string  `json:"manufacturer"`
	Sku            string  `json:"sku"`
	Upc            string  `json:"upc"`
	PricePerUnit   float32 `json:"pricePerUnit"`
	QuantityOnHand int     `json:"quantityOnHand"`
	ProductName    string  `json:"productName"`
}
