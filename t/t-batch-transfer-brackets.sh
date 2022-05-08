#!/usr/bin/env bash

. "$(dirname "$0")/testlib.sh"

# These tests rely on behavior found in 2.3.4 to perform themselves,
# specifically:
#   - Parsing username followed by literal IPv6 address in SSH transport URLs.

ensure_git_version_isnt $VERSION_LOWER "2.3.4"

begin_test "batch transfers with ssh custom port bare endpoint (git-lfs-transfer)"
(
  set -e

  setup_pure_ssh

  reponame="batch-ssh-transfer-bare-port"
  setup_remote_repo "$reponame"
  clone_repo "$reponame" "$reponame"

  sshurl="[git@127.0.0.1:22]:$repodir"
  git config lfs.url "$sshurl"

  contents="test"
  git lfs track "*.dat"
  printf "%s" "$contents" > test.dat
  git add .gitattributes test.dat
  git commit -m "initial commit"

  git push origin main 2>&1
  cd ..
  GIT_TRACE=1 GIT_SSH_VARIANT=ssh \
    git clone "$sshurl" "$reponame-2" 2>&1 | tee trace.log
  grep "lfs-ssh-echo.*git-lfs-transfer .*$reponame.git download" trace.log
  cd "$reponame-2"
  git lfs fsck
)
end_test

begin_test "batch transfers with ssh bracketed host bare endpoint (git-lfs-transfer)"
(
  set -e

  setup_pure_ssh

  reponame="batch-ssh-transfer-bare-brackets"
  setup_remote_repo "$reponame"
  clone_repo "$reponame" "$reponame"

  sshurl="git@[127.0.0.1]:$repodir"
  git config lfs.url "$sshurl"

  contents="test"
  git lfs track "*.dat"
  printf "%s" "$contents" > test.dat
  git add .gitattributes test.dat
  git commit -m "initial commit"

  git push origin main 2>&1
  cd ..
  GIT_TRACE=1 GIT_SSH_VARIANT=ssh \
    git clone "$sshurl" "$reponame-2" 2>&1 | tee trace.log
  grep "lfs-ssh-echo.*git-lfs-transfer .*$reponame.git download" trace.log
  cd "$reponame-2"
  git lfs fsck
)
end_test
