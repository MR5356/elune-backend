package fileutil

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Decoder interface {
	Decode(any) error
}

// UnMarshalToAnyFromFile unmarshal data from a file to a target value.
//
// It takes a file path as input and the target value to unmarshal to.
// The function returns an error if there was a problem opening or decoding the file.
func UnMarshalToAnyFromFile(filePath string, target any) error {
	file, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filePath)
	}

	if file.Size() == 0 {
		return nil
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("open file %s error: %+v", filePath, err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logrus.Warnf("close file %s error: %+v", filePath, err)
		}
	}(f)

	var decoder Decoder
	if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
		decoder = yaml.NewDecoder(f)
	} else if strings.HasSuffix(filePath, ".json") {
		decoder = json.NewDecoder(f)
	} else {
		return fmt.Errorf("unsupported file type %s", filePath)
	}

	err = decoder.Decode(target)
	return err
}
