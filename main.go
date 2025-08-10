package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LineItem struct {
	EAN, ISBN, Title, SKU, ArrivalDate, Quantity, ItemCost, ItemRetail, Discount, CasePack, IOQ string
}

type PO struct {
	Number, Date, ShipTo, BillTo, Terms, CancelAfter, BackOrder, SpecialInfo string
	LineItems                                                                []LineItem
	TotalLines, TotalQty, TotalExtCost, TotalExtRetail                       string
}

func main() {
	f, err := os.Open("POFile.HTM")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		id, exists := s.Attr("id")
		if exists && strings.HasPrefix(id, "PO_") {
			po := extractPO(s)
			writePOHTML(po)
		}
	})
}

func extractPO(poTable *goquery.Selection) PO {
	po := PO{}
	// Extract PO header info
	// Find the first .tbInfo table for Bill To and Ship To
	poTable.Find("table.tbInfo").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() >= 3 {
			billToHtml, _ := tds.Eq(0).Html()
			shipToHtml, _ := tds.Eq(2).Html()
			po.BillTo = htmlToMultiline(billToHtml)
			po.ShipTo = htmlToMultiline(shipToHtml)
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

// Helper to extract detail fields
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

// Helper to format address lines for Bill To / Ship To
func htmlToMultiline(html string) string {
	re := regexp.MustCompile(`(?i)<br\s*/?>`)
	html = re.ReplaceAllString(html, "\n")
	// Remove any remaining HTML tags (if any)
	reTag := regexp.MustCompile(`<[^>]+>`)
	html = reTag.ReplaceAllString(html, "")
	return html
}

func formatAddressLines(addr string) string {
	lines := []string{}
	for _, line := range strings.Split(addr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Remove leading "Bill To:" or "Ship To:"
		line = strings.TrimPrefix(line, "Bill To:")
		line = strings.TrimPrefix(line, "Ship To:")
		line = strings.TrimSpace(line)
		lines = append(lines, line)
	}
	return strings.Join(lines, "<br>")
}

func writePOHTML(po PO) {
	filename := fmt.Sprintf("PO_%s.html", po.Number)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprint(f, `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Purchase Order `+po.Number+`</title>
<style>
body { font-family: Arial, sans-serif; background: #fff; }
.outer-box { border: 3px solid #000; width: 950px; margin: 24px auto; background: #fff; padding: 0; }
.header-table { width: 100%; border-collapse: collapse; }
.header-table td { vertical-align: top; }
.logo-title { text-align: center; font-size: 2em; font-weight: bold; letter-spacing: 2px; padding-top: 12px; }
.address-block { font-size: 1em; text-align: left; padding-left: 24px; padding-bottom: 8px; }
.po-info-box { border: 2px solid #000; padding: 10px 18px; font-size: 1em; margin-top: 12px; margin-right: 24px; float: right; min-width: 260px; }
.po-info-box table { width: 100%; }
.po-info-box td { padding: 2px 6px; }
.section-row { display: flex; justify-content: space-between; margin: 0 24px 0 24px; }
.section-box { border: 2px solid #000; padding: 10px 18px; width: 46%; margin-bottom: 0; }
.section-title { font-weight: bold; margin-bottom: 4px; }
.line-items-table { width: 98%; border-collapse: collapse; margin: 18px auto 0 auto; }
.line-items-table th, .line-items-table td { border: 2px solid #000; padding: 6px 8px; }
.line-items-table th { background: #eee; font-size: 1.05em; }
.item-info-header th { font-size: 1.1em; background: #fff; border: none; border-bottom: 2px solid #000; text-align: left; padding-top: 18px; padding-bottom: 6px; }
.casepack-row td { border: none; font-size: 0.97em; padding-top: 0; padding-bottom: 0; padding-left: 24px; color: #222; }
.line-items-table td.num { text-align: right; font-family: monospace; }
.totals-table { width: 98%; border-collapse: collapse; margin: 18px auto 0 auto; }
.totals-table th, .totals-table td { border: 2px solid #000; padding: 6px 8px; font-weight: bold; background: #eee; }
.footer-box { border: 2px solid #000; margin: 18px 24px 24px 24px; padding: 8px 16px; font-size: 1em; }
</style>
</head>
<body>
<div class="outer-box">
  <table class="header-table">
    <tr>
      <td style="width:60%;">
        <div class="logo-title">BARNES &amp; NOBLE BOOKSELLERS</div>
      </td>
      <td style="width:40%; vertical-align:top;">
        <div class="po-info-box">
           <div style="font-size:1.2em; font-weight:bold; text-align:center;">PURCHASE ORDER</div>
           <table>
             <tr><td><b>NUMBER</b></td><td>`+po.Number+`</td></tr>
             <tr><td><b>DATE</b></td><td>`+po.Date+`</td></tr>
             <tr><td><b>CANCEL AFTER</b></td><td>`+po.CancelAfter+`</td></tr>
           </table>
         </div>
      </td>
    </tr>
  </table>

  <style>
    .address-outer-box {
      border: 3px solid #000;
      margin: 24px 0 0 0;
      background: #fafafa;
      padding: 30px 10px 30px 10px;
    }
    .address-section {
      display: flex;
      justify-content: space-between;
    }
    .address-block {
      width: 48%;
      text-align: center;
    }
    .address-label {
      font-size: 1.1em;
      font-weight: bold;
      margin-bottom: 10px;
      text-align: center;
    }
    .address-lines {
      font-size: 1em;
      line-height: 1.3em;
      text-align: center;
      margin-top: 8px;
    }
  </style>
  <div class="address-outer-box">
    <div class="address-section">
      <div class="address-block">
        <div class="address-label"><b>Bill To:</b></div>
        <div class="address-lines">`+formatAddressLines(po.BillTo)+`</div>
      </div>
      <div class="address-block">
        <div class="address-label"><b>Ship To:</b></div>
        <div class="address-lines">`+formatAddressLines(po.ShipTo)+`</div>
      </div>
    </div>
  </div>

  <table class="line-items-table">
    <tr class="item-info-header">
      <th colspan="8">ITEM INFORMATION</th>
    </tr>
    <tr>
      <th>EAN</th><th>Title</th><th>SKU</th>
      <th>EXPECTED ARRIVAL DATE</th><th>QUANTITY</th>
      <th>ITEM COST</th><th>ITEM RETAIL</th>
      <th>DISC%</th>
    </tr>
`)

	var totalCost float64
	for _, item := range po.LineItems {
		// Calculate total cost: Quantity * ItemCost
		var qty int
		var cost float64
		fmt.Sscanf(item.Quantity, "%d", &qty)
		fmt.Sscanf(item.ItemCost, "%f", &cost)
		totalCost += float64(qty) * cost

		fmt.Fprintf(f, `<tr>
<td>%s</td><td>%s</td><td>%s</td>
<td>%s</td>
<td class="num">%s</td>
<td class="num">%s</td>
<td class="num">%s</td>
<td class="num">%s</td>
</tr>
<tr class="casepack-row"><td colspan="8"><b>CASEPACK QTY:</b> %s &nbsp;&nbsp; <b>IOQ:</b> %s</td></tr>
`, item.EAN, item.Title, item.SKU, item.ArrivalDate, item.Quantity, item.ItemCost, item.ItemRetail, item.Discount, item.CasePack, item.IOQ)
	}

	fmt.Fprintf(f, `
  </table>

  <table class="totals-table">
    <tr>
      <td style="width:20%%;">TOTAL LINES</td>
      <td style="width:20%%;">%s</td>
      <td style="width:20%%;">TOTAL QTY</td>
      <td style="width:20%%;">%s</td>
      <td style="width:20%%;">TOTAL COST</td>
      <td style="width:20%%;">$%.2f</td>
    </tr>
  </table>

  <div class="footer-box">
  </div>
</div>
</body>
</html>
`, po.TotalLines, po.TotalQty, totalCost)
}
