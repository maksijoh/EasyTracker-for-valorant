$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path dist | Out-Null

go run github.com/akavel/rsrc@v0.10.2 `
    -manifest app.manifest `
    -ico assets/easytracker.ico `
    -o rsrc_windows_amd64.syso
if ($LASTEXITCODE -ne 0) { throw "Failed to generate Windows resources" }

go build `
    -trimpath `
    -ldflags "-s -w -H windowsgui -X main.version=1.0.0" `
    -o dist/EasyTracker-v1.0.0-windows-amd64.exe .
if ($LASTEXITCODE -ne 0) { throw "Failed to build EasyTracker" }

Write-Host "Built dist/EasyTracker-v1.0.0-windows-amd64.exe"
