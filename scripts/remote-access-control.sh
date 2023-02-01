#!/bin/bash
# 
# Usage:
# 
#   ./scripts/remote-access-control.sh ./bin/uhppote-cli 405419896
# 
# Ref.:
#   1. https://stackoverflow.com/questions/3004811/how-do-you-run-multiple-programs-in-parallel-from-a-bash-script
#   2. https://unix.stackexchange.com/questions/7558/execute-a-command-once-per-line-of-piped-input

CLI=$1
CONTROLLER=$2
CARDS=(
   8165538
   10058399
   10058400
)

function set_pc_control {
    $CLI set-pc-control $CONTROLLER > /dev/null

    while :; do
        sleep 5
        $CLI get-time $CONTROLLER > /dev/null
    done
}

function listen {
    $CLI listen | while read -r EVENT; do on_event "$EVENT"; done
}

function on_event {
    cid=$(cut -f1 -d ' ' <<< "$1")
    evt=$(cut -f2 -d '|' <<< "$1")

    if [[ "${cid}" =~ ^[0-9]+$ ]]; then
        IFS=' ' read -ra array <<< "$evt"
        card="${array[5]}"
        door="${array[3]}"
        reason="${array[8]}"
        granted="${array[2]}"

        if [[ "${granted}" == "false" && "${reason}" == "5" && "${card}" =~ ^[0-9]+$ && "${door}" =~ ^[1234]$ ]]; then

            if lookup "${card}"; then 
                echo "SWIPE  $cid  $card  $door    ACCESS GRANTED"
                $CLI open "${cid}" "${door}" > /dev/null
            else
                echo "SWIPE  $cid  $card  $door    ACCESS DENIED"
            fi
        fi
    fi
}

function lookup {
    card="$1"

    for c in "${CARDS[@]}"; do
        if [[ "${card}" == "${c}" ]]; then
            return 0
        fi
    done

    return 1
}

(trap 'kill 0' SIGINT; set_pc_control & listen)