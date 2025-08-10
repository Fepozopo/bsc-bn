package po

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractPO parses a PO table and returns a PO struct populated with data.
func ExtractPO(poTable *goquery.Selection) PO {
	po := PO{}
	// Extract PO header info
	// Find the first .tbInfo table for Bill To and Ship To
	poTable.Find("table.tbInfo").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() >= 3 {
			billToHtml, _ := tds.Eq(0).Html()
			shipToHtml, _ := tds.Eq(2).Html()
			po.BillTo = HtmlToMultiline(billToHtml)
			po.ShipTo = HtmlToMultiline(shipToHtml)
		}
	})

	// Find the next .tbborder table for PO Number, Type, Date
	poTable.Find("table.tbborder").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() >= 3 && i == 1 {
			po.Number = strings.TrimSpace(tds.Eq(0).Text())
			po.Date = strings.TrimSpace(tds.Eq(2).Text())
		}
	})

	// Find the .tbborder table that has a header row with "Cancel After"
	poTable.Find("table.tbborder").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, tr *goquery.Selection) {
			if tr.Find("th").Length() > 0 && strings.Contains(tr.Text(), "Cancel After") {
				// The next row should be the detail row
				next := tr.Next()
				tds := next.Find("td")
				if tds.Length() >= 2 {
					po.CancelAfter = strings.TrimSpace(tds.Eq(1).Text())
				}
			}
		})
	})

	// Find the main line items table (the one with class tbborder and many columns)
	var lineTable *goquery.Selection
	poTable.Find("table.tbborder").Each(func(i int, s *goquery.Selection) {
		// Heuristic: the line items table has > 8 th columns
		if s.Find("th").Length() > 8 {
			lineTable = s
		}
	})

	if lineTable == nil {
		return po
	}

	// Extract line items: pair each main row with its detail row
	rows := lineTable.Find("tr")
	for i := 1; i < rows.Length(); i++ {
		row := rows.Eq(i)
		class, _ := row.Attr("class")
		if class == "lineItem" || class == "AltlineItem" {
			// Main row: extract columns
			cols := row.Find("td")
			if cols.Length() < 11 {
				continue
			}
			identifier := extractIdentifier(cols.Eq(2))
			item := LineItem{
				EAN:         identifier["EAN"],
				ISBN:        identifier["ISBN"],
				ArrivalDate: strings.TrimSpace(cols.Eq(10).Text()),
				Quantity:    strings.TrimSpace(cols.Eq(3).Text()),
				ItemCost:    strings.TrimSpace(cols.Eq(6).Text()),
				ItemRetail:  strings.TrimSpace(cols.Eq(7).Text()),
				Discount:    strings.TrimSpace(cols.Eq(8).Text()),
				CasePack:    "",
			}
			// Next row is the detail row
			if i+1 < rows.Length() {
				detailRow := rows.Eq(i + 1)
				detailCols := detailRow.Find("td")
				item.Title = extractDetailField(detailCols, "Title:")
				item.SKU = extractDetailField(detailCols, "Vendor Item Code:")
				item.CasePack = extractDetailField(detailCols, "Case Pack Qty:")
				item.IOQ = extractDetailField(detailCols, "IOQ:")
			}
			po.LineItems = append(po.LineItems, item)
			i++ // skip detail row
		}
	}

	// Find totals table (last .tbborder in this PO)
	poTable.Find("table.tbborder").Each(func(i int, s *goquery.Selection) {
		if s.Find("th").Length() == 2 && s.Find("th").First().Text() == "Total Line Items" {
			tds := s.Find("td")
			if tds.Length() == 2 {
				po.TotalLines = strings.TrimSpace(tds.Eq(0).Text())
				po.TotalQty = strings.TrimSpace(tds.Eq(1).Text())
			}
		}
	})

	return po
}

// extractIdentifier parses EAN/ISBN from a selection.
func extractIdentifier(sel *goquery.Selection) map[string]string {
	result := map[string]string{"EAN": "", "ISBN": ""}
	sel.Find("span").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "EAN:") {
			result["EAN"] = strings.TrimSpace(strings.TrimPrefix(text, "EAN:"))
		}
		if strings.HasPrefix(text, "ISBN:") {
			result["ISBN"] = strings.TrimSpace(strings.TrimPrefix(text, "ISBN:"))
		}
	})
	return result
}

// extractDetailField extracts a detail field from columns based on a label.
func extractDetailField(cols *goquery.Selection, label string) string {
	re := regexp.MustCompile(label + `\s*([^\n\r]*)`)
	for i := 0; i < cols.Length(); i++ {
		m := re.FindStringSubmatch(cols.Eq(i).Text())
		if len(m) == 2 {
			return strings.TrimSpace(m[1])
		}
	}
	return ""
}
