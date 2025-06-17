<#
Builds GoMergeLog4Net for Windows, Linux-amd64, macOS-amd64, macOS-arm64.
Each artefact is placed in   Release\<goos>_<goarch>\GoMergeLog4Net(.exe)
#>

$ErrorActionPreference = "Stop"

#-------------  CONFIG -------------------------------------------------
$SrcDir  = $PSScriptRoot          # ← main.go is right here
$OutRoot = Join-Path $SrcDir "Release\Binaries"
$AppName = "GoMergeLog4Net"
$Version = git -C $SrcDir describe --tags --always 2>$null
if (-not $Version) { $Version = "dev" }
#-----------------------------------------------------------------------

$targets = @(
  @{ GOOS="windows"; GOARCH="amd64"; EXT=".exe" },
  @{ GOOS="linux";   GOARCH="amd64"; EXT=""    },
  @{ GOOS="darwin";  GOARCH="amd64"; EXT=""    },
  @{ GOOS="darwin";  GOARCH="arm64"; EXT=""    }
)
Write-Host "STARTING" 
mkdir $OutRoot -ea 0 | Out-Null
Write-Host "$OutRoot" 
Write-Host "Building version: $Version" 
foreach ($t in $targets) {
  $outDir = "$OutRoot/$($t.GOOS)_$($t.GOARCH)"
  mkdir $outDir -ea 0 | Out-Null

  Write-Host "=> building $($t.GOOS)/$($t.GOARCH)"
  $env:GOOS        = $t.GOOS
  $env:GOARCH      = $t.GOARCH
  $env:CGO_ENABLED = "0"

  go build -ldflags "-s -w -X main.version=$Version" `
           -o "$outDir/$AppName$($t.EXT)" `
           "./"
}

Write-Host "`nBinaries are under $OutRoot"