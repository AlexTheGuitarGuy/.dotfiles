alias gst='git status'
alias gsh='git stash'
alias gsha='git stash apply'
alias gl='git log | less'
alias ga='git add .'
alias gac='git add . && git commit -m $1'
alias gp='git push'
alias gpl='git pull'
alias gch='git checkout'
alias gcd='git checkout dev'
alias gcb='git checkout -b'
alias grv='git checkout .'
alias gfc='git fetch'
alias gfd='git fetch origin dev:dev'
alias gmd='git merge origin/dev'
alias grb='git rebase origin/dev'
alias ghd='git update-index --skip-worktree $1'
alias guhd='git update-index --no-skip-worktree $1'
alias sz='source ~/.zshrc'
alias oz='nvim ~/.zshrc'
alias ob='nvim ~/.bashrc'
alias onc='cd ~/.config/nvim/lua/user/ && nvim .'
alias ffx='firefox-developer-edition'
alias tgs='~/.dotfiles/scripts/tmux-session-generator.sh'
alias d='doppler run --' 
alias srv='echo ${NODE_AUTH_TOKEN} | sudo make serve NODE_AUTH_TOKEN=xargs'
alias rbt='sudo modprobe -r btusb && sudo modprobe btusb'
