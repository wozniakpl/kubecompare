_kubecompare() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    if [[ $cur == -* ]]; then
        COMPREPLY=( $(compgen -W "-n --namespace" -- $cur) )
        return 0
    fi

    if [[ $prev == "kubecompare" || $prev == "-n" || $prev == "--namespace" ]]; then
        COMPREPLY=( $(kubectl get deployment,statefulset,daemonset -o custom-columns=NAME:.metadata.name --no-headers 2>/dev/null) )
        return 0
    fi
}

complete -F _kubecompare kubecompare
