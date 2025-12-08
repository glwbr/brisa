package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Document wraps goquery.Document for HTML parsing.
type Document struct {
	*goquery.Document
}

// ParseHTML creates a Document from HTML bytes.
func ParseHTML(data []byte) (*Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	return &Document{doc}, nil
}

// HasElement returns true if selector matches at least one element.
func (d *Document) HasElement(selector string) bool {
	return d.Find(selector).Length() > 0
}

// Text extracts trimmed text from the first matching element.
func (d *Document) Text(selector string) string {
	return strings.TrimSpace(d.Find(selector).Text())
}

// Attr extracts an attribute value from the first matching element.
func (d *Document) Attr(selector, attr string) string {
	val, _ := d.Find(selector).Attr(attr)
	return val
}

// CollectLabelValues extracts label-span pairs from a selection.
func CollectLabelValues(sel *goquery.Selection, cache map[*html.Node]string) map[string]string {
	values := map[string]string{}

	var walk func(*goquery.Selection)
	walk = func(node *goquery.Selection) {
		node.Children().Each(func(_ int, child *goquery.Selection) {
			if goquery.NodeName(child) == "label" {
				label := CachedText(child, cache)
				valSel := child.Next()
				for valSel.Length() > 0 && goquery.NodeName(valSel) != "span" {
					valSel = valSel.Next()
				}
				if value := CachedText(valSel, cache); label != "" && value != "" {
					if _, ok := values[label]; !ok {
						values[label] = value
					}
				}
			}
			walk(child)
		})
	}
	walk(sel)
	return values
}

// CachedText returns normalized text, using cache for efficiency.
func CachedText(sel *goquery.Selection, cache map[*html.Node]string) string {
	if sel.Length() == 0 {
		return ""
	}
	node := sel.Get(0)
	if node == nil {
		return ""
	}
	if cached, ok := cache[node]; ok {
		return cached
	}
	text := normalizeText(sel.Text())
	cache[node] = text
	return text
}

func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\u00A0", " ")
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return strings.Join(strings.Fields(s), " ")
}
