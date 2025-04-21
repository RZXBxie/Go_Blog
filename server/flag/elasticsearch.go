package flag

import (
	"bufio"
	"fmt"
	"os"
	"server/elasticsearch"
	"server/service"
)

func ElasticSearch() error {
	esService := service.ServiceGroupApp.EsService
	indexExists, err := esService.IndexExist(elasticsearch.ArticleIndex())
	if indexExists {
		fmt.Println("The index already exists. Do you want to delete the data and recreate the index? (y/n)")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "y":
			fmt.Println("Proceeding to delete the data and recreate the index...")
			if err = esService.IndexDelete(elasticsearch.ArticleIndex()); err != nil {
				return err
			}
		case "n":
			fmt.Println("Exiting... the program")
			os.Exit(0)
		default:
			fmt.Println("Invalid input. Please enter 'y' to delete and recreate the index, or 'n' to exit.")
			return ElasticSearch()
		}
	}
	return esService.IndexCreate(elasticsearch.ArticleIndex(), elasticsearch.ArticleMapping())
}
