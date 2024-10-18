package services

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticSearchProvider interface {
	IndexNewBlog()
	IndexNewTag(doc map[string]interface{}, id uint)
	SearchBlogs()
	SearchTags()
	UpdateTagDoc(doc map[string]interface{}, id uint)
}

type ElasticSearchService struct {
	ESClient *elasticsearch.Client
}

type Tag struct {
	TagValue     string `json:"tag_value"`
	BlogsWithTag int    `json:"blogs_with_tag"`
}

func NewElasticSearchService(esclient *elasticsearch.Client) ElasticSearchProvider {
	return &ElasticSearchService{
		ESClient: esclient,
	}
}

func (es *ElasticSearchService) IndexNewBlog() {

}

func (es *ElasticSearchService) UpdateTagDoc(doc map[string]interface{}, id uint) {
	docJSON, err := json.Marshal(doc)
	if err != nil {
		log.Fatal("Err in marshaling of doc ", err)
	}

	req := esapi.UpdateRequest{
		Index:      "tag-index",
		DocumentID: strconv.Itoa(int(id)),
		Body:       strings.NewReader(string(docJSON)),
		Refresh:    "true", // Optional: refresh the index after the update
	}

	res, err := req.Do(context.Background(), es.ESClient)
	if err != nil {
		log.Fatal("Error updating document: ", err)
	}
	defer res.Body.Close()
	// Check for errors in the response
	if res.IsError() {
		log.Fatalf("Error in response: %s", res.String())
	}
	log.Println("Document indexed successfully: ", res.String())
}

func (es *ElasticSearchService) IndexNewTag(doc map[string]interface{}, id uint) {
	docJSON, err := json.Marshal(doc)
	if err != nil {
		log.Fatal("Err in marshaling of doc ", err)
	}

	req := esapi.IndexRequest{
		Index:      "tag-index",
		DocumentID: strconv.Itoa(int(id)),
		Body:       strings.NewReader(string(docJSON)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.ESClient)
	if err != nil {
		log.Fatal("Error indexing document: ", err)
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		log.Fatalf("Error in response: %s", res.String())
	}
	log.Println("Document indexed successfully: ", res.String())
}

func (es *ElasticSearchService) SearchBlogs() {

}

func (es *ElasticSearchService) SearchTags() {
	res, err := es.ESClient.Search(
		es.ESClient.Search.WithContext(context.Background()),
		es.ESClient.Search.WithIndex("tag-index"), // Replace with your index name
		es.ESClient.Search.WithBody(strings.NewReader(`{
			"query": {
				"match_all": {}
			}
		}`)),
		es.ESClient.Search.WithSize(32),
		es.ESClient.Search.WithTrackTotalHits(true), // Track total hits for pagination if needed
		es.ESClient.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response from Elasticsearch: %s", err)
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}
	var resData map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// Iterate through the hits and print each document
	for _, hit := range resData["hits"].(map[string]interface{})["hits"].([]interface{}) {
		// Unmarshal the _source field into the Tag struct
		var tag Tag
		doc := hit.(map[string]interface{})
		source := doc["_source"]
		sourceBytes, _ := json.Marshal(source) // Convert the _source back to JSON
		json.Unmarshal(sourceBytes, &tag)      // Unmarshal the JSON into the struct

		// Print the tag information
		log.Printf("Tag ID: %s, Tag Value: %s, Blogs with Tag: %d\n", doc["_id"], tag.TagValue, tag.BlogsWithTag)
	}
}
