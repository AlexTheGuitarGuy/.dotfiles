#!/bin/sh

[ -f "$HOME/.local/share/zap/zap.zsh" ] && source "$HOME/.local/share/zap/zap.zsh"
# zmodload zsh/zprof

# imports
plug "$HOME/.config/zsh/aliases.zsh"
plug "$HOME/.config/zsh/exports.zsh"
plug "$HOME/.config/zsh/plugins.zsh"

# Disable XON/XOFF flow control so ctrl+s / ctrl+q are usable keys instead of freezing output
[ -t 0 ] && stty -ixon

# Secrets (sops-encrypted, decrypted per command, never exported to the shell)
_sec() {
  local f
  for f in "$HOME/.dotfiles/secrets/secrets.yaml" "$HOME"/.config/secrets.d/*.yaml; do
    [ -e "$f" ] && sops -d --extract "[\"$1\"]" "$f" 2>/dev/null && return
  done
  return 1
}
[ -f "$HOME/.zshrc.private" ] && source "$HOME/.zshrc.private"

bindkey '^ ' autosuggest-accept


# Load Angular CLI autocompletion.
# source <(ng completion script)

# . ~/.nix-profile/etc/profile.d/nix.sh

# Caching zsh compinit for better performance
autoload -Uz compinit
for dump in ~/.zcompdump(N.mh+24); do
  compinit
done
compinit -C

export PATH=$PATH:/home/alex/.spicetify

# starship
eval "$(starship init zsh)"

. "$HOME/.atuin/bin/env"

# eval "$(atuin init zsh)"
