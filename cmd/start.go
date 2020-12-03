package cmd

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis"
	"github.com/lucbarr/leaderboard-manager/api"
	leaderboardBLL "github.com/lucbarr/leaderboard-manager/bll/leaderboard"
	"github.com/lucbarr/leaderboard-manager/dal"
	"github.com/lucbarr/leaderboard-manager/handler"
	"github.com/lucbarr/leaderboard-manager/handler/healthcheck"
	"github.com/lucbarr/leaderboard-manager/handler/leaderboard"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "LB"
)

var (
	configPath     string
	verbosityLevel string

	config      *viper.Viper
	logger      *logrus.Logger
	redisClient *redis.Client
	allDAL      *dal.All
)

var startCmd = &cobra.Command{
	Use:   "api",
	Short: "api command",

	Run: func(cmd *cobra.Command, args []string) {
		err := configure()
		if err != nil {
			panic(err)
		}
		start()

		select {}
	},
}

func start() {
	logger.Debug("starting app")

	allDAL = dal.NewAll(redisClient, config)

	debugHandler := healthcheck.NewHandler(logger)

	leaderboardBLL := leaderboardBLL.NewBLL(allDAL)

	leaderboardHandler := leaderboard.NewHandler(logger, leaderboardBLL, config)

	handlers := &handler.All{
		Debug:       debugHandler,
		Leaderboard: leaderboardHandler,
	}

	api := api.NewAPI(config, handlers)

	api.Start()
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func configure() error {
	readFlags()

	err := readConfigs(configPath)
	if err != nil {
		return err
	}

	err = configureLogger()
	if err != nil {
		return err
	}

	err = configureRedis(config)
	if err != nil {
		return err
	}

	return nil
}

func configureRedis(cfg *viper.Viper) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr: cfg.GetString("redis.url"),
	})

	logger.Debug("connected to redis")

	cmd := redisClient.Ping()
	return cmd.Err()
}

func configureLogger() error {
	logger = logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(verbosityLevel)
	if err != nil {
		return err
	}

	logger.SetLevel(level)

	return nil
}

func readFlags() {
	flag.StringVar(&configPath, "config", "cfg/default.yaml", "path for config yaml file")
	flag.StringVar(&verbosityLevel, "level", "debug", "verbosity level: info, debug, error, panic, fatal")

	flag.Parse()
}

func readConfigs(configPath string) error {
	cfg := viper.New()

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}

	cfgData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	cfg.SetConfigType("yaml")
	cfg.SetEnvPrefix(appName)
	cfg.ReadConfig(bytes.NewBuffer(cfgData))

	config = cfg

	return nil
}
