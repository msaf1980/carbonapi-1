package cfg

import (
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

func ParseAPIConfig(r io.Reader) (API, error) {
	d := yaml.NewDecoder(r)
	d.SetStrict(DEBUG)

	pre := preAPI{
		API:       DefaultAPIConfig,
		Upstreams: DefaultConfig,
	}
	err := d.Decode(&pre)
	if err != nil {
		return API{}, err
	}

	api := pre.API

	// Backwards compatibility is king
	if pre.Concurrency > 0 {
		api.ConcurrencyLimitPerServer = pre.Concurrency
	}

	if pre.CPUs > 0 {
		api.MaxProcs = pre.CPUs
	}

	if pre.IdleConnections > 0 {
		api.MaxIdleConnsPerHost = pre.IdleConnections
	}

	if pre.Upstreams.Buckets != DefaultConfig.Buckets {
		api.Buckets = pre.Upstreams.Buckets
	}

	// Any value set to a non-default in a nested structure means we pick all
	// values from that structure, for the sanity of the ops people.
	if pre.Upstreams.Timeouts != DefaultConfig.Timeouts {
		api.Timeouts = pre.Upstreams.Timeouts
	}

	if len(pre.Upstreams.Backends) >= 1 {
		api.Backends = pre.Upstreams.Backends
	}

	if len(pre.Upstreams.Backends2) >= 1 {
		api.Backends2 = pre.Upstreams.Backends2
	}

	return api, nil
}

var DefaultAPIConfig = defaultAPIConfig()

func defaultAPIConfig() API {
	cfg := API{
		Zipper: DefaultZipperConfig,

		ExtrapolateExperiment: false,
		SendGlobsAsIs:         false,
		AlwaysSendGlobsAsIs:   false,
		MaxBatchSize:          100,
		Cache: CacheConfig{
			Type:              "mem",
			DefaultTimeoutSec: 60,
		},
	}

	cfg.Listen = ":8081"
	cfg.MaxProcs = 0
	cfg.Graphite.Prefix = "carbon.api"

	return cfg
}

type API struct {
	Zipper `yaml:",inline"`

	ExtrapolateExperiment   bool          `yaml:"extrapolateExperiment"`
	SendGlobsAsIs           bool          `yaml:"sendGlobsAsIs"`
	AlwaysSendGlobsAsIs     bool          `yaml:"alwaysSendGlobsAsIs"`
	MaxBatchSize            int           `yaml:"maxBatchSize"`
	Cache                   CacheConfig   `yaml:"cache"`
	TimezoneString          string        `yaml:"tz"`
	PidFile                 string        `yaml:"pidFile"`
	BlockHeaderFile         string        `yaml:"blockHeaderFile"`
	BlockHeaderUpdatePeriod time.Duration `yaml:"blockHeaderUpdatePeriod"`
	HeadersToLog            []string      `yaml:"headersToLog"`

	UnicodeRangeTables  []string          `yaml:"unicodeRangeTables"`
	IgnoreClientTimeout bool              `yaml:"ignoreClientTimeout"`
	DefaultColors       map[string]string `yaml:"defaultColors"`
	FunctionsConfigs    map[string]string `yaml:"functionsConfig"`
}

type CacheConfig struct {
	Type              string   `yaml:"type"`
	Size              int      `yaml:"size_mb"`
	MemcachedServers  []string `yaml:"memcachedServers"`
	DefaultTimeoutSec int32    `yaml:"defaultTimeoutSec"`
}

type preAPI struct {
	API             `yaml:",inline"`
	Concurrency     int    `yaml:"concurency"`
	CPUs            int    `yaml:"cpus"`
	IdleConnections int    `yaml:"idleConnections"`
	Upstreams       Common `yaml:"upstreams"`
}
