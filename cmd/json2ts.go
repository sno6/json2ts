package cmd

import (
	"log"

	"github.com/sno6/json2ts/parse"
	"github.com/sno6/json2ts/transform"

	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:   "github.com/sno6/json2ts",
		Short: "Transform JSON into typescript classes",
		Run: func(cmd *cobra.Command, args []string) {
			input := cmd.Flag("input").Value.String()

			nodes, err := (parse.Parser{}).Parse(input)
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
