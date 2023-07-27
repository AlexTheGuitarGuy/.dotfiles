#!/bin/sh

[ -f "$HOME/.local/share/zap/zap.zsh" ] && source "$HOME/.local/share/zap/zap.zsh"
# zmodload zsh/zprof
eval "$(fnm env --use-on-cd)"

# imports
plug "$HOME/.config/zsh/aliases.zsh"
plug "$HOME/.config/zsh/exports.zsh"
plug "$HOME/.config/zsh/plugins.zsh"

bindkey '^ ' autosuggest-accept


# Load Angular CLI autocompletion.
# source <(ng completion script)


# Caching zsh compinit for better performance
autoload -Uz compinit
for dump in ~/.zcompdump(N.mh+24); do
  compinit
done
compinit -C
# zprof
