package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/goccy/go-yaml"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/build"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/l3uddz/nabarr/media"
	"github.com/l3uddz/nabarr/rss"
	"github.com/lefelys/state"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	Media media.Config       `yaml:"media"`
	Pvrs  []nabarr.PvrConfig `yaml:"pvrs"`
	Rss   rss.Config         `yaml:"rss"`
}

var (
	// CLI
	cli struct {
		globals

		// flags
		Config    string `type:"path" default:"${config_file}" env:"APP_CONFIG" help:"Config file path"`
		Cache     string `type:"path" default:"${cache_file}" env:"APP_CACHE" help:"Cache file path"`
		Log       string `type:"path" default:"${log_file}" env:"APP_LOG" help:"Log file path"`
		Verbosity int    `type:"counter" default:"0" short:"v" env:"APP_VERBOSITY" help:"Log level verbosity"`

		// commands
		Run  struct{} `cmd help:"Run"`
		Test struct {
			Pvr      string `type:"string" required:"1" help:"PVR to test item against" placeholder:"sonarr"`
			Id       string `type:"string" required:"1" help:"Metadata ID of item to test" placeholder:"tvdb:121361"`
			AllowAdd bool   `type:"bool" default:"0" required:"0" help:"Add item"`
		} `cmd help:"Test your filters and stop"`
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
			"version":     fmt.Sprintf("%s (%s@%s)", build.Version, build.GitCommit, build.Timestamp),
			"config_file": filepath.Join(defaultConfigPath(), "config.yml"),
			"cache_file":  filepath.Join(defaultConfigPath(), "cache"),
			"log_file":    filepath.Join(defaultConfigPath(), "activity.log"),
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		return
	}

	if ctx.Command() == "test" && cli.Verbosity == 0 {
		// default to debug verbosity in test mode
		cli.Verbosity = 1
	}

	// logger
	logger := log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out:        os.Stderr,
	}, zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
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
		log.Error().
			Err(err).
			Msg("Failed opening config")
		return
	}
	defer file.Close()

	cfg := config{}
	decoder := yaml.NewDecoder(file, yaml.Strict())
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Error().Msg("Failed decoding configuration")
		log.Error().Msg(err.Error())
		return
	}

	// cache
	c, err := cache.New(cli.Cache)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed initialising cache")
		return
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Error().
				Err(err).
				Msg("Failed closing cache gracefully")
		}
	}()

	// media
	log.Trace().Msg("Initialising media")
	m, err := media.New(&cfg.Media)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed initialising media")
		return
	}

	// states
	pvrStates := make([]state.State, 0)
	rssState := state.Empty()

	// pvrs
	log.Trace().Msg("Initialising pvrs")

	pvrs := make(map[string]pvr.PVR, 0)
	for _, p := range cfg.Pvrs {
		if ctx.Command() == "run" || (ctx.Command() == "test" && strings.EqualFold(cli.Test.Pvr, p.Name)) {
			// init pvr
			mode := ctx.Command()
			if ctx.Command() == "test" && cli.Test.AllowAdd {
				mode = "test-add"
			}

			po, err := pvr.NewPVR(p, mode, m, c)
			if err != nil {
				log.Error().
					Err(err).
					Str("pvr", p.Name).
					Msg("Failed initialising pvr")
				return
			}

			// start pvr processor
			pvrStates = append(pvrStates, po.Start())

			// add pvr to map
			pvrs[p.Name] = po
		}
	}

	// run mode (start rss scheduler and wait for shutdown signal)
	if ctx.Command() == "run" {
		// rss
		log.Trace().Msg("Initialising rss")
		r := rss.New(cfg.Rss, c, pvrs)
		for _, feed := range cfg.Rss.Feeds {
			if err := r.AddJob(feed); err != nil {
				log.Error().
					Err(err).
					Msg("Failed initialising rss")
				return
			}
		}
		rssState = r.Start()

		// wait for shutdown signal
		waitShutdown()
	} else {
		// test mode
		idParts := strings.Split(cli.Test.Id, ":")
		if len(idParts) < 2 {
			log.Error().
				Str("id", cli.Test.Id).
				Msg("An invalid id was provided")
			return
		}

		// prepare test item
		testItem := new(media.FeedItem)
		switch strings.ToLower(idParts[0]) {
		case "imdb":
			testItem.ImdbId = idParts[1]
		case "tmdb":
			testItem.TmdbId = idParts[1]
		case "tvdb":
			testItem.TvdbId = idParts[1]
		default:
			log.Error().
				Str("agent", idParts[0]).
				Str("id", idParts[1]).
				Msg("Unsupported agent was provided")
			return
		}

		// queue test item
		for _, p := range pvrs {
			// set test item title based on pvr type
			if p.Type() == "sonarr" {
				testItem.Title = "Test.Mode.S01E01.1080p.DTS-HD.MA.5.1.AVC.REMUX-FraMeSToR"
			} else {
				testItem.Title = "Test.Mode.2021.BluRay.1080p.TrueHD.Atmos.7.1.AVC.HYBRID.REMUX-FraMeSToR"
			}
			// queue test item
			p.QueueFeedItem(testItem)
		}

		// sleep for a moment
		time.Sleep(1 * time.Second)
	}

	// shutdown
	appCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	appState := state.Merge(pvrStates...).DependsOn(rssState)
	if err := appState.Shutdown(appCtx); err != nil {
		log.Error().
			Err(err).
			Msg("Failed shutting down gracefully")
		return
	}
}
