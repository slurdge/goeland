FOR /F "tokens=* USEBACKQ" %%F IN (`git rev-parse HEAD`) DO (SET GIT_COMMIT=%%F)
FOR /F "tokens=* USEBACKQ" %%F IN (`powershell get-date -format "{yyyyMMddHHmmss}"`) DO (SET BUILD_DATE=%%F)
gox.py -parallel 4 -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags "-X github.com/slurdge/goeland/version.GitCommit=%GIT_COMMIT% -X github.com/slurdge/goeland/version.BuildDate=%BUILD_DATE% -s -w"
