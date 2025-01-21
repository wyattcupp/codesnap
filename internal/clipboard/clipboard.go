package clipboard

import (
	"errors"
	"os/exec"
	"runtime"
)

// WriteClipboard copies the given content to the system clipboard,
// supporting macOS, Linux (xclip or wl-copy), and Windows (clip).
func WriteClipboard(content string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		in, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		_, _ = in.Write([]byte(content))
		_ = in.Close()
		return cmd.Wait()
	case "windows":
		cmd := exec.Command("clip")
		in, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		_, _ = in.Write([]byte(content))
		_ = in.Close()
		return cmd.Wait()
	default:
		cmd := exec.Command("xclip", "-selection", "clipboard")
		in, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			// fallback to wl-copy
			cmd2 := exec.Command("wl-copy")
			in2, err2 := cmd2.StdinPipe()
			if err2 != nil {
				return err2
			}
			if err2 := cmd2.Start(); err2 != nil {
				return errors.New("failed to use xclip or wl-copy for clipboard copying")
			}
			_, _ = in2.Write([]byte(content))
			_ = in2.Close()
			return cmd2.Wait()
		}
		_, _ = in.Write([]byte(content))
		_ = in.Close()
		return cmd.Wait()
	}
}
