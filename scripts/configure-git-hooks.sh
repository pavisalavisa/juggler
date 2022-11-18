#!/usr/bin/env bash

PROJECT_ROOT_DIR=$(pwd)
GIT_HOOKS_PATH="${PROJECT_ROOT_DIR}/.git/hooks"
CUSTOM_GIT_HOOKS_PATH="${PROJECT_ROOT_DIR}/scripts/git-hooks"

printf "Configuring pre-commit hook found\n"
if [[ -e "${GIT_HOOKS_PATH}/pre-commit" ]]; then
    printf "The pre-commit hook is already configured.\n"
    exit 0
fi

ln -s "$CUSTOM_GIT_HOOKS_PATH/pre-commit" "$GIT_HOOKS_PATH/pre-commit" 
if [ $? -ne 0 ]; then
    printf 1>&2 " Error on copying the pre-commit script.\n"
    exit 1
fi

printf "pre-commit hook configured.\n"
exit 0
