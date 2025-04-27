package ssh

import (
	"bufio"
	"encoding/binary"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/Microsoft/go-winio"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

const (
	// Windows constats
	invalidHandleValue = ^windows.Handle(0)
	pageReadWrite      = 0x4
	fileMapWrite       = 0x2

	// ssh-agent/Pageant constants
	agentMaxMessageLength = 8192
	agentCopyDataID       = 0x804e50ba
)

var failureMessage = [...]byte{0, 0, 0, 1, 5}

// copyDataStruct is used to pass data in the WM_COPYDATA message.
// We directly pass a pointer to our copyDataStruct type, we need to be
// careful that it matches the Windows type exactly
type copyDataStruct struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

// #nosec G103: Using unsafe here to write to shared memory
func queryPageant(buf []byte) (result []byte, err error) {
	if len(buf) > agentMaxMessageLength {
		err = NewErrMessageLength(len(buf))
		return
	}

	// Static string, should not return any error
	pageantPtr, _ := syscall.UTF16PtrFromString("Pageant")
	hwnd := win.FindWindow(pageantPtr, pageantPtr)

	// Launch gpg-connect-agent
	if hwnd == 0 {
		slog.Info("Launching gpg-connect-agent")
		err = exec.Command("gpg-connect-agent", "/bye").Run()
		if err != nil {
			err = NewErrGPGAgentLaunch(err)
			return
		}
	}

	hwnd = win.FindWindow(pageantPtr, pageantPtr)
	if hwnd == 0 {
		err = NewErrNoPageant()
		return
	}

	// Adding process id in order to support parrallel requests.
	mapName := "WSLPageantRequest" + strconv.Itoa(os.Getpid())
	mapNamePtr, err := syscall.UTF16PtrFromString(mapName)
	if err != nil {
		return
	}

	fileMap, err := windows.CreateFileMapping(invalidHandleValue, nil, pageReadWrite, 0, agentMaxMessageLength, mapNamePtr)
	if err != nil {
		return
	}
	defer func() {
		_ = windows.CloseHandle(fileMap)
	}()

	sharedMemory, err := windows.MapViewOfFile(fileMap, fileMapWrite, 0, 0, 0)
	if err != nil {
		return
	}
	defer func() {
		_ = windows.UnmapViewOfFile(sharedMemory)
	}()

	// Ignore vet warning "possible misuse of unsafe.Pointer". Shared memory should not be visible to garbage collector
	sharedMemoryArray := unsafe.Slice((*byte)(unsafe.Pointer(sharedMemory)), agentMaxMessageLength)
	copy(sharedMemoryArray[:], buf)

	mapNameWithNul := mapName + "\000"

	// We use our knowledge of Go strings to get the length and pointer to the
	// data and the length directly
	cds := copyDataStruct{
		dwData: agentCopyDataID,
		cbData: uint32(len(mapNameWithNul)),
		lpData: uintptr(unsafe.Pointer(unsafe.StringData(mapNameWithNul))),
	}

	ret := win.SendMessage(hwnd, win.WM_COPYDATA, 0, uintptr(unsafe.Pointer(&cds)))
	if ret == 0 {
		err = NewErrWMCopydata()
		return
	}

	len := binary.BigEndian.Uint32(sharedMemoryArray[:4])
	len += 4

	if len > agentMaxMessageLength {
		err = NewErrMessageLength(int(len))
		return
	}

	result = make([]byte, len)
	copy(result, sharedMemoryArray[:len])

	return
}

func handleSSH() {
	reader := bufio.NewReader(os.Stdin)
	for {
		lenBuf := make([]byte, 4)
		_, err := io.ReadFull(reader, lenBuf)
		if err != nil {
			if err == io.EOF {
				slog.Debug("Got EOF error when reading message length from stdin")
			} else {
				slog.Error("Failed to read message length from stdin", "err", err)
			}
			return
		}

		len := binary.BigEndian.Uint32(lenBuf)
		slog.Info("Reading input", slog.Int("length", int(len)))
		buf := make([]byte, len)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			if err == io.EOF {
				slog.Debug("Got EOF error when reading message from stdin")
			} else {
				slog.Error("Failed to read message from stdin", "err", err)
			}
			return
		}

		slog.Info("Querying pageant")
		result, err := queryPageant(append(lenBuf, buf...))
		if err != nil {
			// If for some reason talking to Pageant fails we fall back to
			// sending an agent error to the client
			slog.Error("Failed to query pageant", "err", err)
			result = failureMessage[:]
		}

		_, err = os.Stdout.Write(result)
		if err != nil {
			slog.Error("Failed to write result to stdout", "err", err)
			return
		}
	}
}

func handlePipedSSH(pipe string) {
	conn, err := winio.DialPipe(pipe, nil)
	if err != nil {
		slog.Error("Failed to dial ssh pipe", slog.String("pipe", pipe), "err", err)
		return
	}
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil && err != io.EOF {
			slog.Error("Failed to copy from pipe to stdout", "err", err)
			os.Exit(1)
		}
	}()
	_, err = io.Copy(conn, os.Stdin)
	if err != nil && err != io.EOF {
		slog.Error("Failed to copy from stdin to pipe", "err", err)
		os.Exit(1)
	}
}
