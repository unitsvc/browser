package browser

import (
	"errors"
	"fmt"
	"os/exec"
)

func openBrowser(url string) error {
	err := runCmd("xdg-open", url)
	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("xdg-open: command not found - install xdg-utils from pkgsrc(7)")
	}
	return err
}
