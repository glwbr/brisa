package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glwbr/brisa/core"
	"github.com/glwbr/brisa/core/scraper"
)

func main() {
	// Mode flags
	modeFlag := flag.String("mode", "parse", "Mode: 'parse' (from file) or 'scrape' (from portal)")

	// Parse mode flags
	portalFlag := flag.String("portal", "BA", "NFC-e portal (e.g. BA, SP)")
	fileFlag := flag.String("file", "", "Path to HTML file to parse (parse mode)")

	// Scrape mode flags
	accessKeyFlag := flag.String("key", "", "NFC-e access key (scrape mode)")
	outputFlag := flag.String("output", "", "Output directory for scraped HTML (scrape mode)")
	captchaFileFlag := flag.String("captcha-output", "captcha.png", "Path to save captcha image (scrape mode)")

	flag.Parse()

	switch *modeFlag {
	case "parse":
		if *fileFlag == "" {
			fmt.Fprintln(os.Stderr, "missing --file for parse mode")
			os.Exit(1)
		}
		runParseMode(*portalFlag, *fileFlag)

	case "scrape":
		if *accessKeyFlag == "" {
			fmt.Fprintln(os.Stderr, "missing --key for scrape mode")
			os.Exit(1)
		}
		runScrapeMode(*portalFlag, *accessKeyFlag, *outputFlag, *captchaFileFlag)

	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", *modeFlag)
		os.Exit(1)
	}
}

func runParseMode(portal, filePath string) {
	html, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	receipt, err := core.ParseFromHTML(html, core.ParseOptions{
		Portal: core.Portal(portal),
	})
	if err != nil {
		log.Fatalf("failed to parse receipt: %v", err)
	}

	printReceipt(receipt)
}

func runScrapeMode(portal, accessKey, outputDir, captchaFile string) {
	if portal != "BA" {
		log.Fatalf("scraping only supported for BA portal currently")
	}

	ctx := context.Background()

	// Create scraper with interactive captcha solver
	baScraper, err := core.NewBAScraper(
		core.WithCaptchaSolver(&scraper.ManualSolver{
			PromptFunc: func(_ context.Context, challenge *scraper.CaptchaChallenge) (string, error) {
				// Save captcha to file
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
		}),
	)
	if err != nil {
		log.Fatalf("failed to create scraper: %v", err)
	}

	fmt.Printf("Fetching invoice: %s\n", accessKey)

	result, err := baScraper.FetchByAccessKey(ctx, accessKey)
	if err != nil {
		log.Fatalf("failed to fetch invoice: %v", err)
	}

	// Save HTML files if output directory specified
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("failed to create output directory: %v", err)
		}

		files := map[string][]byte{
			"danfe.html":    result.DanfeHTML,
			"nfe_tab.html":  result.NFETabHTML,
			"products.html": result.ProductsTabHTML,
		}

		for name, content := range files {
			path := outputDir + "/" + name
			if err := os.WriteFile(path, content, 0644); err != nil {
				log.Printf("warning: failed to save %s: %v", name, err)
			} else {
				fmt.Printf("Saved: %s\n", path)
			}
		}
	}

	printReceipt(result.Receipt)
}

func printReceipt(receipt *core.Receipt) {
	fmt.Printf("\n=== Invoice Details ===\n\n")
	fmt.Printf("Key: %s\n", receipt.Key)
	fmt.Printf("Portal: %s\n", receipt.Portal)

	if receipt.IssueDate.IsZero() {
		fmt.Println("IssueDate: -")
	} else {
		fmt.Printf("IssueDate: %s\n", receipt.IssueDate.Format(time.RFC3339))
	}

	fmt.Printf("Series: %s, Number: %s\n", receipt.Series, receipt.ReceiptNumber)

	fmt.Println("\nIssuer:")
	fmt.Printf("  Name: %s\n", receipt.Issuer.Name)
	fmt.Printf("  CNPJ: %s\n", receipt.Issuer.CNPJ)
	if receipt.Issuer.StateRegID != "" {
		fmt.Printf("  State Reg ID: %s\n", receipt.Issuer.StateRegID)
	}
	if receipt.Issuer.Address.State != "" {
		fmt.Printf("  State: %s\n", receipt.Issuer.Address.State)
	}

	fmt.Println("\nConsumer:")
	if receipt.Consumer.Document == "" && receipt.Consumer.Name == "" {
		fmt.Println("  (none)")
	} else {
		if receipt.Consumer.Document != "" {
			fmt.Printf("  Document: %s\n", receipt.Consumer.Document)
		}
		if receipt.Consumer.Name != "" {
			fmt.Printf("  Name: %s\n", receipt.Consumer.Name)
		}
	}

	fmt.Printf("\nSubtotal: %s\n", receipt.Subtotal.String())
	fmt.Printf("Discount: %s\n", receipt.Discount.String())
	fmt.Printf("Total: %s\n", receipt.Total.String())

	fmt.Println("\nPayments:")
	if len(receipt.Payments) == 0 {
		fmt.Println("  (none)")
	} else {
		for i, p := range receipt.Payments {
			fmt.Printf("  %d. %s - %s", i+1, p.Method, p.Amount.String())
			if p.Installments > 0 {
				fmt.Printf(" (%d installments)", p.Installments)
			}
			fmt.Println()
		}
	}

	fmt.Println("\nItems:")
	if len(receipt.Items) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, item := range receipt.Items {
			fmt.Printf("  %d. %s\n", item.LineNumber, item.Description)
			fmt.Printf("     Qty: %.3f %s x %s = %s\n",
				item.Quantity,
				item.Unit,
				item.UnitPrice.String(),
				item.Total.String(),
			)
			if item.Code != "" {
				fmt.Printf("     Code: %s", item.Code)
				if item.NCM != "" {
					fmt.Printf(" | NCM: %s", item.NCM)
				}
				if item.GTIN != "" {
					fmt.Printf(" | GTIN: %s", item.GTIN)
				}
				fmt.Println()
			}
		}
	}
}
