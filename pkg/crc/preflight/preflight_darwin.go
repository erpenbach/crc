package preflight

import (
	"fmt"
	cmdConfig "github.com/code-ready/crc/cmd/crc/cmd/config"
	"github.com/code-ready/crc/pkg/crc/config"
)

// StartPreflightChecks performs the preflight checks before starting the cluster
func StartPreflightChecks() {
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckRootUser.Name),
		checkIfRunningAsNormalUser,
		"Checking if running as non-root",
		config.GetBool(cmdConfig.WarnCheckRootUser.Name),
	)
	preflightCheckSucceedsOrFails(false,
		checkOcBinaryCached,
		"Checking if oc binary is cached",
		false,
	)

	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckHyperKitInstalled.Name),
		checkHyperKitInstalled,
		"Checking if HyperKit is installed",
		config.GetBool(cmdConfig.WarnCheckHyperKitInstalled.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckHyperKitDriver.Name),
		checkMachineDriverHyperKitInstalled,
		"Checking if crc-driver-hyperkit is installed",
		config.GetBool(cmdConfig.WarnCheckHyperKitDriver.Name),
	)

	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckHostsFilePermissions.Name),
		checkHostsFilePermissions,
		fmt.Sprintf("Checking file permissions for %s", resolverFile),
		config.GetBool(cmdConfig.WarnCheckHostsFilePermissions.Name),
	)

	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckHostsFilePermissions.Name),
		checkHostsFilePermissions,
		fmt.Sprintf("Checking file permissions for %s", hostFile),
		config.GetBool(cmdConfig.WarnCheckHostsFilePermissions.Name),
	)
}

// SetupHost performs the prerequisite checks and setups the host to run the cluster
func SetupHost() {
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckRootUser.Name),
		checkIfRunningAsNormalUser,
		fixRunAsNormalUser,
		"Checking if running as non-root",
		config.GetBool(cmdConfig.WarnCheckRootUser.Name),
	)
	preflightCheckAndFix(false,
		checkOcBinaryCached,
		fixOcBinaryCached,
		"Caching oc binary",
		false,
	)

	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckHyperKitInstalled.Name),
		checkHyperKitInstalled,
		fixHyperKitInstallation,
		"Setting up virtualization with HyperKit",
		config.GetBool(cmdConfig.WarnCheckHyperKitInstalled.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckHyperKitDriver.Name),
		checkMachineDriverHyperKitInstalled,
		fixMachineDriverHyperKitInstalled,
		"Installing crc-machine-hyperkit",
		config.GetBool(cmdConfig.WarnCheckHyperKitDriver.Name),
	)

	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckResolverFilePermissions.Name),
		checkResolverFilePermissions,
		fixResolverFilePermissions,
		fmt.Sprintf("Setting file permissions for %s", resolverFile),
		config.GetBool(cmdConfig.WarnCheckResolverFilePermissions.Name),
	)

	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckHostsFilePermissions.Name),
		checkHostsFilePermissions,
		fixHostsFilePermissions,
		fmt.Sprintf("Setting file permissions for %s", hostFile),
		config.GetBool(cmdConfig.WarnCheckHostsFilePermissions.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckBundleCached.Name),
		checkBundleCached,
		fixBundleCached,
		"Unpacking bundle from the CRC binary",
		config.GetBool(cmdConfig.WarnCheckBundleCached.Name),
	)
}
