package main

import (
	"fmt"
	"github.com/adrg/xdg"
	"github.com/alexflint/go-arg"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type RunnableCmd interface {
	Run(c *elastic.Client, l *logrus.Logger) error
}

type InstanceConfig struct {
	URL *string `yaml:"url"`
	Sniff *bool `yaml:"sniff" default:"false"`
}

type Config struct {
	Instance *InstanceConfig `yaml:"instance"`
}

var args struct {
	Index *IndexCmd `arg:"subcommand:index"`
	InstanceURL string `arg:"--url,-u"`
	Sniff *bool `arg:"--sniff,-s" help:"Sniff nodes, then connect to them"`
	LogLevel string `arg:"env:LOG_LEVEL" default:"info"`
}

func main(){
	logger := logrus.New()
	arg.MustParse(&args)

	switch strings.ToLower(args.LogLevel) {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	}

	// Load config
	fileConfig, err := loadConfig()
	if err != nil {
		logger.Warnf("unable to load config: %v", err)
	}

	if fileConfig != nil {
		// Replace args
		if args.InstanceURL == "" && fileConfig.Instance != nil {
			if fileConfig.Instance.URL != nil {
				args.InstanceURL = *fileConfig.Instance.URL
			}

			if fileConfig.Instance.Sniff != nil {
				args.Sniff = fileConfig.Instance.Sniff
			}
		}
	}

	// Init ES Client
	c, err := elastic.NewClient(
		elastic.SetURL(),
		elastic.SetSniff(false),
	)
	if err != nil {
		logger.Fatalf("unable to create ES client: %v", err)
	}

	if args.Index != nil {
		// Index Command selected
		err = args.Index.Run(c, logger)
	}
	if err != nil {
		logger.Fatal(err)
	}
}

func loadConfig() (*Config, error) {
	filePath, err := xdg.ConfigFile("elastically/config.yml")
	if err != nil {
		return nil, fmt.Errorf("unable to get config path: %v", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %v", err)
	}

	var config Config

	dec := yaml.NewDecoder(f)
	dec.SetStrict(true)
	err = dec.Decode(&config)

	if err != nil {
		return nil, fmt.Errorf("unable to decode YAML config file: %v", err)
	}

	return &config, nil
}