<#
Creates a .zip archive for every sub-folder found inside the sibling Release directory.
Each archive is placed into a ZIPS directory next to Release.

Example result structure:
  Release/
    windows_amd64/  (binary files…)
    linux_amd64/    (binary files…)
  ZIPS/
    windows_amd64.zip
    linux_amd64.zip

Run from any location:
  pwsh path/to/zip_release.ps1
#>

$ErrorActionPreference = "Stop"

#-------------  LOCATIONS -----------------------------------------------
# Script directory (repository root when script lives at repo root)
$RepoRoot  = $PSScriptRoot

# Folder that holds platform-specific build outputs
$BinariesDir = Join-Path -Path $RepoRoot -ChildPath "Release\Binaries"

# Where the generated archives will be stored
$ZipsDir     = Join-Path -Path $RepoRoot -ChildPath "Release\ZIPS"
#-----------------------------------------------------------------------

# Ensure the binaries folder exists
if (-not (Test-Path $BinariesDir)) {
    throw "Binaries directory not found: $BinariesDir"
}

# Create the ZIPS folder if missing
if (-not (Test-Path $ZipsDir)) {
    New-Item -ItemType Directory -Path $ZipsDir | Out-Null
    Write-Host "Created ZIPS directory at $ZipsDir"
}

# Enumerate each immediate sub-folder in Release/Binaries
Get-ChildItem -Path $BinariesDir -Directory | ForEach-Object {
    $folder       = $_
    $destination  = Join-Path -Path $ZipsDir -ChildPath ($folder.Name + '.zip')

    # If a zip with this name already exists, overwrite it
    if (Test-Path $destination) {
        Remove-Item $destination -Force
    }

    # Each folder should contain exactly one binary (with or without .exe)
    $binary = Get-ChildItem -Path $folder.FullName -File | Select-Object -First 1
    if (-not $binary) {
        Write-Warning "No binary found in $($folder.FullName). Skipping."
        return
    }

    Write-Host "Zipping $($binary.Name) from $($folder.Name) -> $destination"
    Compress-Archive -Path $binary.FullName -DestinationPath $destination -Force
}

Write-Host "`nDone. Archives created in $ZipsDir" 