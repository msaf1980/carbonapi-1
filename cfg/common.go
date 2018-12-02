package cfg

import (
	"io"
	"time"

	"github.com/lomik/zapwriter"
	"gopkg.in/yaml.v2"
)

var DEBUG bool = false

type GraphiteConfig struct {
	Pattern  string
	Host     string
	Interval time.Duration
	Prefix   string
}

func ParseCommon(r io.Reader) (Common, error) {
	d := yaml.NewDecoder(r)
	d.SetStrict(DEBUG)

	c := DefaultConfig
	err := d.Decode(&c)

	return c, err
}

type Backends2 struct {
	groupName string  `yaml:"groupName"`
	servers map[string]int `yaml:"servers"`
}

type Common struct {
	Listen         string      `yaml:"listen"`
	ListenInternal string      `yaml:"listenInternal"`
	Backends       []string    `yaml:"backends"`
	Backends2      []Backends2 `yaml:"backends2"`

	MaxProcs                  int           `yaml:"maxProcs"`
	Timeouts                  Timeouts      `yaml:"timeouts"`
	ConcurrencyLimitPerServer int           `yaml:"concurrencyLimit"`
	KeepAliveInterval         time.Duration `yaml:"keepAliveInterval"`
	MaxIdleConnsPerHost       int           `yaml:"maxIdleConnsPerHost"`

	ExpireDelaySec             int32   `yaml:"expireDelaySec"`
	GraphiteWeb09Compatibility bool    `yaml:"graphite09compat"`
	CorruptionThreshold        float64 `yaml:"corruptionThreshold"`

	Buckets  int                `yaml:"buckets"`
	Graphite GraphiteConfig     `yaml:"graphite"`
	Logger   []zapwriter.Config `yaml:"logger"`
}

type Timeouts struct {
	Global       time.Duration `yaml:"global"`
	AfterStarted time.Duration `yaml:"afterStarted"`
	Connect      time.Duration `yaml:"connect"`
}

var DefaultConfig = Common{
	Listen:         ":8080",
	ListenInternal: ":7080",

	MaxProcs: 1,
	Timeouts: Timeouts{
		Global:       10000 * time.Millisecond,
		AfterStarted: 2 * time.Second,
		Connect:      200 * time.Millisecond,
	},
	ConcurrencyLimitPerServer: 20,
	KeepAliveInterval:         30 * time.Second,
	MaxIdleConnsPerHost:       100,

	ExpireDelaySec: int32(10 * time.Minute / time.Second),

	Buckets: 10,
	Graphite: GraphiteConfig{
		Interval: 60 * time.Second,
		Host:     "127.0.0.1:3002",
		Prefix:   "carbon.zipper",
		Pattern:  "{prefix}.{fqdn}",
	},
	Logger: []zapwriter.Config{DefaultLoggerConfig},
}

var DefaultLoggerConfig = zapwriter.Config{
	Logger:           "",
	File:             "stdout",
	Level:            "info",
	Encoding:         "console",
	EncodingTime:     "iso8601",
	EncodingDuration: "seconds",
}
