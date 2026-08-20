package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigFlagName = "config"

func BindFlags(rootCMD *cobra.Command) {
	rootCMD.PersistentFlags().String(ConfigFlagName, "", "Path to env configuration file")
	_ = viper.BindPFlag(ConfigFlagName, rootCMD.PersistentFlags().Lookup(ConfigFlagName))
}
