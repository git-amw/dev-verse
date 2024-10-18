package db

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func ESClientConnection() *elasticsearch.Client {
	username := os.Getenv("ESUSER")
	password := os.Getenv("ESPASSWORD")
	url := os.Getenv("ESURL")
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip certificate verification
			},
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic("Error in creating ESClient")
	}
	res, err := client.Info()
	if err != nil {
		panic("Error in connecting with Elastic Search " + err.Error())
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Fatal("error in url")
	} else {
		log.Println("response body ", res)
	}
	return client
}

func ESCheackIndexExists(client *elasticsearch.Client) {
	indices := GetIndices()
	for indexName, mapping := range indices {
		res, err := client.Indices.Exists([]string{indexName})
		if err != nil {
			panic("Error in checking of index" + err.Error())
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusNotFound {
			log.Println(indexName)
			res, err := client.Indices.Create(indexName, client.Indices.Create.WithBody(strings.NewReader(mapping)))
			if err != nil {
				panic("error in index creation")
			}
			if res.IsError() {
				log.Fatalln("Error in index creation response ", res.StatusCode)
			} else {
				log.Println("Index created!!")
			}
		}
		if res.StatusCode == http.StatusOK {
			log.Println("exists - ", indexName)
		}
	}
}

func GetIndices() map[string]string {
	indices := map[string]string{
		"tag-index": `{
			"settings": {
				"number_of_shards": 2,
				"number_of_replicas": 1
			},
			"mappings": {
				"properties": {
					"id": { "type": "integer" },
					"tag_value": { "type": "text" },
					"blogs_with_tag": { "type": "integer" }
				}
			}
		}`,
		"blog-index": `{
			"settings": {
				"number_of_shards": 5,
				"number_of_replicas": 2
			},
			"mappings": {
    			"properties": {
					"id": 		   { "type": "integer" },
      				"title":       { "type": "text" },
      				"content":     { "type": "text" },
      				"likes":       { "type": "integer" },
     				"blogTags":    { 
						"type": "nested" 
						"properties": {
          					"tag_id":   { "type": "integer" },
        				}
		   			},
					"is_deleted": {         // New field to indicate soft delete status
        				"type": "boolean",
       					"null_value": false    // Defaults to false if not provided
      				}
    			}
  			}
		}`,
	}
	return indices
}
