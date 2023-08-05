package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
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

	data, err := ioutil.ReadAll(r)
	if err != nil {
		handleErr(err)
	}

	packed := "" + string(data) // TO REMOVE

	err = ioutil.WriteFile(packedFileName(filePath), []byte(packed), 0644)
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
