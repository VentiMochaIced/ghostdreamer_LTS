<#
.SYNOPSIS
    Automated build and release script for ghostdreamer_LTS
.DESCRIPTION
    Injects dynamic versioning, commit hashes, and build dates via ldflags,
    compiles standalone binaries for Windows and Linux, and manages Git tags.
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory=$false)]
    [string]$Version = "v0.1.0-beta",

    [Parameter(Mandatory=$false)]
    [switch]$CreateGitTag
)

# 1. Ensure output directory exists
$OutputDir = "dist"
if (!(Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# 2. Gather Git metadata
$GitCommit = "unknown"
$BuildDate = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")

if (Get-Command git -ErrorAction SilentlyContinue) {
    try {
        $GitCommit = (git rev-parse --short HEAD).Trim()
    } catch {
        Write-Host "Warning: Could not resolve git commit hash." -ForegroundColor Yellow
    }
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " Building Ghostdreamer LTS Release" -ForegroundColor Cyan
Write-Host " Version:     $Version" -ForegroundColor Green
Write-Host " Git Commit:  $GitCommit" -ForegroundColor Green
Write-Host " Build Date:  $BuildDate" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan

# 3. Configure Linker Flags
$LdFlags = "-s -w -X main.AppVersion=$Version -X main.GitCommit=$GitCommit -X main.BuildDate=$BuildDate"

# 4. Compile for Windows (AMD64)
Write-Host "`n[1/2] Compiling Windows executable..." -ForegroundColor Yellow
$Env:CGO_ENABLED = 0
$Env:GOOS = "windows"
$Env:GOARCH = "amd64"
go build -ldflags $LdFlags -o "$OutputDir/ghostdreamerv1b-windows-amd64.exe" main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "-> Success: $OutputDir/ghostdreamerv1b-windows-amd64.exe" -ForegroundColor Green
} else {
    Write-Error "Windows compilation failed."
    exit 1
}

# 5. Compile for Linux (AMD64)
Write-Host "`n[2/2] Compiling Linux binary..." -ForegroundColor Yellow
$Env:GOOS = "linux"
$Env:GOARCH = "amd64"
go build -ldflags $LdFlags -o "$OutputDir/ghostdreamerv1b-linux-amd64" main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "-> Success: $OutputDir/ghostdreamerv1b-linux-amd64" -ForegroundColor Green
} else {
    Write-Error "Linux compilation failed."
    exit 1
}

# 6. Optional Git Tagging & Push
if ($CreateGitTag) {
    Write-Host "`n[+] Creating Git tag: $Version" -ForegroundColor Cyan
    git tag -a $Version -m "Release $Version"
    git push origin $Version
    Write-Host "-> Tag pushed to remote." -ForegroundColor Green
}

Write-Host "`nBuild complete! Binaries are located in the './$OutputDir' directory." -ForegroundColor Green