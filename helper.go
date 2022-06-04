package main

import (
	"log"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/anaskhan96/soup"
)

type Job struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	Status   int    `json:"status"`
	Result   int    `json:"result"`
}

type Page struct {
	Link            string
	Content         string
	SimilarityScore int
	Title           string
	Description     string
}

type PageList []Page

func (a PageList) Len() int { return len(a) }
func (a PageList) Less(i, j int) bool {
	return a[i].SimilarityScore > a[j].SimilarityScore
}
func (a PageList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetContentFromURL(url string) string {
	resp, err := soup.Get(url)
	if err != nil {
		return ""
	}

	doc := soup.HTMLParse(resp)
	ps := doc.FindAll("p")
	var content string
	for _, p := range ps {
		content += p.Text()
	}
	return content
}

func GetSimilarityScore(content string, text string) int {
	oc := metrics.NewOverlapCoefficient()
	return int(strutil.Similarity(content, text, oc) * 100)
}
