# dwm installer for Windows
# Usage: iex (irm https://raw.githubusercontent.com/jinkp/dbeaver-go-mcp/main/install.ps1)

$ErrorActionPreference = 'Stop'

$Repo   = "jinkp/dbeaver-go-mcp"
$Binary = "dwm"

function Get-LatestVersion {
    $url = "https://api.github.com/repos/$Repo/releases/latest"
    $release = Invoke-RestMethod -Uri $url -Headers @{ "User-Agent" = "dwm-installer" }
    return $release.tag_name
}

function Get-Arch {
    if ([System.Environment]::Is64BitOperatingSystem) { return "x86_64" }
    return "x86"
}

$Version = Get-LatestVersion
$Arch    = Get-Arch
$Archive = "${Binary}_${Version}_windows_${Arch}.zip"
$Url     = "https://github.com/$Repo/releases/download/$Version/$Archive"

$InstallDir = "$env:LOCALAPPDATA\Programs\dwm"
$TmpDir     = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()

Write-Host "Installing dwm $Version..." -ForegroundColor Cyan

# Download
New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null
$ZipPath = Join-Path $TmpDir $Archive
Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing

# Extract
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

# Install
New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
Copy-Item -Path (Join-Path $TmpDir "dwm.exe") -Destination $InstallDir -Force

# Add to PATH (user scope)
$UserPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [System.Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to PATH" -ForegroundColor Green
}

# Cleanup
Remove-Item -Recurse -Force $TmpDir

Write-Host ""
Write-Host "  dwm $Version installed to $InstallDir\dwm.exe" -ForegroundColor Green
Write-Host "  Restart your terminal, then run: dwm --help" -ForegroundColor Yellow
Write-Host ""
