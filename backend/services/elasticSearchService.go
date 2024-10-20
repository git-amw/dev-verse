package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/git-amw/backend/models"
)

type ElasticSearchProvider interface {
	IndexNewBlog(doc map[string]interface{}, id uint) error
	IndexNewTag(doc map[string]interface{}, id uint)
	SearchBlogs(searchTitle string, searchContent string, tagId int) []models.BlogSearchResponse
	SearchTags(tagValue string) models.TagSearchResponse
	UpdateTagDoc(doc map[string]interface{}, id uint)
}

type ElasticSearchService struct {
	ESClient *elasticsearch.Client
}

func NewElasticSearchService(esclient *elasticsearch.Client) ElasticSearchProvider {
	return &ElasticSearchService{
		ESClient: esclient,
	}
}

func (es *ElasticSearchService) IndexNewBlog(doc map[string]interface{}, id uint) error {
	docJSON, err := json.Marshal(doc)
	if err != nil {
		log.Fatal("Err in marshaling of doc ", err)
	}
	req := esapi.IndexRequest{
		Index:      "blog-index",
		DocumentID: strconv.Itoa(int(id)),
		Body:       strings.NewReader(string(docJSON)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.ESClient)
	if err != nil {
		log.Println("Error indexing document: ", err)
		return err
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		log.Printf("Error in response: %s", res.String())
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}
	log.Println("Document indexed successfully: ", res.String())
	return nil
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

func (es *ElasticSearchService) SearchBlogs(searchTitle string, searchContent string, tagId int) []models.BlogSearchResponse {
	// Construct the query as a string
	query := fmt.Sprintf(`{
        "query": {
            "bool": {
                "should": [
                    {
                        "match": {
                            "title": "%s"
                        }
                    },
                    {
                        "match": {
                            "content": "%s"
                        }
                    },
                    {
                        "nested": {
                            "path": "blogTags",
                            "query": {
                                "match": {
                                    "blogTags.tag_id": %d
                                }
                            }
                        }
                    }
                ],
                "minimum_should_match": 1
            }
        }
    }`, searchTitle, searchContent, tagId)

	var blogsdata []models.BlogSearchResponse
	var blog models.BlogSearchResponse
	resData := SearchQuery(es, query, "blog-index")
	for _, hit := range resData["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})
		source := doc["_source"]
		// log.Println(source)
		sourceBytes, _ := json.Marshal(source)    // Convert the _source back to JSON
		err := json.Unmarshal(sourceBytes, &blog) // Unmarshal the JSON into the struct
		if err != nil {
			log.Println("Error:", err)
		}
		log.Println(blog)
		blogsdata = append(blogsdata, blog)
	}
	return blogsdata
}

func (es *ElasticSearchService) SearchTags(tagValue string) models.TagSearchResponse {
	query := fmt.Sprintf(`{
		"query": {
			"bool": {
				"should": [
					"match": {
						"tag_value": %s
					}
				]
			}
		}
	 }`, tagValue)
	log.Println(query)
	// resData := SearchQuery(es, query, "tag-index")
	var tag models.TagSearchResponse
	/* hit := resData["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hit) > 0 {
		doc := hit[0].(map[string]interface{})
		log.Println(doc)
		source := doc["_source"]
		sourceBytes, _ := json.Marshal(source) // Convert the _source back to JSON
		json.Unmarshal(sourceBytes, &tag)      // Unmarshal the JSON into the struct
		log.Printf("Tag ID: %d, Tag Value: %s, Blogs with Tag: %d\n", tag.TagId, tag.TagValue, tag.BlogsWithTag)
	} */
	return tag
}

func SearchQuery(es *ElasticSearchService, query string, index string) map[string]interface{} {
	res, err := es.ESClient.Search(
		es.ESClient.Search.WithContext(context.Background()),
		es.ESClient.Search.WithIndex(index), // Replace with your index name
		es.ESClient.Search.WithBody(strings.NewReader(query)),
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
	return resData
}
