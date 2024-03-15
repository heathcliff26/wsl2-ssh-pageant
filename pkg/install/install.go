package install

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"io"
	"os"
)

const (
	WinExeFilepath    = "/mnt/c/ProgramData/heathcliff26/wsl-ssh-pageant/wsl-ssh-pageant-windows.exe"
	WinExeZipFilename = "wsl-ssh-pageant-windows.exe"
	WinExeChecksum    = "0a661b7587859c06019a582c6ad7ca12f25e5baf5f42c0135567d75f51bfa4ec"
)

var (
	winExeFilepath    = WinExeFilepath
	winExeZipFilename = WinExeZipFilename
)

//go:embed embed.zip
var embeddedWinExeZip []byte

func CheckWinExeIntegrity() (bool, error) {
	data, err := os.ReadFile(winExeFilepath)
	if err != nil {
		return false, err
	}

	h := sha256.New()
	_, err = h.Write(data)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString(h.Sum(nil)) == WinExeChecksum, nil
}

func Install() error {
	data, err := unzipEmbeddedFile()
	if err != nil {
		return err
	}

	err = os.WriteFile(winExeFilepath, data, 0755)
	if err != nil {
		return err
	}

	return nil
}

// Unzip the embedded file and return the data
func unzipEmbeddedFile() ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(embeddedWinExeZip), int64(len(embeddedWinExeZip)))
	if err != nil {
		return nil, err
	}
	f, err := r.Open(winExeZipFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
