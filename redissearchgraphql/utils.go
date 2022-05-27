package redissearchgraphql

import (
	"fmt"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func checkIndexNames(indices []string) error {
	for _, x := range indices {
		if strings.Contains(x, ":") || strings.Contains(x, "-") {
			err := fmt.Errorf("Index name %s is not GraphQL compliant with a : or - in name", x)
			return err
		}
	}
	return nil
}

// CheckIndexNames ensures we don't have indices that are not GraphQL compliant objects
func GetIndices(searchClient *redisearch.Client) ([]string, error) {
	indices, err := searchClient.List()
	if err != nil {
		return nil, err
	}
	err = checkIndexNames(indices)
	return indices, err
}
