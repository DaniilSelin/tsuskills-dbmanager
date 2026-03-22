package opensearch

import (
	"fmt"
	"time"

	"tsuskills-dbmanager/config"

	osclient "github.com/opensearch-project/opensearch-go"
)

func NewClient(cfg osclient.Config) (*osclient.Client, error) {
	return osclient.NewClient(cfg)
}

func Ping(client *osclient.Client, cfg config.SearchConnect) error {
	var lastErr error
	for i := 0; i < cfg.Retries; i++ {
		res, err := client.Info()
		if err == nil && res.StatusCode < 400 {
			res.Body.Close()
			return nil
		}
		if res != nil {
			res.Body.Close()
			lastErr = fmt.Errorf("ping returned bad status: %d", res.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(cfg.Delay)
	}
	return fmt.Errorf("all ping attempts failed: %w", lastErr)
}
