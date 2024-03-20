package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/kevinmcconnell/mproxy/internal/server"
)

type runCommand struct {
	cmd              *cobra.Command
	debugLogsEnabled bool
}

func newRunCommand() *runCommand {
	runCommand := &runCommand{}
	runCommand.cmd = &cobra.Command{
		Use:   "run",
		Short: "Run the server",
		RunE:  runCommand.run,
	}

	runCommand.cmd.Flags().BoolVar(&runCommand.debugLogsEnabled, "debug", false, "Include debugging logs")
	runCommand.cmd.Flags().IntVar(&globalConfig.HttpPort, "http-port", server.DefaultHttpPort, "Port to serve HTTP traffic on")
	runCommand.cmd.Flags().IntVar(&globalConfig.HttpsPort, "https-port", server.DefaultHttpsPort, "Port to serve HTTPS traffic on")
	runCommand.cmd.Flags().DurationVar(&globalConfig.HttpIdleTimeout, "http-idle-timeout", server.DefaultHttpIdleTimeout, "Timeout before idle connection is closed")
	runCommand.cmd.Flags().DurationVar(&globalConfig.HttpReadTimeout, "http-read-timeout", server.DefaultHttpReadTimeout, "Tiemout for client to send a request")
	runCommand.cmd.Flags().DurationVar(&globalConfig.HttpWriteTimeout, "http-write-timeout", server.DefaultHttpWriteTimeout, "Timeout for client to receive a response")

	return runCommand
}

func (c *runCommand) run(cmd *cobra.Command, args []string) error {
	c.setLogger()

	router := server.NewRouter(globalConfig.StatePath())
	router.RestoreLastSavedState()

	s := server.NewServer(&globalConfig, router)
	s.Start()
	defer s.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	return nil
}

func (c *runCommand) setLogger() {
	level := slog.LevelInfo
	if c.debugLogsEnabled {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}
