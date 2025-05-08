package testhelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/nsf/jsondiff"
	"github.com/pkg/errors"
)

// FileToStruct reads a json file to a struct
// As an additive utility a reader for the file bytes is returned to the caller
func FileToStruct(t *testing.T, filepath string, s interface{}) io.Reader {
	bb, err := os.ReadFile(filepath)
	if err != nil {
		t.Errorf(errors.Wrapf(err, "\033[93m FileToStruct: unable to read filepath %v \033[0m ", filepath).Error())
		t.FailNow()
	}

	err = json.Unmarshal(bb, s)
	if err != nil {
		t.Errorf(errors.Wrapf(err, "\033[93m FileToStruct: unable to unmarshal val=%s \033[0m ", bb).Error())
		t.FailNow()
	}

	return bytes.NewReader(bb)
}

// AssertJSONEquals asserts a struct and a json file are equal and pretty prints the difference if not equal
func AssertJSONEquals(t *testing.T, resFilePath string, actual interface{}) {
	actualStr, err := json.Marshal(actual)
	if err != nil {
		log.Fatalf("\033[93m Error: %s \033[0m ", err)
	}
	expected, _ := os.ReadFile(resFilePath)
	opts := jsondiff.DefaultConsoleOptions()
	diff, diffStr := jsondiff.Compare(expected, actualStr, &opts)
	if diff != jsondiff.FullMatch {
		t.Errorf(fmt.Sprintf("Diff=%s", diffStr))
		t.FailNow()
	}
}

type runFunc func(name string, f func(t *testing.T))

// Run replaces t.Run for panic and stopping test runtime to start debugging when a sub test fails
func Run(t *testing.T) runFunc {
	return func(name string, f func(t *testing.T)) {
		if pass := t.Run(name, f); !pass {
			t.Errorf("Subtest is not passing - ( %v )", name)
			t.FailNow()
		}
	}
}
