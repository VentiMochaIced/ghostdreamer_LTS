<#
.SYNOPSIS
    Quick-syncs local changes, increments version tags, and builds the Go binary.
#>
param(
    [string]$Version = "v0.1.1-beta",
    [string]$Message = "Routine LTS sync and build update"
)

Write-Host "--- 1. Pulling latest upstream/fork changes ---" -ForegroundColor Cyan
git pull origin main

Write-Host "--- 2. Staging and committing local changes ---" -ForegroundColor Cyan
git add .
git commit -m "$Message"
git push origin main

Write-Host "--- 3. Building isolated standalone binary ---" -ForegroundColor Cyan
$GitCommit = (git rev-parse --short HEAD).Trim()
$BuildDate = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
$LdFlags = "-s -w -X main.AppVersion=$Version -X main.GitCommit=$GitCommit -X main.BuildDate=$BuildDate"

$Env:CGO_ENABLED = 0
$Env:GOOS = "windows"
$Env:GOARCH = "amd64"
go build -ldflags $LdFlags -o "dist/ghostdreamerv1b-windows-amd64.exe" main.go

Write-Host "-> Build complete: dist/ghostdreamerv1b-windows-amd64.exe [Version: $Version]" -ForegroundColor Green