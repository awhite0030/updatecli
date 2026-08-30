package version

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	sv "github.com/Masterminds/semver/v3"
)

var (
	// Version contains application version
	Version string

	// BuildTime contains application build time
	BuildTime string

	// GoVersion contains the golang version uses to build this binary
	GoVersion string

	// DisableDevWarning is used to identify if we already notify that we use a dev version
	isDevWarningDisabled bool
)

// Show displays various version information
func Show() {
	goVer := strings.TrimPrefix(runtime.Version(), "go")
	logrus.Infof("")
	logrus.Infof("Application:\t%s", Version)
	if GoVersion != "" {
		baseGo := strings.TrimSpace(strings.ReplaceAll(GoVersion, "go version go", ""))
		// baseGo might already contain the os/arch, since the old build script was setting GoVersion to $(go version)
		// if it does, don't append it again
		archString := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
		if strings.Contains(baseGo, "/") {
			// Just print it as it is
			logrus.Infof("Golang     :\t%s", baseGo)
		} else {
			logrus.Infof("Golang     :\t%s %s", baseGo, archString)
		}
	} else {
		logrus.Infof("Golang     :\t%s %s/%s", goVer, runtime.GOOS, runtime.GOARCH)
	}
	logrus.Infof("Build Time :\t%s", BuildTime)
	logrus.Infof("")
}

// IsGreaterThan test if an updatecli manifest required version is greater or equal to the current updatecli binary version
// Please not that empty version are set to 0.0.0
func IsGreaterThan(binaryVersion, manifestVersion string) (bool, error) {
	if len(manifestVersion) == 0 {
		manifestVersion = "0.0.0"
	}

	if len(binaryVersion) == 0 {
		if !isDevWarningDisabled {
			logrus.Warningf("Updatecli binary version is unset. This means you are using a development version that ignores manifest version constraint.\n")
			isDevWarningDisabled = true
		}

		return true, nil
	}

	mv, err := sv.NewVersion(manifestVersion)
	if err != nil {
		return false, fmt.Errorf("can't parse Updatecli manifest version %q - %q", manifestVersion, err)
	}

	bv, err := sv.NewVersion(binaryVersion)
	if err != nil {
		return false, fmt.Errorf("can't parse Updatecli binary version %q - %q", binaryVersion, err)
	}

	return bv.GreaterThan(mv) || bv.Equal(mv), nil
}
