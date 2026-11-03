$src = Join-Path $PSScriptRoot "config"
$lw = Join-Path $env:USERPROFILE ".librewolf"
$tri = Join-Path $env:USERPROFILE ".config\tridactyl"
$triThemes = Join-Path $tri "themes"

New-Item -ItemType Directory -Force -Path $lw, $tri, $triThemes | Out-Null
Copy-Item (Join-Path $src "librewolf.overrides.cfg") $lw -Force
Copy-Item (Join-Path $src "tridactylrc") $tri -Force
Copy-Item (Join-Path $src "tridactyl.css") (Join-Path $triThemes "catppuccin.css") -Force

Get-ChildItem (Join-Path $env:APPDATA "librewolf\Profiles") -Directory | ForEach-Object {
    $chrome = Join-Path $_.FullName "chrome"
    New-Item -ItemType Directory -Force -Path $chrome | Out-Null
    Copy-Item (Join-Path $src "userChrome.css"), (Join-Path $src "userContent.css") $chrome -Force
    Write-Host "copied chrome to $chrome"
}
