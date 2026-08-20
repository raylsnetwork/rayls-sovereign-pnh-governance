package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/flags"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/withstack"
)

var (
	rootCMD = &cobra.Command{
		Use: "",
	}

	runCMD = &cobra.Command{
		Use:   "run",
		Short: "Run API",
		Long:  "Run API",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Run(viper.GetString(flags.ConfigFlagName)); err != nil {
				return withstack.Wrap(err)
			}
			return nil
		},
	}
)

func init() {
	flags.BindFlags(runCMD)
}

func Execute() {
	rootCMD.AddCommand(runCMD)
	if err := rootCMD.Execute(); err != nil {
		logger.Error("failed to execute root cmd", "error", err)
	}
}
