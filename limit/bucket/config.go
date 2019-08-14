package limitBucket

import (
	"encoding/json"
)

type Config struct {
	Limit    int
	Interval float64
	changed  bool
}

func (c *Config) Bytes() ([]byte, error) {
	return json.Marshal(c)
}

func ConfigFromJSON(b []byte) (*Config, error) {
	c := &Config{}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}
