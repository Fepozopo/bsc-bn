package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Fepozopo/bsc-bn/po"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	var filename string
	flag.StringVar(&filename, "file", "", "Path to the PO HTML file")
	flag.Parse()

	if filename == "" {
		fmt.Println("Usage: bsc-bn -file <POFile.HTM>")
		os.Exit(1)
	}

	f, err := os.Open(filename)
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
			poObj := po.ExtractPO(s)
			if err := po.WritePOHTML(poObj); err != nil {
				log.Printf("Failed to write PO HTML for %s: %v", poObj.Number, err)
			}
		}
	})
}
