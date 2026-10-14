# Tải binary mới nhất từ GitHub Releases và chạy.
# Dùng: irm https://raw.githubusercontent.com/dominhhieu1405/Simple-Auction-system/main/scripts/run.ps1 | iex
$ErrorActionPreference = 'Stop'

$Repo    = 'dominhhieu1405/Simple-Auction-system'
$Version = if ($env:VERSION) { $env:VERSION } else { 'latest' }
$Dir     = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { 'simple-auction' }

$arch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }
$name = "auction-server-windows-$arch.exe"

$url = if ($Version -eq 'latest') {
    "https://github.com/$Repo/releases/latest/download/$name"
} else {
    "https://github.com/$Repo/releases/download/$Version/$name"
}

New-Item -ItemType Directory -Force -Path $Dir | Out-Null
Set-Location $Dir

Write-Host "Downloading $url"
Invoke-WebRequest -Uri $url -OutFile 'auction-server.exe' -UseBasicParsing

if (-not (Test-Path '.env')) {
    $bytes = New-Object byte[] 32
    [System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
    $secret = -join ($bytes | ForEach-Object { $_.ToString('x2') })
    "JWT_SECRET=$secret`nDB_PATH=data/auction.db`nSERVER_PORT=8080" |
        Set-Content -Encoding ascii '.env'
    Write-Host 'Created .env with random JWT_SECRET'
}

Write-Host "Starting server in $(Get-Location)"
.\auction-server.exe