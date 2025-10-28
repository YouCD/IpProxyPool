package cmd

import (
	"IpProxyPool/api"
	"IpProxyPool/common"
	"IpProxyPool/middleware/config"
	"IpProxyPool/middleware/database"
	"IpProxyPool/run"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/youcd/toolkit/log"
)

const name = "IpProxyPool"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   name,
	Short: "Main application",
	Run: func(_ *cobra.Command, _ []string) {
		if config.ConfigFile == "" {
			fmt.Println("config file is empty")
			os.Exit(1)
		}
		setting := config.ServerSetting

		// 信号处理应该放在服务器启动之后
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer func() {
			log.Info("shutting down server")
			stop()
		}()

		// 初始化数据库连接
		database.InitDB(ctx, &setting.Database)

		// Start HTTP
		go func(ctx context.Context) {
			run.Task(ctx)
		}(ctx)

		api.Run(ctx, &setting.System)
	},
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "f", "conf/config.yaml", "config file")
	rootCmd.AddCommand(versionCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of " + name,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("Version:   %s\n", common.Version)
		fmt.Printf("CommitID:  %s\n", common.CommitID)
		fmt.Printf("BuildTime: %s\n", common.BuildTime)
		fmt.Printf("GoVersion: %s\n", common.GoVersion)
		fmt.Printf("BuildUser: %s\n", common.BuildUser)
	},
}
