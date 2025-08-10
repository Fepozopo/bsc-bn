package po

import (
	"fmt"
	"os"
)

// WritePOHTML generates an HTML file for the given PO.
func WritePOHTML(po PO) error {
	filename := fmt.Sprintf("PO_%s.html", po.Number)
	f, err := os.Create(filename)
	if err != nil {
		return err
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
        <div class="address-lines">`+FormatAddressLines(po.BillTo)+`</div>
      </div>
      <div class="address-block">
        <div class="address-label"><b>Ship To:</b></div>
        <div class="address-lines">`+FormatAddressLines(po.ShipTo)+`</div>
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

	return nil
}
