package firebase_test

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
)

const FirestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"
const StartMsg = "Dev App Server is now running"

func TestMain(m *testing.M) {
	cmd, stderr := startFirestoreEmulator()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		monitorEmulatorOutput(stderr)
		defer wg.Done()
	}()

	wg.Wait()

	// Now it's running, we can run our unit tests
	result := m.Run()

	defer killCommand(cmd, result)
	defer stderr.Close()
}

func startFirestoreEmulator() (*exec.Cmd, io.ReadCloser) {
	cmd := exec.Command("gcloud", "emulators", "firestore", "start", "--host-port=localhost:9090", "--quiet")

	// this makes the emulator killable
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Failed to get stderr pipe: %v", err)
	}
	// No need to Close() stderr, it will be closed when the cmd is killed

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start emulator: %v", err)
	}

	return cmd, stderr
}

func monitorEmulatorOutput(r io.ReadCloser) {
	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Error reading stderr: %v", err)
			}
			break
		}

		if n > 0 {
			d := string(buf[:n])
			log.Printf("%s", d)

			if strings.Contains(d, StartMsg) {
				return
			}

			pos := strings.Index(d, FirestoreEmulatorHost+"=")
			if pos > 0 {
				host := d[pos+len(FirestoreEmulatorHost)+1 : len(d)-1]
				os.Setenv(FirestoreEmulatorHost, host)
			}
		}
	}
}

func killCommand(cmd *exec.Cmd, result int) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	os.Exit(result)
}
