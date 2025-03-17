package pkgElsticsearch

import "github.com/elastic/go-elasticsearch/v7"

func NewDefaultConfig(addresses []string) elasticsearch.Config {
	return elasticsearch.Config{
		Addresses:            addresses,
		EnableRetryOnTimeout: true, // Default: false.
		MaxRetries:           5,    // Default: 3.
	}
}
