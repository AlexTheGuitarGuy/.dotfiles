# Setting for the new UTF-8 terminal support in Lion
export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8

export GEM_HOME="$(ruby -e 'puts Gem.user_dir')"
export PATH="$PATH:$GEM_HOME/bin"
export PATH="$PATH:/home/alex/.cargo/bin"

PATH="/home/alex/perl5/bin${PATH:+:${PATH}}"
export PATH
PERL5LIB="/home/alex/perl5/lib/perl5${PERL5LIB:+:${PERL5LIB}}"
export PERL5LIB
PERL_LOCAL_LIB_ROOT="/home/alex/perl5${PERL_LOCAL_LIB_ROOT:+:${PERL_LOCAL_LIB_ROOT}}"
export PERL_LOCAL_LIB_ROOT
PERL_MB_OPT="--install_base \"/home/alex/perl5\""
export PERL_MB_OPT
PERL_MM_OPT="INSTALL_BASE=/home/alex/perl5"
export PERL_MM_OPT
export PATH="/usr/pgsql-14/bin:$PATH"

# Load nvm
export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
# This loads nvm bash_completion
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"

# go
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

export DISPLAY=$(grep -m 1 nameserver /etc/resolv.conf | awk '{print $2}'):0
export LIBGL_ALWAYS_INDIRECT=1

export NODE_OPTIONS="--max-old-space-size=5120" # Increase to 5 GB
export BROWSER=wslview
export PATH="$PATH:/opt/nvim-linux64/bin"


# ===== i3 stuff =====
export DISPLAY=:0
export LIBGL_ALWAYS_INDIRECT=1
export XDG_RUNTIME_DIR=/tmp/xdg
export RUNLEVEL=3
# ===== fix mkdir wrong perms =====
if grep -q Microsoft /proc/version; then
  if [ "0022" == '0000' ]; then
    umask 0022
  fi
fi

autoload -U +X bashcompinit && bashcompinit
complete -o nospace -C /usr/bin/mcli mcli

# pnpm
export PNPM_HOME="/home/alex/.local/share/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac
# pnpm end

# bun completions
[ -s "/home/alex/.bun/_bun" ] && source "/home/alex/.bun/_bun"

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"
eval "$(starship init zsh)"

# zprof

# fnm
# FNM_PATH="/home/alex/.local/share/fnm"
# if [ -d "$FNM_PATH" ]; then
#   export PATH="/home/alex/.local/share/fnm:$PATH"
#   eval "`fnm env`"
# fi

export NVM_DIR="$HOME/.nvm"
 -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion

export PATH=$PATH:~/.npm-global/bin

export MANPAGER="nvim +Man!"
