# Server Discovery

One of the design goals for Git LFS is that it should work using as few
required configuration properties as possible.  In particular, Git LFS will
attempt to determine the Git LFS server to use based on your Git remote
settings, but you can also configure a custom Git LFS server if your Git
remote doesn't support one, or you just want to use a different one.

Look for the `Endpoint` properties in `git lfs env` to see your current LFS
servers.

## Guessing the Remote

By default, to find the name of the Git remote from which it will derive
a Git LFS server URL, Git LFS tries each of the following steps, stopping
when a remote name is selected.

1. If the current branch tracks a remote branch, the name of that remote
   as set by the `branch.<branch-name>.remote` Git configuration option
   will be used.
2. If the `remote.lfsDefault` Git configuration option is defined, its
   value will be used.
3. If exactly one Git remote is defined, its name will be used.
4. If none of the above conditions apply, the name `origin` will be used.

For Git LFS upload operations, some additional steps to select the name of
the remote will be tried first, before the steps listed above.

1. If the current branch tracks a remote branch for push operations, the
   name of that remote as set by the `branch.<branch-name>.pushRemote`
   Git configuration option will be used.
2. If the `remote.lfsPushDefault` Git configuration option is defined, its
   value will be used.
3. If the `remote.pushDefault` Git configuration option is defined, its
   value will be used.
4. If none of the above conditions apply, the other steps are followed,
   starting with checking for a defined `branch.<branch-name>.remote` Git
   configuration option.



* `lfs.remote.autodetect` - default false

This boolean option enables the remote autodetect feaure within Git LFS.
LFS tries to derive the corresponding remote from the commit information
and, in case of success, ignores the settings defined by
`remote.lfsdefault` and `remote.<remote>.lfsurl`.

* `lfs.remote.searchall` - default false

This boolean option enables Git LFS to search all registered remotes to
find LFS data. This is a fallback mechanism executed only if the LFS
data could not be found via the ordinary heuristics as described in
`remote.lfsdefault`, `remote.<remote>.lfsurl` and, if enabled,
`lfs.remote.autodetect`.


## Guessing the Server

In general, to construct the server URL it will use, Git LFS appends
`.git/info/lfs` to the end of the URL returned by Git for the remote
chosen by Git LFS, as shown in the examples below.

Git Remote: `https://git-server.com/foo/bar`<br>
LFS Server: `https://git-server.com/foo/bar.git/info/lfs`

Git Remote: `https://git-server.com/foo/bar.git`<br>
LFS Server: `https://git-server.com/foo/bar.git/info/lfs`

Git Remote: `git@git-server.com:foo/bar.git`<br>
LFS Server: `https://git-server.com/foo/bar.git/info/lfs`

Git Remote: `ssh://git-server.com/foo/bar.git`<br>
LFS Server: `https://git-server.com/foo/bar.git/info/lfs`

Git Remote: `file://foo/bar`<br>
LFS Server: `file://foo/bar`

The Git remote URL will first be rewritten according to the longest-matching
`url.<base>.insteadOf` Git configuration, if any.  (For Git LFS upload
operations, the longest-matching `url.<base>.pushInsteadOf` Git configuration,
if any, will take precedence over any `url.<base>.insteadOf` configurations.)

If the result is a URL with a `file://` scheme, that will be used as the
Git LFS endpoint without further modification.  Otherwise, `/info/lfs` will
be appended to the URL, preceded by `.git` unless the URL already ends with
`.git`, and if the scheme is anything other than `file://` or `https://`
the URL will be converted into an `https://` URL.

Git LFS also supports various non-URL Git remotes.  Simple file paths
are converted to `file://` URLs, if the referenced location exists on
the local system.  Remotes specified using the SSH syntax
`[user@]host:/path/to/repo.git` are converted to `ssh://` URLs.

## SSH

If Git LFS detects an SSH remote, it will run the `git-lfs-authenticate`
command. This allows supporting Git servers to give the Git LFS client
alternative authentication so the user does not have to setup a git credential
helper.

Git LFS runs the following command:

    $ ssh [{user}@]{server} git-lfs-authenticate {path} {operation}

The `user`, `server`, and `path` properties are taken from the SSH remote. The
`operation` can either be "download" or "upload". The SSH command can be
tweaked with the `GIT_SSH` or `GIT_SSH_COMMAND` environment variables. The
output for successful commands is JSON, and matches the schema as an `action`
in a Batch API response. Git LFS will dump the STDERR from the `ssh` command if
it returns a non-zero exit code.

Examples:

The `git-lfs-authenticate` command can even suggest an LFS endpoint that does
not match the Git remote by specifying an `href` property.

```bash
# Called for remotes like:
#   * git@git-server.com:foo/bar.git
#   * ssh://git@git-server.com/foo/bar.git
$ ssh git@git-server.com git-lfs-authenticate foo/bar.git download
{
  "href": "https://lfs-server.com/foo/bar",
  "header": {
    "Authorization": "RemoteAuth some-token"
  },
  "expires_in": 86400
}
```

Git LFS will output the STDERR if `git-lfs-authenticate` returns a non-zero
exit code:

```bash
$ ssh git@git-server.com git-lfs-authenticate foo/bar.git wat
Invalid LFS operation: "wat"
```

## Custom Configuration

If Git LFS can't guess your LFS server, or you aren't using the
`git-lfs-authenticate` command, you can specify the LFS server using Git config.

Set `lfs.url` to set the LFS server, regardless of Git remote.

```bash
$ git config lfs.url https://lfs-server.com/foo/bar
```

You can set `remote.<name>.lfsurl` to set the LFS server for that specific
remote only:

```bash
$ git config remote.dev.lfsurl http://lfs-server.dev/foo/bar
$ git lfs env
...

Endpoint=https://git-server.com/foo/bar.git/info/lfs (auth=none)
Endpoint (dev)=http://lfs-server.dev/foo/bar (auth=none)
```

To set a distinct server URL for Git LFS upload operations only, use
`lfs.pushurl` or `remote.<name>.lfspushurl`.

Any of these four configuration options (`lfs.url`, `lfs.pushurl`,
`remote.<name>.lfsurl`, and `remote.<name>.lfspushurl`), as well as
a restricted set of other options, will be read by Git LFS if they
are defined in an `.lfsconfig` file in the root of your repository.
Unlike Git's own configuration files, the `.lfsconfig` file may be
committed to the repository, which allows it to be shared with all
users of the repository.  Note that if these options are also defined
in a Git configuration file, those settings will take precedence over
the ones in the `.lfsconfig` file.

```bash
$ git config --file=.lfsconfig lfs.url https://lfs-server.com/foo/bar
```
