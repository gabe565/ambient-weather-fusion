package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gabe565.com/ambient-weather-fusion/internal/ambientweather"
	"gabe565.com/ambient-weather-fusion/internal/config"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func New(options ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ambient-weather-fusion",
		Short: "Integrate consensus-based Ambient Weather readings into Home Assistant",
		RunE:  run,
		Args:  cobra.NoArgs,

		DisableAutoGenTag: true,
	}
	conf := config.New()
	conf.RegisterFlags(cmd)
	if cmd.Context() == nil {
		cmd.SetContext(context.Background())
	}
	cmd.SetContext(config.NewContext(cmd.Context(), conf))
	for _, option := range options {
		option(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		panic("command missing config")
	}

	if conf.Latitude == 0 || conf.Longitude == 0 || conf.Radius == 0 || conf.BaseTopic == "" {
		return cmd.Help()
	}

	cmd.SilenceUsage = true

	slog.Info("Ambient Weather Fusion", "version", cobrax.GetVersion(cmd), "commit", cobrax.GetCommit(cmd))

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	server := ambientweather.NewServer(conf,
		ambientweather.WithVersion(cobrax.GetVersion(cmd)),
		ambientweather.WithUserAgent(cobrax.BuildUserAgent(cmd)),
	)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Close(ctx); err != nil {
			slog.Error("Failed to clean up retained data", "error", err)
		}
	}()

	return server.Run(ctx)
}
