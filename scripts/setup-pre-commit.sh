#!/bin/bash
ROOT_DIR=$(git rev-parse --show-toplevel)
cd $ROOT_DIR

if ! command -v pre-commit &> /dev/null
then
		brew install pre-commit
		printf '%s\n' "Installed pre-commit"
else
		printf '%s\n' "pre-commit is already installed"
fi

if pre-commit install --hook-type pre-commit --config .pre-commit-config.yml; then
		printf '%s\n' "Configured pre-commit"
else
		printf '\t%s\n' "FAIL"
fi

if pre-commit install --hook-type pre-push --config .pre-push-config.yml;  then
		printf '%s\n' "Configured pre-push"
else
		printf '\t%s\n' "FAIL"
fi
