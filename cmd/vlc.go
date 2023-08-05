package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"packer/lib"
	"path/filepath"
	"strings"
)

var vlcCmd = &cobra.Command{
	Use:   "vlc",
	Short: "Pack file using variable-length code",
	Run:   pack,
}

const (
	packedExtension  = "vlc"
	pathNotSpecified = "path to file is not specified"
)

func pack(cmd *cobra.Command, args []string) {
	if len(args) == 0 || args[0] == "" {
		handleErr(errors.New(pathNotSpecified))
	}

	filePath := args[0]

	r, err := os.Open(filePath)
	if err != nil {
		handleErr(err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		handleErr(err)
	}

	packed := lib.Encode(string(data))

	err = os.WriteFile(packedFileName(filePath), []byte(packed), 0644)
}

func packedFileName(path string) string {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)
	return baseName + "." + packedExtension
}

func init() {
	packCmd.AddCommand(vlcCmd)
}
