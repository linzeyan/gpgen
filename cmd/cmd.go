package cmd

import (
	"github.com/linzeyan/gpgen/src"
	"github.com/spf13/cobra"
)

func root() *cobra.Command {
	var (
		graphql    bool
		protobuf   bool
		srcFile    string
		dstFile    string
		structName string
	)
	var cmd = &cobra.Command{
		Use:   "gpgen",
		Short: "Generate graphql or protobuf schema from model file.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !graphql && !protobuf {
				return cmd.Help()
			}
			src.Analysis(&src.Options{
				Graphql:    graphql,
				Protobuf:   protobuf,
				SrcFile:    srcFile,
				DstFile:    dstFile,
				StructName: structName,
			})
			return nil
		},
	}

	cmd.Flags().BoolVarP(&graphql, "graphql", "g", false, "Generate GraphQL schema.")
	cmd.Flags().BoolVarP(&protobuf, "protobuf", "p", false, "Generate protobuf schema.")
	cmd.Flags().StringVarP(&srcFile, "src", "s", "", "Specify source file path.")
	cmd.Flags().StringVarP(&dstFile, "dst", "d", "", "Specify destination file path.")
	cmd.Flags().StringVarP(&structName, "struct", "t", "", "Specify struct name.")
	return cmd
}

func Run() *cobra.Command {
	return root()
}
