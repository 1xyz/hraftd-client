package config

import "time"

type Config struct {
	URL     string
	Timeout time.Duration
}
