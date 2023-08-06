package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"packer/lib/compression"
	"packer/lib/compression/vlc"
	"packer/lib/compression/vlc/table/shannon_fano"
	"path/filepath"
)

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack file",
	Run:   pack,
}

const (
	packedExtension  = "shannon"
	pathNotSpecified = "path to file is not specified"
)

func pack(cmd *cobra.Command, args []string) {
	if len(args) == 0 || args[0] == "" {
		handleErr(errors.New(pathNotSpecified))
	}

	method := cmd.Flag("method").Value.String()

	var encoder compression.Encoder

	switch method {
	case "shannon":
		encoder = vlc.New(shannon_fano.NewGenerator())
	default:
		panic("no such encoder: " + method)
	}

	filePath := args[0]

	r, err := os.Open(filePath)
	if err != nil {
		handleErr(err)
	}
	defer os.Remove(filePath) // remove old file
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		handleErr(err)
	}

	packed := encoder.Encode(string(data))

	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil {
		handleErr(err)
	}

}

func packedFileName(path string) string {
	fileName := filepath.Base(path)
	return fileName + "." + packedExtension
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringP("method", "m", "",
		"compression method: vlc")

	if err := packCmd.MarkFlagRequired("method"); err != nil {
		panic(err.Error())
	}
}
