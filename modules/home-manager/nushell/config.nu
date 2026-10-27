def gac [message] {
    git add .
    git commit -m $message
}

def ghd [file] {
    git update-index --skip-worktree $file
}

def guhd [file] {
    git update-index --no-skip-worktree $file
}

def --wrapped lazygit [...rest] {
    let main = ($env.PWD | path join main)
    if ($env.PWD | path join .bare | path exists) and ($main | path exists) { cd $main }
    ^lazygit ...$rest
}

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

if (^tty | complete | get exit_code) == 0 {
    ^stty -ixon
}
