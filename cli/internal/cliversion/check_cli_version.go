package cliversion

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Checks if this cli version is out of date
func CheckCliVersion(verbose bool) error {

	cliVersion := &struct {
		CompatibleVersion string `json:"compatible_version"`
	}{}

	if err := api.GetJSON(config.CliVersionCheckEndpoint, cliVersion, "", nil); err != nil {
		if verbose {
			return fmt.Errorf("version check failed: %w", err)
		}
		return errors.New("unable to reach HashDrop server (use --verbose for details)")
	}

	// compare the version recevied from our version
	if cliVersion.CompatibleVersion != config.CurrentCliVersion {
		fmt.Println("Update Your Hashdrop CLI.")
		fmt.Printf("Current version: %s\n", config.CurrentCliVersion)
		fmt.Printf("Latest  version: %s\n\n", cliVersion.CompatibleVersion)
		fmt.Println("To update, run:")
		fmt.Println("go install github.com/anxhukumar/hashdrop/cli/cmd/hashdrop@latest")

		return fmt.Errorf("cli version outdated")
	}

	return nil
}
