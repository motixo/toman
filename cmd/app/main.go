package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/motixo/toman/pkg/format"
	"github.com/motixo/toman/pkg/spinner"
)

type Target struct {
	Label    string
	Slug     string
	FlagName string
	Desc     string
	Enabled  *bool
}

type Result struct {
	Label string
	Price string
	Error error
}

func main() {
	targets := []*Target{
		{Label: "USD", Slug: "price_dollar_rl", FlagName: "usd", Desc: "Get United States Dollar price"},
		{Label: "EUR", Slug: "price_eur", FlagName: "eur", Desc: "Get Euro price"},
		{Label: "GOLD/COIN", Slug: "sekee", FlagName: "gold", Desc: "Get Gold Coin price"},
		{Label: "TETHER", Slug: "crypto-tether", FlagName: "tether", Desc: "Get Tether (USDT) price"},
	}

	for _, t := range targets {
		t.Enabled = flag.Bool(t.FlagName, false, t.Desc)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  By default, all currencies are fetched. Use flags to filter.")
		fmt.Fprintln(os.Stderr, "\nFlags:")

		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
		for _, t := range targets {
			fmt.Fprintf(w, "  -%s\t%s\n", t.FlagName, t.Desc)
		}
		w.Flush()
		fmt.Fprintln(os.Stderr, "")
	}

	flag.Parse()

	var targetsToFetch []*Target
	anyFlagSet := false

	for _, t := range targets {
		if *t.Enabled {
			anyFlagSet = true
			targetsToFetch = append(targetsToFetch, t)
		}
	}

	if !anyFlagSet {
		targetsToFetch = targets
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	spin := spinner.New()
	spin.Start()

	doc, err := fetchDocument(ctx, "https://www.tgju.org/")
	if err != nil {
		spin.Stop()
		log.Fatalf("Failed to fetch data: %v", err)
	}
	spin.Stop()
	results := processTargets(doc, targetsToFetch)
	printResults(results)
}

func processTargets(doc *goquery.Document, targets []*Target) []Result {
	var results []Result
	for _, t := range targets {
		price, err := extractPrice(doc, t.Slug)
		results = append(results, Result{
			Label: t.Label,
			Price: price,
			Error: err,
		})
	}
	return results
}

func printResults(results []Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, res := range results {
		if res.Error != nil {
			fmt.Fprintf(w, "%s\tError: %v\n", res.Label, res.Error)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\n", res.Label, res.Price)
	}
	w.Flush()
}

func fetchDocument(ctx context.Context, url string) (*goquery.Document, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fa-IR,fa;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func extractPrice(doc *goquery.Document, slug string) (string, error) {
	selector := fmt.Sprintf("tr[data-market-nameslug='%s']", slug)
	s := doc.Find(selector).First()

	if s.Length() == 0 {
		return "", fmt.Errorf("row not found")
	}

	irrPrice, exists := s.Attr("data-price")
	if !exists {
		return "", fmt.Errorf("data-price attribute missing")
	}

	return parseToToman(irrPrice)
}

func parseToToman(rawPrice string) (string, error) {
	clean := strings.ReplaceAll(rawPrice, ",", "")
	num, err := strconv.Atoi(clean)
	if err != nil {
		return "", fmt.Errorf("invalid number format: %s", rawPrice)
	}
	toman := num / 10
	return format.FormatWithCommas(toman), nil
}
