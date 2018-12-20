package cmd

import (
	"fmt"
	"runtime"

	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/exec"
	"github.com/drud/ddev/pkg/util"
	"github.com/spf13/cobra"
)

var (
	httpsArg bool
)

var OpenCmd = &cobra.Command{
	Use:   "open [project]",
	Short: "Opens a running project in the default browser",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var app *ddevapp.DdevApp

		if len(args) > 1 {
			util.Failed("Accepts one optional arg: project name")
		}

		if len(args) == 0 {
			app, err = ddevapp.GetActiveApp("")
			if err != nil {
				util.Failed("Failed to get current project: %v", err)
			}
		}

		if len(args) == 1 {
			app, err = ddevapp.GetActiveApp(args[0])
			if err != nil {
				util.Failed("Failed to get %s: %v", args[0], err)
			}
		}

		if app.SiteStatus() != ddevapp.SiteRunning {
			util.Failed("Project %s is not running", app.Name)
		}

		browserCmd, err := getBrowserCommand()
		if err != nil {
			util.Failed(err.Error())
		}

		var url string
		if httpsArg {
			url = app.GetHTTPSURL()
		} else {
			url = app.GetHTTPURL()
		}

		if _, err := exec.RunCommandPipe(browserCmd, []string{url}); err != nil {
			util.Failed("failed to open %s: %v", app.Name, err)
		}
	},
}

// getBrowserCommand will return the appropriate command to open a URL in the default browser for the OS
func getBrowserCommand() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return "xdg-open", nil

	case "darwin":
		return "open", nil

	case "windows":
		return "start", nil
	}

	return "", fmt.Errorf("no browser command for %s", runtime.GOOS)
}

func init() {
	RootCmd.AddCommand(OpenCmd)

	fs := OpenCmd.Flags()
	fs.BoolVarP(&httpsArg, "https", "s", false, "Open the https URL (default http)")
}
