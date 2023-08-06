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
	"strings"
)

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack file",
	Run:   unpack,
}

func unpack(cmd *cobra.Command, args []string) {
	var decoder compression.Decoder

	if len(args) == 0 || args[0] == "" {
		handleErr(errors.New(pathNotSpecified))
	}

	method := cmd.Flag("method").Value.String()

	switch method {
	case "shannon":
		decoder = vlc.New(shannon_fano.NewGenerator())
	default:
		panic("no such dencoder: " + method)
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

	unpacked := decoder.Decode(data)

	err = os.WriteFile(unpackedFileName(filePath), []byte(unpacked), 0644)
	if err != nil {
		handleErr(err)
	}
}

func unpackedFileName(path string) string {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)
	return baseName
}

func init() {
	rootCmd.AddCommand(unpackCmd)

	unpackCmd.Flags().StringP("method", "m", "",
		"decompression method: vlc")

	if err := unpackCmd.MarkFlagRequired("method"); err != nil {
		panic(err.Error())
	}
}
