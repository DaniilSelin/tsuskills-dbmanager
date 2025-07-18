package search

import (
    "fmt"
    opensearch "github.com/opensearch-project/opensearch-go"
)

func NewClient(cfg opensearch.Config) (*opensearch.Client, error) {
    return opensearch.NewClient(cfg)
}

func Ping(client *opensearch.Client) error {
    res, err := client.Ping()
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("ping failed: %s", res.String())
    }
    return nil
}
