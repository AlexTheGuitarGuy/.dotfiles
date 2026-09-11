#!/bin/sh

[ -f "$HOME/.local/share/zap/zap.zsh" ] && source "$HOME/.local/share/zap/zap.zsh"
# zmodload zsh/zprof

# imports
plug "$HOME/.config/zsh/aliases.zsh"
plug "$HOME/.config/zsh/exports.zsh"
plug "$HOME/.config/zsh/plugins.zsh"

# Disable XON/XOFF flow control so ctrl+s / ctrl+q are usable keys instead of freezing output
[ -t 0 ] && stty -ixon

# Secrets (sops-encrypted, decrypted into env at shell start)
if [ -f "$HOME/.config/sops/age/keys.txt" ]; then
  while IFS= read -r line; do
    export "${line%%: *}=${line#*: }"
  done < <(sops -d "$HOME/.dotfiles/secrets/secrets.yaml")
fi

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
