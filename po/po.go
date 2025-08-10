package po

// LineItem represents a single line item in a purchase order.
type LineItem struct {
	EAN         string
	ISBN        string
	Title       string
	SKU         string
	ArrivalDate string
	Quantity    string
	ItemCost    string
	ItemRetail  string
	Discount    string
	CasePack    string
	IOQ         string
}

// PO represents a purchase order with its header and line items.
type PO struct {
	Number         string
	Date           string
	ShipTo         string
	BillTo         string
	Terms          string
	CancelAfter    string
	BackOrder      string
	SpecialInfo    string
	LineItems      []LineItem
	TotalLines     string
	TotalQty       string
	TotalExtCost   string
	TotalExtRetail string
}
