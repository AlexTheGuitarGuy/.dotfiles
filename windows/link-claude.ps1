# Links Claude Code's user rules and skills on Windows to this clone, so `git pull` updates them, and merges
# the shared settings (plugins, model) into settings.json.
# Run from PowerShell (no admin needed): powershell -ExecutionPolicy Bypass -File windows\link-claude.ps1
# Rerun after a pull that changed settings.portable.json, or after moving the clone.
# An existing real folder is renamed to <name>.bak-<date>, never deleted; settings set only on Windows are kept.

$ErrorActionPreference = 'Stop'
$repo = Split-Path -Parent $PSScriptRoot
$claude = Join-Path $env:USERPROFILE '.claude'
New-Item -ItemType Directory -Force -Path $claude | Out-Null

foreach ($name in 'rules', 'skills') {
    $target = Join-Path $repo "modules\home-manager\claude\$name"
    $link = Join-Path $claude $name
    $item = Get-Item -LiteralPath $link -Force -ErrorAction SilentlyContinue
    if ($item -and $item.LinkType -eq 'Junction' -and $item.Target -eq $target) {
        Write-Host "ok      $link -> $target"
        continue
    }
    if ($item -and $item.LinkType) {
        $item.Delete()  # removes only the old link, not what it points to
    } elseif ($item) {
        $backup = "$link.bak-$(Get-Date -Format yyyyMMdd-HHmmss)"
        Rename-Item -LiteralPath $link -NewName (Split-Path -Leaf $backup)
        Write-Host "moved   $link -> $backup"
    }
    New-Item -ItemType Junction -Path $link -Target $target | Out-Null
    Write-Host "linked  $link -> $target"
}

function Merge-Json($base, $over) {
    foreach ($p in $over.PSObject.Properties) {
        $cur = $base.PSObject.Properties[$p.Name]
        if ($cur -and $cur.Value -is [pscustomobject] -and $p.Value -is [pscustomobject]) {
            Merge-Json $cur.Value $p.Value
        } elseif ($cur) {
            $cur.Value = $p.Value
        } else {
            $base | Add-Member -NotePropertyName $p.Name -NotePropertyValue $p.Value
        }
    }
}

$portable = Join-Path $repo 'modules\home-manager\claude\settings.portable.json'
$settingsFile = Join-Path $claude 'settings.json'
$settings = if (Test-Path -LiteralPath $settingsFile) {
    Copy-Item -LiteralPath $settingsFile "$settingsFile.bak" -Force
    Get-Content -Raw -LiteralPath $settingsFile | ConvertFrom-Json
} else { $null }
if (-not $settings) { $settings = [pscustomobject]@{} }
Merge-Json $settings (Get-Content -Raw -LiteralPath $portable | ConvertFrom-Json)
# UTF-8 without BOM: Windows PowerShell 5 would otherwise write one
[IO.File]::WriteAllText($settingsFile, ($settings | ConvertTo-Json -Depth 20), (New-Object Text.UTF8Encoding $false))
Write-Host "merged  $portable -> $settingsFile"
