@echo off
FOR /F "tokens=* USEBACKQ" %%F IN (`git rev-parse HEAD`) DO (SET GIT_COMMIT=%%F)
FOR /F "tokens=* USEBACKQ" %%F IN (`powershell get-date -format "{yyyyMMddHHmmss}"`) DO (SET BUILD_DATE=%%F)
FOR /F "tokens=* USEBACKQ" %%F IN (`type clientid.txt`) DO (SET CLIENT_ID=%%F)
del /s /q build\*.*
python gox.py -osarch="linux/386 linux/amd64 linux/arm64 linux/arm darwin/amd64 darwin/arm64 windows/amd64 windows/386" -parallel 4 -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags "-X github.com/slurdge/goeland/goeland.clientID=%CLIENT_ID% -X github.com/slurdge/goeland/version.GitCommit=%GIT_COMMIT% -X github.com/slurdge/goeland/version.BuildDate=%BUILD_DATE% -s -w"
