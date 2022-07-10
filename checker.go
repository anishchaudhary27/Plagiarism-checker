package main

import (
	"database/sql"
	"log"
	"os/exec"
	"sort"
	"strings"

	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/streadway/amqp"
	"google.golang.org/api/customsearch/v1"
)

const maxPhrases = 20
const maxSearchResults = 2

func handleMsg(msg []string, db *sql.DB,
	customsearchService *customsearch.Service, customSearchEngineID string) {
	location := msg[0]
	jobId := msg[1]
	log.Printf("Processing job: %s", jobId)
	rawText, err := exec.Command("pdf2txt", "./files/"+location).Output()
	if err != nil {
		//update failed status in database and return
		_, err = db.Exec("UPDATE jobs SET status = $1, result = $2 WHERE id = $3", -1, 0, jobId)
		FailOnError(err, "Failed to update job status")
		log.Printf("Failed to process job %s", jobId)
		return
	}
	text := string(rawText)

	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()
	tr.Populate(text, language, rule)
	tr.Ranking(algorithmDef)
	rankedPhrases := textrank.FindPhrases(tr)

	var pageList PageList

	i := 0
	var added map[string]bool = make(map[string]bool)
	for _, phrase := range rankedPhrases {
		searchResult, err := customsearchService.Cse.List().Q(phrase.Left + " " + phrase.Right).Cx(customSearchEngineID).Do()
		if err != nil {
			continue
		}
		j := 0
		for _, item := range searchResult.Items {
			if added[item.Link] {
				continue
			}
			added[item.Link] = true
			pageList = append(pageList, Page{Link: item.Link, Title: item.Title, Description: item.Snippet})
			j++
			if j > maxSearchResults {
				break
			}
		}

		i++
		if i > maxPhrases {
			break
		}
	}

	for i, page := range pageList {
		pageList[i].Content = GetContentFromURL(page.Link)
		pageList[i].SimilarityScore = GetSimilarityScore(pageList[i].Content, text)
	}

	sort.Sort(pageList)
	k := 0
	maxSimilarity := 0
	for j := 0; j < pageList.Len(); j++ {
		if k > 4 {
			break
		}
		if pageList[j].SimilarityScore < 10 {
			continue
		}
		_, err = db.Exec("INSERT INTO refs (jobId, link, title, description, similarity) VALUES ($1, $2, $3, $4, $5)",
			jobId, pageList[j].Link, pageList[j].Title, pageList[j].Description, pageList[j].SimilarityScore)
		FailOnError(err, "Failed to insert ref")
		if pageList[j].SimilarityScore > maxSimilarity {
			maxSimilarity = pageList[j].SimilarityScore
		}
		k++
	}
	_, err = db.Exec("UPDATE jobs SET status = $1, result = $2 WHERE id = $3", 1, maxSimilarity, jobId)
	FailOnError(err, "Failed to update job status")
	log.Printf("Processed job: %s", jobId)
}

func Checker(ch *amqp.Channel, q amqp.Queue, db *sql.DB,
	customsearchService *customsearch.Service, customSearchEngineID string) {
	//return
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")
	var forever chan struct{}
	go func() {
		for d := range msgs {
			msg := strings.Split(string(d.Body), "$")
			go handleMsg(msg, db, customsearchService, customSearchEngineID)
		}
	}()
	log.Printf("Checker started")
	<-forever
}
