gox -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags "-X github.com/slurdge/goeland/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/slurdge/goeland/version.BuildDate=${BUILD_DATE}"
