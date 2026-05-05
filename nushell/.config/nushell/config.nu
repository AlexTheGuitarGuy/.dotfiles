# config.nu
#
# Installed by:
# version = "0.103.0"
#
# This file is used to override default Nushell settings, define
# (or import) custom commands, or run any other startup tasks.
# See https://www.nushell.sh/book/configuration.html
#
# This file is loaded after env.nu and before login.nu
#
# You can open this file in your default editor using:
# config nu
#
# See `help config nu` for more options
#
# You can remove these comments if you want or leave
# them for future reference.

# Configure env
$env.config = {
  edit_mode: vi
  show_banner: false
  highlight_resolved_externals: true
  color_config: {
    shape_external_resolved: "green"
    shape_external: "red"
  }
}
$env.EDITOR = "nvim"
$env.ENV_CONVERSIONS = {
    PATH: {
        from_string: { |s| $s | split row ':' },
        to_string: { |v| $v | str join ':' }
    }
}

# Carapace
source ~/.cache/carapace/init.nu

# Zoxide
source ~/.zoxide.nu

# Secrets
source ~/.config/nushell/secrets.nu

# Aliases
alias gst = git status
alias gsh = git stash
alias gsha = git stash apply
alias gl = git log | less
alias ga = git add .
def gac [message] {
    git add .
    git commit -m $message
}
alias gp = git push
alias gpl = git pull
alias gch = git checkout
alias gcd = git checkout dev
alias gcb = git checkout -b
alias grv = git checkout .
alias gfc = git fetch
alias gfd = git fetch origin dev:dev
alias gmd = git merge origin/dev
alias grb = git rebase origin/dev
def ghd [file] {
    git update-index --skip-worktree $file
}
def guhd [file] {
    git update-index --no-skip-worktree $file
}
alias tgs = ~/.dotfiles/scripts/tmux-session-generator.sh
def rbt [] {
    sudo modprobe -r btusb
    sudo modprobe btusb
}
def zvim [pattern?: string] {
  let dir = if ($pattern | is-empty) {
      zoxide query -i | str trim
  } else {
      zoxide query $pattern | str trim
  }
  
  if ($dir | is-empty) {
      error make {msg: "No matching directory found"}
      return
  }
  
  nvim $dir
}

# Plugins
# source ~/.local/share/atuin/init.nu

# FNM
# Only run if fnm is installed
if not (which fnm | is-empty) {
    # Load environment variables from fnm in JSON format
    ^fnm env --json | from json | load-env

    # Add Node's bin directory to PATH
    let node_path = if $nu.os-info.name == "windows" {
        $env.FNM_MULTISHELL_PATH
    } else {
        $"($env.FNM_MULTISHELL_PATH)/bin"
    }
    $env.PATH = ($env.PATH | prepend [ $node_path ])
}

def connect-to-dvag-vpn [] {
  let $orig_dir = $env.PWD

  cd /home/alex/Apps/openfortivpn-webview/openfortivpn-webview-electron

  $env.DVAG_VPN_COOKIE = (npx electron index.js fg.zentrale.dvag:443)

  cd $orig_dir

  sudo openconnect --protocol=fortinet fg.zentrale.dvag:443 --cookie $env.DVAG_VPN_COOKIE
}
source $"($nu.home-dir)/.cargo/env.nu"
