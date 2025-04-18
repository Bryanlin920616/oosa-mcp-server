package main

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Bryanlin920616/oosa-mcp-server/config"
	iolog "github.com/Bryanlin920616/oosa-mcp-server/pkg/log"
	"github.com/Bryanlin920616/oosa-mcp-server/pkg/oosa"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// ServerTransport 定義服務器的傳輸方式
type ServerTransport string

const (
	// ServerTransportStdio 使用標準輸入輸出進行通信
	ServerTransportStdio ServerTransport = "stdio"
	// ServerTransportSSE 使用 SSE 進行通信
	ServerTransportSSE ServerTransport = "sse"
)

var (
	rootCmd = &cobra.Command{
		Use:     "server",
		Short:   "OOSA MCP Server",
		Long:    `A OOSA MCP server that handles various tools and resources.`,
		Version: fmt.Sprintf("%s (%s) %s", config.Version, config.Commit, config.Date),
	}

	serverCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  `Start the server in the specified mode (stdio or sse)`,
		Run: func(_ *cobra.Command, _ []string) {
			logFile := viper.GetString("log-file")
			logger, err := initLogger(logFile)
			if err != nil {
				stdlog.Fatal("Failed to initialize logger:", err)
			}

			cfg := runConfig{
				logger:      logger,
				logCommands: viper.GetBool("enable-command-logging"),
				transport:   ServerTransport(viper.GetString("server.transport")),
				addr:        viper.GetString("server.addr"),
				baseURL:     viper.GetString("server.base_url"),
			}

			if err := runServer(cfg); err != nil {
				stdlog.Fatal("failed to run server:", err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// 設定配置檔路徑
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oosa-mcp-server.yaml)")

	// 設定 server 相關的 flag
	serverCmd.Flags().StringP("transport", "t", "stdio", "服務器傳輸方式 (stdio 或 sse)")
	serverCmd.Flags().StringP("addr", "a", "0.0.0.0:8080", "SSE 服務器監聽地址")
	serverCmd.Flags().StringP("base-url", "b", "http://localhost:8080", "SSE 服務器 base URL（用於 origin 驗證）")
	serverCmd.Flags().String("log-file", "", "Path to log file")
	serverCmd.Flags().Bool("enable-command-logging", false, "When enabled, the server will log all command requests and responses")

	// 綁定 flag 到 viper
	_ = viper.BindPFlag("server.mode", serverCmd.Flags().Lookup("transport"))
	_ = viper.BindPFlag("server.addr", serverCmd.Flags().Lookup("addr"))
	_ = viper.BindPFlag("server.base_url", serverCmd.Flags().Lookup("base-url"))
	_ = viper.BindPFlag("log-file", serverCmd.Flags().Lookup("log-file"))
	_ = viper.BindPFlag("enable-command-logging", serverCmd.Flags().Lookup("enable-command-logging"))

	rootCmd.AddCommand(serverCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name ".oosa-mcp-server" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".oosa-mcp-server")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initLogger(outPath string) (*log.Logger, error) {
	logger := log.New()

	// 設定日誌級別
	level := viper.GetString("log.level")
	if level != "" {
		logLevel, err := log.ParseLevel(level)
		if err != nil {
			return nil, fmt.Errorf("invalid log level: %w", err)
		}
		logger.SetLevel(logLevel)
	} else {
		logger.SetLevel(log.DebugLevel)
	}

	// 設定日誌輸出
	target := viper.GetString("log.target")
	switch target {
	case "os":
		logger.SetOutput(os.Stderr)
		return logger, nil
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger.SetLevel(log.DebugLevel)
	logger.SetOutput(file)

	return logger, nil
}

type runConfig struct {
	logger      *log.Logger
	logCommands bool
	transport   ServerTransport
	addr        string
	baseURL     string
}

func runServer(cfg runConfig) error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create client like github, mongodb, etc.
	client := &oosa.Client{} // TODO: use real client

	// Create server
	mcpServer := oosa.NewServer(client, config.Version)

	// Create error logger
	stdLogger := stdlog.New(cfg.logger.Writer(), "server", 0)

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		switch cfg.transport {
		case ServerTransportStdio:
			// 建立 stdio server
			stdioServer := server.NewStdioServer(mcpServer)
			stdioServer.SetErrorLogger(stdLogger)
			cfg.logger.Info("Starting server in stdio mode")

			// 設定輸入輸出
			in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)
			if cfg.logCommands {
				// TODO: Add command logging
				loggedIO := iolog.NewIOLogger(in, out, cfg.logger)
				in, out = loggedIO, loggedIO
			}

			errC <- stdioServer.Listen(ctx, in, out)

		case ServerTransportSSE:
			// 建立 SSE server
			cfg.logger.Infof("Set base URL: %s", cfg.baseURL)
			sseServer := server.NewSSEServer(mcpServer,
				server.WithBaseURL(cfg.baseURL),
			)
			cfg.logger.Infof("Starting server in SSE mode on %s", cfg.addr)

			errC <- sseServer.Start(cfg.addr) // server listen 的地址

		default:
			errC <- fmt.Errorf("unsupported server transport: %s", cfg.transport)
		}
	}()

	// Output server info
	_, _ = fmt.Fprintf(os.Stderr, "OOSA MCP Server running in %s mode\n", cfg.transport)

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		cfg.logger.Infof("shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
