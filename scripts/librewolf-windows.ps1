$repo = Split-Path -Parent $PSScriptRoot
$src = Join-Path $repo "librewolf\.config\librewolf\librewolf"
$lw = Join-Path $env:USERPROFILE ".librewolf"
$tri = Join-Path $env:USERPROFILE ".config\tridactyl"

New-Item -ItemType Directory -Force -Path $lw, $tri | Out-Null
Copy-Item (Join-Path $src "librewolf.overrides.cfg") $lw -Force
Copy-Item (Join-Path $repo "librewolf\.config\tridactyl\tridactylrc") $tri -Force

Get-ChildItem (Join-Path $env:APPDATA "librewolf\Profiles") -Directory | ForEach-Object {
    $chrome = Join-Path $_.FullName "chrome"
    New-Item -ItemType Directory -Force -Path $chrome | Out-Null
    Copy-Item (Join-Path $src "chrome\*") $chrome -Force
    Write-Host "copied chrome to $chrome"
}
