#!/usr/bin/env bash
set -eu

prefix="/usr/local"

if [ "${PREFIX:-}" != "" ] ; then
  prefix=${PREFIX:-}
elif [ "${BOXEN_HOME:-}" != "" ] ; then
  prefix=${BOXEN_HOME:-}
fi

while [[ $# -gt 0 ]]; do
  case "$1" in
    --local)
      prefix="$HOME/.local"
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Check if the user has permission to install in the specified prefix
if [ ! -w "$prefix" ]; then
  echo "Error: Insufficient permissions to install in $prefix. Try running with sudo or choose a different prefix.">&2
  exit 1
fi

mkdir -p "$prefix/bin"
rm -rf "$prefix/bin/git-lfs*"

pushd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null
  for g in git*; do
    install "$g" "$prefix/bin/$g"
  done
popd > /dev/null

if [ "$(id -u)" -eq 0 ]; then
  cat <<-EOF
	The Git LFS binary is now installed in '$prefix/bin'.

	To configure Git LFS as a regular user, run 'git lfs install'.

	To configure Git LFS for all users of the system, run:

	  $ git config set --system core.hooksPath <system-hooks-path>
	  $ git lfs install --system

	Note that you may need to run these commands with 'sudo'.

	If a system-wide Git hooks location is already configured, you can skip
	the first command.  Otherwise, <system-hooks-path> should specify an
	absolute path to your preferred location for system-wide Git hooks.

EOF
else
  PATH+=:"$prefix/bin"
  git lfs install
fi
