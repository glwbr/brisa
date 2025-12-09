package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glwbr/brisa/invoice"
	"github.com/glwbr/brisa/portal/ba"
	"github.com/glwbr/brisa/scraper"
	"github.com/glwbr/brisa/server"
)

func main() {
	mode := flag.String("mode", "parse", "Mode: 'parse' (from file), 'scrape' (from portal), or 'server' (http api)")
	portal := flag.String("portal", "BA", "NFC-e portal (e.g. BA)")
	file := flag.String("file", "", "Path to HTML file to parse (parse mode)")
	key := flag.String("key", "", "NFC-e access key (scrape mode)")
	output := flag.String("output", "", "Output directory for scraped HTML (scrape mode)")
	captchaFile := flag.String("captcha-output", "captcha.png", "Path to save captcha image (scrape mode)")
	addr := flag.String("addr", ":8080", "Server address (server mode)")

	flag.Parse()

	switch *mode {
	case "server":
		if err := server.NewServer().Start(*addr); err != nil {
			log.Fatal(err)
		}
	case "parse":
		if *file == "" {
			log.Fatal("missing --file for parse mode")
		}
		runParseMode(*portal, *file)
	case "scrape":
		if *key == "" {
			log.Fatal("missing --key for scrape mode")
		}
		runScrapeMode(*portal, *key, *output, *captchaFile)
	default:
		log.Fatalf("unknown mode: %s", *mode)
	}
}

func runParseMode(portalName, filePath string) {
	if portalName != "BA" {
		log.Fatalf("unsupported portal: %s", portalName)
	}

	html, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}

	// For now, assuming we are parsing the NFe tab HTML
	receipt, err := ba.ParseNFeTab(html)
	if err != nil {
		log.Fatalf("parse receipt: %v", err)
	}

	printReceipt(receipt)
}

func runScrapeMode(portalName, accessKey, outputDir, captchaFile string) {
	if portalName != "BA" {
		log.Fatalf("unsupported portal: %s", portalName)
	}

	ctx := context.Background()

	s, err := ba.New(ba.WithCaptchaSolver(&scraper.ManualSolver{
		PromptFunc: func(_ context.Context, challenge *scraper.CaptchaChallenge) (string, error) {
			if err := os.WriteFile(captchaFile, challenge.Image, 0644); err != nil {
				return "", fmt.Errorf("save captcha: %w", err)
			}
			fmt.Printf("Captcha saved to: %s\n", captchaFile)
			fmt.Print("Enter captcha solution: ")

			var solution string
			if _, err := fmt.Scanln(&solution); err != nil {
				return "", fmt.Errorf("read input: %w", err)
			}
			return solution, nil
		},
	}))
	if err != nil {
		log.Fatalf("failed to create scraper: %v", err)
	}

	fmt.Printf("Fetching invoice: %s\n", accessKey)

	result, err := s.FetchByAccessKey(ctx, accessKey)

	if err != nil {
		log.Fatalf("fetch invoice: %v", err)
	}

	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("create output dir: %v", err)
		}
		for name, content := range result.RawHTML {
			path := fmt.Sprintf("%s/%s.html", outputDir, name)
			if err := os.WriteFile(path, content, 0644); err != nil {
				log.Printf("warning: save %s: %v", name, err)
			} else {
				fmt.Printf("Saved: %s\n", path)
			}
		}
	}

	printReceipt(result.Receipt)
}

func printReceipt(r *invoice.Receipt) {
	fmt.Printf("\n=== Invoice Details ===\n\n")
	fmt.Printf("Key: %s\n", r.Key)
	fmt.Printf("Portal: %s\n", r.Portal)

	if !r.IssueDate.IsZero() {
		fmt.Printf("IssueDate: %s\n", r.IssueDate.Format(time.RFC3339))
	}

	fmt.Printf("Series: %s, Number: %s\n", r.Series, r.ReceiptNumber)

	fmt.Println("\nIssuer:")
	fmt.Printf("  Name: %s\n", r.Issuer.Name)
	fmt.Printf("  CNPJ: %s\n", r.Issuer.CNPJ)
	if r.Issuer.StateRegID != "" {
		fmt.Printf("  State Reg ID: %s\n", r.Issuer.StateRegID)
	}
	if r.Issuer.Address.State != "" {
		fmt.Printf("  State: %s\n", r.Issuer.Address.State)
	}

	fmt.Println("\nConsumer:")
	if r.Consumer.Document != "" {
		fmt.Printf("  Document: %s\n", r.Consumer.Document)
	}
	if r.Consumer.Name != "" {
		fmt.Printf("  Name: %s\n", r.Consumer.Name)
	}

	fmt.Printf("\nSubtotal: %s\n", r.Subtotal.String())
	fmt.Printf("Discount: %s\n", r.Discount.String())
	fmt.Printf("Total: %s\n", r.Total.String())

	fmt.Println("\nItems:")
	if len(r.Items) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, item := range r.Items {
			fmt.Printf("  %d. %s\n", item.LineNumber, item.Description)
			fmt.Printf("     Qty: %.3f %s x %s = %s\n",
				item.Quantity, item.Unit, item.UnitPrice.String(), item.Total.String())
		}
	}
}
