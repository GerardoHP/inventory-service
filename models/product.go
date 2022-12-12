package models

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

func (p Product) ToMap() (obj map[string]interface{}) {
	obj = make(map[string]interface{})
	obj["productId"] = p.ProductID
	obj["manufacturer"] = p.Manufacturer
	obj["sku"] = p.Sku
	obj["upc"] = p.Upc
	obj["pricePerUnit"] = p.PricePerUnit
	obj["quantityOnHand"] = p.QuantityOnHand
	obj["productName"] = p.ProductName

	return
}

func ProductFromMap(obj map[string]interface{}) *Product {
	product := Product{}
	product.ProductID = mapFromInt("productId", obj)
	product.Manufacturer = mapFromString("manufacturer", obj)
	product.Sku = mapFromString("sku", obj)
	product.Upc = mapFromString("upc", obj)
	product.PricePerUnit = mapFromString("pricePerUnit", obj)
	product.QuantityOnHand = mapFromInt("quantityOnHand", obj)
	product.ProductName = mapFromString("productName", obj)

	return &product
}

func mapFromString(key string, obj map[string]interface{}) (value string) {
	if str, ok := obj[key].(string); ok {
		value = str
	} else {
		value = ""
	}

	return
}

func mapFromInt(key string, obj map[string]interface{}) (value int) {
	if number, ok := obj[key].(int64); ok {
		value = int(number)
	} else {
		value = -1
	}

	return
}
