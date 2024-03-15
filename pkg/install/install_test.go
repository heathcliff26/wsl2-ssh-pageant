package install

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/test.zip
var embeddedTestFile []byte

func init() {
	embeddedWinExeZip = embeddedTestFile
	winExeFilepath = "./test.txt"
	winExeZipFilename = "test.txt"
}

func TestInstall(t *testing.T) {
	err := Install()

	if assert.Nil(t, err) {
		err = os.Remove(winExeFilepath)
		assert.Nil(t, err)
	}
}

func TestCheckWinExeIntegrity(t *testing.T) {
	tMatrix := []struct {
		Name, File string
		Result     bool
	}{
		{"Unmodified", "test.txt", true},
		{"Modified", "modified.txt", false},
	}
	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			winExeFilepath = "./testdata/" + tCase.File
			t.Cleanup(func() {
				winExeFilepath = "./test.txt"
			})

			res, err := CheckWinExeIntegrity()

			assert := assert.New(t)

			assert.Nil(err)
			assert.Equal(tCase.Result, res)
		})
	}
}
