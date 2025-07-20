package search

import (
    "fmt"
    opensearch "github.com/opensearch-project/opensearch-go"
)

func NewClient(cfg opensearch.Config) (*opensearch.Client, error) {
    return opensearch.NewClient(cfg)
}

func Ping(client *opensearch.Client) error {
    res, err := client.Info()
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.StatusCode >= 400 {
        return fmt.Errorf("ping returned bad status: %d", res.StatusCode)
    }
    return nil
}
