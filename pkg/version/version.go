package version

import (
	"fmt"
	"strings"
)

var (
	Product   = "control-plane"
	version   = "unknown"
	gitTag    = "unknown"
	gitCommit = "unknown"
	buildDate = "unknown"
)

type BuildInfo struct {
	Product   string
	Version   string
	GitTag    string
	GitCommit string
	BuildDate string
}

func (b BuildInfo) FormatDetailedProductInfo() string {
	base := []string{
		fmt.Sprintf("Product:       %s", b.Product),
		fmt.Sprintf("Version:       %s", b.Version),
		fmt.Sprintf("Git Tag:       %s", b.GitTag),
		fmt.Sprintf("Git Commit:    %s", b.GitCommit),
		fmt.Sprintf("Build Date:    %s", b.BuildDate),
	}
	return strings.Join(
		base,
		"\n",
	)
}

var Build BuildInfo

func init() {
	Build = BuildInfo{
		Product:   Product,
		Version:   version,
		GitTag:    gitTag,
		GitCommit: gitCommit,
		BuildDate: buildDate,
	}
}
