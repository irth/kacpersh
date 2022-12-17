_kacpersh_start() {
    curl -s --unix-socket "$KACPERSH_SOCK" http://kacpersh/start
}
_kacpersh_stop() {
    curl -s --unix-socket "$KACPERSH_SOCK" http://kacpersh/stop
}

if [[ "$KACPERSH_CUSTOM_HOOKS" != "1" ]] && [[ -n "$KACPERSH_SOCK" ]]; then
    autoload add-zsh-hook
    add-zsh-hook precmd _kacpersh_stop
    add-zsh-hook preexec _kacpersh_start

    kacpersh_last() {
        if [[ -n "$KACPERSH_SOCK" ]]; then 
            curl -s --output - --unix-socket "$KACPERSH_SOCK" http://kacpersh/last 
        fi
    }

    alias _=kacpersh_last
fi
