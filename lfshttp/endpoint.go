package lfshttp

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/git-lfs/git-lfs/v3/git"
	"github.com/git-lfs/git-lfs/v3/ssh"
)

const UrlUnknown = "<unknown>"

// An Endpoint describes how to access a Git LFS server.
type Endpoint struct {
	Url         string
	SSHMetadata ssh.SSHMetadata
	Operation   string
}

func endpointOperation(e Endpoint, method string) string {
	if len(e.Operation) > 0 {
		return e.Operation
	}

	switch method {
	case "GET", "HEAD":
		return "download"
	default:
		return "upload"
	}
}

var sshHostPortRE = regexp.MustCompile(`^([^\:]+)(?:\:(\d+))?$`)

// EndpointFromSshUrl constructs a new endpoint from an ssh:// URL
func EndpointFromSshUrl(u *url.URL) Endpoint {
	var endpoint Endpoint
	// Pull out port now, we need it separately for SSH
	match := sshHostPortRE.FindStringSubmatch(u.Host)
	if match == nil || len(match) < 2 {
		endpoint.Url = UrlUnknown
		return endpoint
	}

	host := match[1]
	if u.User != nil && u.User.Username() != "" {
		endpoint.SSHMetadata.UserAndHost = fmt.Sprintf("%s@%s", u.User.Username(), host)
	} else {
		endpoint.SSHMetadata.UserAndHost = host
	}

	if len(match) > 2 {
		endpoint.SSHMetadata.Port = match[2]
	}

	endpoint.SSHMetadata.Path = u.Path

	// Always use ssh scheme instead of deprecated git+ssh or ssh+git.
	endpoint.SSHMetadata.Scheme = "ssh"

	// Fallback URL for using HTTPS while still using SSH for git
	// u.Host includes host & port so can't use SSH port
	endpoint.Url = fmt.Sprintf("https://%s%s", host, u.Path)

	return endpoint
}

// EndpointFromBareSshUrl constructs a new endpoint from a bare SSH URL:
//
//   user@host.com:path/to/repo.git or
//   [user@host.com:port]:path/to/repo.git
//
// We emulate the relevant logic from Git's parse_connect_url() and
// host_end() functions in connect.c:
//   https://github.com/git/git/blob/0f828332d5ac36fc63b7d8202652efa152809856/connect.c#L673-L695
//   https://github.com/git/git/blob/0f828332d5ac36fc63b7d8202652efa152809856/connect.c#L1051
func EndpointFromBareSshUrl(s string) Endpoint {
	var userHostAndPort string
	leftToParse := s
	if i := strings.Index(s, "@["); i >= 0 {
		userHostAndPort = s[0 : i+1]
		leftToParse = s[i+1:]
	}

	var bracketed bool
	if leftToParse[0] == '[' {
		if i := strings.Index(leftToParse, "]"); i >= 0 {
			userHostAndPort += leftToParse[1:i]
			leftToParse = leftToParse[i+1:]
			bracketed = true
		}
	}

	i := strings.Index(leftToParse, ":")
	if i < 0 {
		return Endpoint{Url: s}
	}
	path := leftToParse[i+1:]
	if !bracketed {
		userHostAndPort += leftToParse[0:i]
	}

//// DEBUG: rename functions to use SSH, GitSyntax, etc.
// https://github.com/golang/go/wiki/CodeReviewComments#initialisms

//// DEBUG: rename Bare and bare in _test.go and t/* and comments above

//// DEBUG: test with IPv6 [::1] and add tests in _test.go, t/t-env, t/batch-transfer

//// Documentation/config/protocol.txt
//// t/t5603-clone-dirname.sh

//// 86ceb337ec340c7db9b060b90bfab05a08b8251b in git/git
//// 6b6c5f7a2f66751a93afce54277a1f30ab0dc521 in 2.3.7

	match := sshHostPortRE.FindStringSubmatch(userHostAndPort)
	if match == nil || len(match) < 2 {
		return Endpoint{Url: UrlUnknown}
	}

	var endpoint Endpoint
	endpoint.SSHMetadata.UserAndHost = match[1]
	if len(match) > 2 {
		endpoint.SSHMetadata.Port = match[2]
	}
	endpoint.SSHMetadata.Path = path

	// Fallback URL for using HTTPS while still using SSH for git
	host := endpoint.SSHMetadata.UserAndHost
	if i = strings.Index(host, "@"); i >= 0 {
		host = host[i+1:]
	}
	endpoint.Url = fmt.Sprintf("https://%s/%s", host, strings.TrimLeft(path, "/"))

	return endpoint
}

// Construct a new endpoint from a HTTP URL
func EndpointFromHttpUrl(u *url.URL) Endpoint {
	// just pass this straight through
	return Endpoint{Url: u.String()}
}

func EndpointFromLocalPath(path string) Endpoint {
	return Endpoint{Url: git.RewriteLocalPathAsURL(path)}
}

// Construct a new endpoint from a file URL
func EndpointFromFileUrl(u *url.URL) Endpoint {
	// just pass this straight through
	return Endpoint{Url: u.String()}
}
