package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"

	"github.com/spf13/cobra"
)

// ProtoCommand generate protobuf code
func ProtoCommand() *cobra.Command {
	var (
		moduleName string // go.mod文件的module名称
		serverName string // 服务名称
		outPath    string // 输出目录
		sqlArgs    = sql2code.Args{}
	)

	cmd := &cobra.Command{
		Use:   "proto",
		Short: "Generate protobuf code",
		Long: `generate protobuf code.

Examples:
  # generate protobuf code.
  sponge proto --module-name=yourModuleName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate protobuf code and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge proto --module-name=yourModuleName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			return runGenProtoCommand(moduleName, serverName, codes, outPath)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "p", "", "module-name is the name of the module in the 'go.mod' file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./proto_<time>")

	return cmd
}

func runGenProtoCommand(moduleName string, serverName string, codes map[string]string, outPath string) error {
	subTplName := "proto"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	if serverName == "" {
		serverName = moduleName
	}

	// 设置模板信息
	subDirs := []string{"api/serverNameExample"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{}                     // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"userExample.pb.go", "userExample.pb.validate.go",
		"userExample_grpc.pb.go"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addProtoFields(moduleName, serverName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, out = %s\n\n", subTplName, r.GetOutputDir())
	return nil
}

func addProtoFields(moduleName string, serverName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, protoFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // 替换v1/userExample.proto文件内容
			Old: protoFileMark,
			New: codes[parser.CodeTypeProto],
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: moduleName,
		},
		// 替换目录名称
		{
			Old: strings.Join([]string{"api", "serverNameExample", "v1"}, gofile.GetPathDelimiter()),
			New: strings.Join([]string{"api", serverName, "v1"}, gofile.GetPathDelimiter()),
		},
		{
			Old: "api/serverNameExample/v1",
			New: fmt.Sprintf("api/%s/v1", serverName),
		},
		{
			Old: "api.serverNameExample.v1",
			New: fmt.Sprintf("api.%s.v1", strings.ReplaceAll(serverName, "-", "_")), // proto package 不能存在"-"号
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
