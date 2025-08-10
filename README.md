# bsc-bn

**bsc-bn** is a command-line tool written in Go for parsing Barnes & Noble Purchase Order (PO) HTML files. It extracts structured purchase order data and generates clean, readable HTML summary files for each PO.

## Features

- Parses Barnes & Noble PO HTML files (as downloaded from their vendor portal).
- Extracts key information: PO number, date, cancel after, bill to/ship to addresses, line items, totals, and more.
- Outputs a well-formatted HTML summary for each PO, suitable for printing or archiving.
- Processes multiple POs in a single HTML file if present.

## Requirements

- Go 1.18 or newer (tested with Go 1.24.5)

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/Fepozopo/bsc-bn.git
   cd bsc-bn
   ```

2. **Install dependencies:**

   ```sh
   go mod tidy
   ```

3. **Build the tool:**
   ```sh
   go build -o bsc-bn
   ```

## Usage

To process a Barnes & Noble PO HTML file and generate summary HTML files:

```sh
./bsc-bn -file <POFile.HTM>
```

- Replace `<POFile.HTM>` with the path to your PO HTML file.
- The tool will extract each PO in the file and generate a corresponding HTML summary in the `html/` directory, named by PO number (e.g., `html/123456.html`).

### Example

```sh
./bsc-bn -file sample_po.htm
```

This will create one or more files like `html/123456.html`, `html/123457.html`, etc., depending on how many POs are present in the input file.

## Output

- Output files are placed in the `html/` directory (created automatically if it doesn't exist).
- Each output file is a standalone HTML document with:
  - PO header (number, date, cancel after)
  - Bill To and Ship To addresses
  - Line items (EAN, Title, SKU, Arrival Date, Quantity, Cost, Retail, Discount, Casepack, IOQ)
  - Totals (lines, quantity, total cost)
  - Clean, print-friendly formatting

## Project Structure

- `main.go` — CLI entry point and main logic
- `po/` — PO parsing and HTML generation logic
- `html/` — Output directory for generated HTML summaries
- `files/` — (Unused by default, may be for future extensions)
- `go.mod`, `go.sum` — Go module/dependency files

---

**bsc-bn** is not affiliated with Barnes & Noble. This tool is provided as-is for convenience in managing and archiving PO data.
