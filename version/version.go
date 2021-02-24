package version

import (
	"bufio"
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// GitCommit returns the git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// ChangeLog returns the full changelog
var ChangeLog string

// Version returns the main version number that is being run at the moment.
// Given a version number MAJOR.MINOR.PATCH, increment the:
//     MAJOR version when you make incompatible API changes,
//     MINOR version when you add functionality in a backwards compatible manner, and
//     PATCH version when you make backwards compatible bug fixes.
// Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.
var Version = "dev"

// BuildDate returns the date the binary was built
var BuildDate = ""

// GoVersion returns the version of the go runtime used to compile the binary
var GoVersion = runtime.Version()

// OsArch returns the os and arch used to build the binary
var OsArch = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

// ExtractVersionFromChangelog is used to extract the latest version referenced in the Changelog, this is our version
func ExtractVersionFromChangelog(changeLog string) {
	ChangeLog = changeLog
	scanner := bufio.NewScanner(strings.NewReader(ChangeLog))
	versionRegExp := regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
	for scanner.Scan() {
		line := scanner.Bytes()
		if versionRegExp.Match(line) {
			Version = string(line)
			break
		}
	}
}
