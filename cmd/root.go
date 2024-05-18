package cmd

import (
	"dilu-gateway/config"
	"dilu-gateway/proxy"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var (
	// Used for flags.
	cfgFile string
	//userLicense string

	rootCmd = &cobra.Command{
		Use:          "go-gateway -c config.yaml",
		Short:        "go-gateway",
		Long:         `go-gateway`,
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				tip()
				return errors.New("requires at least one arg")
			}
			return nil
		},
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		Run: func(cmd *cobra.Command, args []string) {
			tip()
		},
	}
)

func tip() {
	usageStr := `欢迎使用 go-gateway 查看命令：go-gateway --help`
	fmt.Printf("%s\n", usageStr)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./resources/config.dev.yaml", "go-gateway -c config.yaml")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	v := viper.New()
	v.SetConfigFile(cfgFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: %v \n", err))
	}

	var cfg config.AppConfig

	if err = v.Unmarshal(&cfg); err != nil {
		fmt.Println(err)
	}

	if cfg.RemoteCfg.Enable {
		rviper := viper.New()
		if cfg.RemoteCfg.SecretKeyring == "" {
			err = rviper.AddRemoteProvider(cfg.RemoteCfg.Provider, cfg.RemoteCfg.Endpoint, cfg.RemoteCfg.Path)
		} else {
			err = rviper.AddSecureRemoteProvider(cfg.RemoteCfg.Provider, cfg.RemoteCfg.Endpoint, cfg.RemoteCfg.Path, cfg.RemoteCfg.SecretKeyring)
		}
		if err != nil {
			panic(fmt.Sprintf("Fatal error remote config : %v \n", err))
		}
		rviper.SetConfigType(cfg.RemoteCfg.GetConfigType())
		err = rviper.ReadRemoteConfig()
		if err != nil {
			panic(fmt.Sprintf("Fatal error remote config : %v \n", err))
		}
		var remoteCfg config.AppConfig
		rviper.Unmarshal(&remoteCfg)
		mergeCfg(&cfg, &remoteCfg)

		go func() {
			for {
				time.Sleep(cfg.RemoteCfg.GetDuration()) // delay after each request
				err := rviper.WatchRemoteConfig()
				if err != nil {
					fmt.Println(err)
					continue
				}
				var watchCfg config.AppConfig
				err = rviper.Unmarshal(&watchCfg)
				if err != nil {
					slog.Error("watch", "err", err)
					continue
				}
				if !reflect.DeepEqual(remoteCfg, watchCfg) {
					slog.Debug("watch config changed", "config", watchCfg)
					//fmt.Println("watch config changed", remoteCfg, watchCfg)
					mergeCfg(&cfg, &watchCfg)
					// if !reflect.DeepEqual(remoteCfg.Logger, watchCfg.Logger) {
					// 	proxy.InitLog()
					// }
					if !reflect.DeepEqual(remoteCfg.RdConfig, watchCfg.RdConfig) {
						proxy.InitRd()
					}
					if !reflect.DeepEqual(remoteCfg.JWT, watchCfg.JWT) ||
						!reflect.DeepEqual(remoteCfg.Rules, watchCfg.Rules) ||
						!reflect.DeepEqual(remoteCfg.Extend, watchCfg.Extend) {
						proxy.InitHandler()
					}
					remoteCfg = watchCfg
				}
			}
		}()
	} else {
		mergeCfg(&cfg, nil)
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			slog.Debug("config file changed", "config", e.String())
			if err = v.Unmarshal(&cfg); err != nil {
				fmt.Println(err)
			}
			mergeCfg(&cfg, nil)
		})
	}
	proxy.Run()
}

func mergeCfg(local, remote *config.AppConfig) {
	if remote != nil {
		proxy.Cfg = local
		//proxy.Cfg.Logger = remote.Logger
		proxy.Cfg.JWT = remote.JWT
		proxy.Cfg.Rules = remote.Rules
		proxy.Cfg.Extend = remote.Extend
		proxy.Cfg.RdConfig = remote.RdConfig
	} else {
		proxy.Cfg = local
	}
}
