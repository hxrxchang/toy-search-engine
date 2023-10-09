package main

import (
	"fmt"
	"log"

	"github.com/kotaroooo0/stalefish"
)

func main() {
	db, err := stalefish.NewDBClient(stalefish.NewDBConfig("root", "password", "127.0.0.1", "3306", "stalefish"))
	if err != nil {
		log.Fatal(err)
	}
	storage := stalefish.NewStorageRdbImpl(db)
	analyzer := stalefish.NewAnalyzer([]stalefish.CharFilter{}, stalefish.NewStandardTokenizer(), []stalefish.TokenFilter{stalefish.NewLowercaseFilter()})

	indexer := stalefish.NewIndexer(storage, analyzer, 1)
	for _, body := range []string{"Ruby PHP JS", "Go Ruby", "Ruby Go PHP", "Go PHP"} {
		if err := indexer.AddDocument(stalefish.NewDocument(body)); err != nil {
			log.Fatal(err)
		}
	}

	sorter := stalefish.NewTfIdfSorter(storage)
	mq := stalefish.NewMatchQuery("GO Ruby", stalefish.OR, analyzer, sorter)
	mseacher := mq.Searcher(storage)
	result, err := mseacher.Search()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result) // [{2 Go Ruby 2} {3 Ruby Go PHP 3} {4 Go PHP 2} {1 Ruby PHP JS 3}]

	pq := stalefish.NewPhraseQuery("go RUBY", analyzer, nil)
	pseacher := pq.Searcher(storage)
	result, err = pseacher.Search()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result) // [{2 Go Ruby 2}
}
