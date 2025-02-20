package commands

import (
	"os"

	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

// Version 命令版本号
var (
	version     = "v0.0.0"
	versionFile = os.TempDir() + "/sponge/.github/version"
)

// NewRootCMD 命令入口
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "sponge",
		Long:          "sponge management tools",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       getVersion(),
	}

	cmd.AddCommand(
		generate.ModelCommand(),
		generate.DaoCommand(),
		generate.HandlerCommand(),
		generate.HTTPCommand(),
		generate.ProtoCommand(),
		generate.ServiceCommand(),
		generate.GRPCCommand(),
		generate.ConfigCommand(),
		UpdateCommand(),
	)
	return cmd
}

func getVersion() string {
	data, _ := os.ReadFile(versionFile)
	v := string(data)
	if v != "" {
		return v
	}
	return version
}
