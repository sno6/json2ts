package cmd

import (
	"json2ts/parse"
	"json2ts/transform"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:   "json2ts",
		Short: "Transform JSON into typescript classes",
		Run: func(cmd *cobra.Command, args []string) {
			nodes, err := (parse.Parser{}).Parse(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}

			err = (transform.Transformer{}).Transform(nodes, &transform.Config{
				BaseClassName:   cmd.Flag("root").Value.String(),
				PrefixClassName: cmd.Flag("prefix").Value.String(),
				Output:          cmd.Flag("output").Value.String(),
				Decorators:      cmd.Flag("decorators").Value.String() == "true",
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	root.Flags().StringP("input", "i", "", "Optional input file")
	root.Flags().StringP("output", "o", "", "Optional output file (defaults to stdout)")
	root.Flags().StringP("root", "r", "", "Name of the root class")
	root.Flags().StringP("prefix", "p", "", "Prefix for all class names")
	root.Flags().BoolP("decorators", "d", false, "Add decorators to class parameters")

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
