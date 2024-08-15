#!/bin/bash

golangci_lint_download_url=github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Function to check and install a command if missing
install_if_missing() {
  if ! command -v "$1" &> /dev/null; then
    printf "\n%s missing. installing...\n" "$1"
    go install $2
  fi
}

install_if_missing golangci-lint $golangci_lint_download_url

nodejs_version=v20.12.0
nvm_installation_script_url=https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh

if ! command -v npm &> /dev/null; then

    # shellcheck disable=SC2059
    printf "\nnpm not available. installing nodejs$nodejs_version using nvm...\n"
    nodejs_installation_dir=$installation_dir/nodejs

    if ! command -v node &> /dev/null; then

#       TODO: Fork this script, save nvm to /tmp/bin and all other dependencies to this directory, then install node and npm to this directory as well
        curl -o- "$nvm_installation_script_url" | bash

        # shellcheck disable=SC1090
        source ~/.bashrc

        printf "\ninstalling nvm v0.39.3\n"
        nvm install "$nodejs_version"

    fi

    NODE_HOME=$nodejs_installation_dir
    PATH="$NODE_HOME/bin:$PATH"

fi

if ! command -v commitlint &> /dev/null; then
  printf "\ncommitlint not available. installing latest commitlint globally...\n"
  npm install --global @commitlint/cli
fi

printf "\nsetting up npm repository for commitlint\n"
npm init -y &> /dev/null
printf "\ninstalling commit lint as dev dependency\n"
npm install --save-dev @commitlint/config-conventional &> /dev/null

printf "\nall git hook dependencies installed/available. proceeding with setup...\n"

cp "$parent_dir/.githooks/commit-msg" "$parent_dir/.git/hooks/"
chmod +x "$parent_dir/.git/hooks/commit-msg"

cp "$parent_dir/.githooks/pre-commit" "$parent_dir/.git/hooks/pre-commit"
chmod +x "$parent_dir/.git/hooks/pre-commit"

printf "\nsetup hook dependencies and local hooks\n"
