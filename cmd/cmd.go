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
	"gabe565.com/ambient-weather-fusion/internal/mqtt"
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

	if conf.Latitude == 0 || conf.Longitude == 0 || conf.Radius == 0 || conf.TopicPrefix == "" {
		return cmd.Help()
	}

	cmd.SilenceUsage = true

	slog.Info("Ambient Weather Fusion", "version", cobrax.GetVersion(cmd), "commit", cobrax.GetCommit(cmd))

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	client, err := mqtt.Connect(ctx, conf)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := ambientweather.Cleanup(ctx, conf, client); err != nil {
			slog.Error("Failed to clean up ambient-weather payloads", "error", err)
		}

		if err := mqtt.Disconnect(ctx, conf, client); err != nil {
			slog.Error("Failed to disconnect from mqtt", "error", err)
		}
	}()

	if err := mqtt.PublishDiscovery(ctx, cmd, conf, client); err != nil {
		return err
	}

	ticker := time.NewTicker(1)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			ticker.Reset(5 * time.Minute)
			if err := ambientweather.Process(ctx, cmd, conf, client); err != nil {
				slog.Error("Failed to process ambient-weather data", "error", err)
			}
		}
	}
}
