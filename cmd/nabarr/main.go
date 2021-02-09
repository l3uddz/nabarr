package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/goccy/go-yaml"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/l3uddz/nabarr/cmd/nabarr/rss"
	"github.com/l3uddz/nabarr/trakt"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
)

type config struct {
	Trakt trakt.Config       `yaml:"trakt"`
	Pvrs  []nabarr.PvrConfig `yaml:"pvrs"`
	Rss   rss.Config         `yaml:"rss"`
}

var (
	Version   string
	Timestamp string
	GitCommit string

	// CLI
	cli struct {
		globals

		// flags
		Config    string `type:"path" default:"${config_file}" env:"APP_CONFIG" help:"Config file path"`
		Cache     string `type:"path" default:"${cache_file}" env:"APP_CACHE" help:"Cache file path"`
		Log       string `type:"path" default:"${log_file}" env:"APP_LOG" help:"Log file path"`
		Verbosity int    `type:"counter" default:"0" short:"v" env:"APP_VERBOSITY" help:"Log level verbosity"`
	}
)

type globals struct {
	Version versionFlag `name:"version" help:"Print version information and quit"`
	Update  updateFlag  `name:"update" help:"Update if newer version is available and quit"`
}

func main() {
	// cli
	ctx := kong.Parse(&cli,
		kong.Name("nabarr"),
		kong.Description("Monitor newznab/torznab rss and add new media to sonarr/radarr"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
			Compact: true,
		}),
		kong.Vars{
			"version":     fmt.Sprintf("%s (%s@%s)", Version, GitCommit, Timestamp),
			"config_file": filepath.Join(defaultConfigPath(), "config.yml"),
			"cache_file":  filepath.Join(defaultConfigPath(), "cache"),
			"token_file":  filepath.Join(defaultConfigPath(), "token.json"),
			"log_file":    filepath.Join(defaultConfigPath(), "activity.log"),
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		return
	}

	// logger
	logger := log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}, zerolog.ConsoleWriter{
		Out: &lumberjack.Logger{
			Filename:   cli.Log,
			MaxSize:    5,
			MaxAge:     14,
			MaxBackups: 5,
		},
		NoColor: true,
	}))

	switch {
	case cli.Verbosity == 1:
		log.Logger = logger.Level(zerolog.DebugLevel)
	case cli.Verbosity > 1:
		log.Logger = logger.Level(zerolog.TraceLevel)
	default:
		log.Logger = logger.Level(zerolog.InfoLevel)
	}

	// config
	log.Trace().Msg("Initialising config")
	file, err := os.Open(cli.Config)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed opening config")
	}
	defer file.Close()

	cfg := config{}
	decoder := yaml.NewDecoder(file, yaml.Strict())
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Error().Msg("Failed decoding configuration")
		log.Fatal().
			Msg(err.Error())
	}

	// cache
	c, err := cache.New(cli.Cache)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed initialising cache")
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Error().
				Err(err).
				Msg("Failed closing cache gracefully")
		}
	}()

	// trakt
	log.Trace().Msg("Initialising trakt")
	t := trakt.New(&cfg.Trakt)

	// pvrs
	log.Trace().Msg("Initialising pvrs")
	pvrs := make(map[string]pvr.PVR, 0)
	for _, p := range cfg.Pvrs {
		// init pvr
		po, err := pvr.NewPVR(p, t, c)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("pvr", p.Name).
				Msg("Failed initialising pvr")
		}

		// start pvr processor
		po.Start()

		// add pvr to map
		pvrs[p.Name] = po
	}

	// rss
	log.Trace().Msg("Initialising rss")
	r := rss.New(cfg.Rss, pvrs)
	for _, feed := range cfg.Rss.Feeds {
		if err := r.AddJob(feed); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed initialising rss")
		}
	}
	r.Start()

	// wait for shutdown signal
	waitShutdown()

	// stop feed scheduler
	r.Stop()

	// stop pvr queue processors
	for _, p := range pvrs {
		p.Stop()
	}
}
