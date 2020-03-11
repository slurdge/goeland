package version

import (
	"fmt"
	"runtime"
)

// GitCommit returns the git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version returns the main version number that is being run at the moment.
// Given a version number MAJOR.MINOR.PATCH, increment the:
//     MAJOR version when you make incompatible API changes,
//     MINOR version when you add functionality in a backwards compatible manner, and
//     PATCH version when you make backwards compatible bug fixes.
// Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.
const Version = "0.4.0"

// BuildDate returns the date the binary was built
var BuildDate = ""

// GoVersion returns the version of the go runtime used to compile the binary
var GoVersion = runtime.Version()

// OsArch returns the os and arch used to build the binary
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
