package cmd

import (
	"runtime"

	"fmt"

	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/exec"
	"github.com/drud/ddev/pkg/util"
	"github.com/spf13/cobra"
)

var DdevLaunchCmd = &cobra.Command{
	Use:   "launch <appname>",
	Short: "Launches an app in the default browser.",
	Long:  `Launches an app in the default browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {

		// No arguments - try to launch the current app.
		case 0:
			app, err := ddevapp.GetActiveApp("")
			if err != nil {
				util.Failed("Failed to launch browser: %v", err)
			}

			launchBrowser(app)
			break

		// One argument - try to launch the requested app.
		case 1:
			appName := args[0]

			for _, a := range ddevapp.GetApps() {
				if a.Name == appName {
					if err := launchBrowser(a); err != nil {
						util.Failed("Failed to launch browser: %v", err)
					}

					return
				}
			}

			util.Failed("App not found: %s", appName)
			break

		// Other - print usage and exit.
		default:
			util.Error("Incorrect arguments provided. Please use 'ddev launch [appname]'")
			break
		}
	},
}

// getBrowserCommand will return the relevant command used to open the default browser
// in the current operating system.
func getBrowserCommand() (string, error) {
	goos := runtime.GOOS

	switch goos {
	case "linux":
		return "xdg-open", nil

	case "darwin":
		return "open", nil

	case "win32":
		return "start", nil

	default:
		return "", fmt.Errorf("unable to build launch command for OS: %s", goos)
	}
}

// launchBrowser will attempt to open the site in the default browser.
func launchBrowser(app *ddevapp.DdevApp) error {
	baseCommand, err := getBrowserCommand()
	if err != nil {
		return err
	}

	if _, err := exec.RunCommand(baseCommand, []string{app.GetHTTPURL()}); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCmd.AddCommand(DdevLaunchCmd)
}
