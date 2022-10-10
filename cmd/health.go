package cmd

import (
	"context"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Checks the health or readiness status of the Prometheus server",
	Long: `Checks the health or readiness status of the Prometheus server.
	Health: Checks the health of an endpoint, which always returns 200 and should be used to check Prometheus health.
    Ready: Checks the readiness of an endpoint, which returns 200 when Prometheus is ready to serve traffic (i.e. respond to queries).`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			rc         int
			output     string
			statuscode int
		)

		c := cliConfig.Client()
		err := c.Connect()
		if err != nil {
			check.ExitError(err)
		}

		// Ready status
		if cliConfig.PReady {
			ready, err := c.Ready()
			if err != nil {
				check.ExitError(err)
			}

			statuscode = ready.StatusCode

			rc, output, err = c.GetStatus(ready)
			if err != nil {
				check.ExitError(err)
			}
		} else {
			// Health status
			health, err := c.Health()
			if err != nil {
				check.ExitError(err)
			}

			statuscode = health.StatusCode

			rc, output, err = c.GetStatus(health)
			if err != nil {
				check.ExitError(err)
			}
		}

		if cliConfig.Info {
			// Displays various build information properties about the Prometheus server
			info, err := c.Api.Buildinfo(context.Background())
			if err != nil {
				check.ExitError(err)
			}

			output += "\n\n" +
				"Version: " + info.Version + "\n" +
				"Branch: " + info.Branch + "\n" +
				"BuildDate: " + info.BuildDate + "\n" +
				"BuildUser: " + info.BuildUser + "\n" +
				"Revision: " + info.Revision
		}

		// Statuscode 200 && "Prometheus Server is Healthy." -> 0 OK
		// Statuscode 200 && "Prometheus Server is Ready." -> 0 OK

		p := perfdata.PerfdataList{
			{Label: "status", Value: rc},
			{Label: "output", Value: output},
			{Label: "statuscode", Value: statuscode},
		}

		check.ExitRaw(rc, output, "|", p.String())
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)

	fs := healthCmd.Flags()
	fs.BoolVarP(&cliConfig.PReady, "ready", "r", false,
		"Checks the readiness of an endpoint")
	fs.BoolVarP(&cliConfig.Info, "info", "i", false,
		"Displays various build information properties about the Prometheus server")

	fs.SortFlags = false
	healthCmd.DisableFlagsInUseLine = true
}