# Links Claude Code's user rules and skills on Windows to this clone, so `git pull` updates them.
# Run once from PowerShell (no admin needed): powershell -ExecutionPolicy Bypass -File windows\link-claude.ps1
# Rerun after moving the clone. An existing real folder is renamed to <name>.bak-<date>, never deleted.

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
