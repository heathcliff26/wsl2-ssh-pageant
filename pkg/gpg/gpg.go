package gpg

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strconv"
)

func handleGPG(path string) {
	var port int
	var nonce [16]byte

	file, err := os.Open(path)
	if err != nil {
		slog.Error("Failed to open gpg directory", "err", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(file)
	tmp, _, err := reader.ReadLine()
	if err != nil {
		slog.Error("Could not read port from gpg nonce", "err", err)
		return
	}
	port, err = strconv.Atoi(string(tmp))
	if err != nil {
		slog.Error("Could not read port from gpg nonce", "err", err)
		return
	}
	n, err := reader.Read(nonce[:])
	if err != nil {
		slog.Error("Could not read port from gpg nonce", "err", err)
		return
	}

	if n != 16 {
		slog.Error("Could not connet gpg: incorrect number of bytes for nonceRead incorrect number of bytes for nonce", slog.Int("bytes", n))
		return
	}

	gpgConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		slog.Error("Could not connect to gpg", "err", err)
		return
	}

	_, err = gpgConn.Write(nonce[:])
	if err != nil {
		slog.Error("Could not authenticate gpg", "err", err)
		return
	}

	go func() {
		_, err := io.Copy(gpgConn, os.Stdin)
		if err != nil {
			slog.Error("Could not copy gpg data from assuan socket to socket", "err", err)
			return
		}
	}()

	_, err = io.Copy(os.Stdout, gpgConn)
	if err != nil {
		slog.Error("Could not copy gpg data from socket to assuan socket", "err", err)
		return
	}
}
