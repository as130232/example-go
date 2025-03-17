package config

import "time"

type JobLock struct {
	Key        string
	Expiration time.Duration
}
