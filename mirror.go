package oppositus

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"eagain.net/go/oppositus/channels"
	"eagain.net/go/oppositus/internal/atomic"
	"eagain.net/go/oppositus/internal/href"
	"eagain.net/go/oppositus/sig"
	"eagain.net/go/oppositus/versionfile"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const defaultBaseURL = "http://release.core-os.net/amd64-usr/"

var baseURL *url.URL

func init() {
	var err error
	baseURL, err = url.Parse(defaultBaseURL)
	if err != nil {
		panic(err)
	}
}

// Option is passed to Mirror to change its behavior.
type Option option

type option func(*config) error

type config struct {
	chans  []channels.Channel
	filter func(basename string) bool
	errFn  func(error) error
}

// WithChannels sets the channels to mirror. Caller must not mutate
// chans after the call.
func WithChannels(chans ...channels.Channel) Option {
	return func(conf *config) error {
		conf.chans = chans
		return nil
	}
}

// WithFilter sets a filter files must pass, or they won't be
// mirrored.
func WithFilter(fn func(basename string) bool) Option {
	return func(conf *config) error {
		conf.filter = fn
		return nil
	}
}

// WithErrorHandler sets a function that decides which errors are
// fatal. If it returns a non-nil error, the mirroring process aborts;
// otherwise, as much progress is made as possible.
//
// A typical use would be to log errors and return nil.
func WithErrorHandler(fn func(error) error) Option {
	return func(conf *config) error {
		conf.errFn = fn
		return nil
	}
}

// Mirror fetches CoreOS releases, verifies signatures, and stores
// them locally under the directory dst.
func Mirror(ctx context.Context, dst string, opts ...Option) error {
	conf := config{
		chans: channels.All(),
	}
	for _, opt := range opts {
		if err := opt(&conf); err != nil {
			return err
		}
	}
	for _, channel := range conf.chans {
		if err := mirrorChannel(ctx, dst, conf.filter, conf.errFn, channel); err != nil {
			if err := conf.errFn(err); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func mirrorChannel(ctx context.Context, dst string, filter func(string) bool, errFn func(error) error, channel channels.Channel) error {
	chanURL := new(url.URL)
	*chanURL = *baseURL
	chanURL.Host = channel.String() + "." + chanURL.Host

	current := chanURL.ResolveReference(&url.URL{Path: "current/version.txt"})
	resp, err := ctxhttp.Get(ctx, nil, current.String())
	if err != nil {
		return fmt.Errorf("cannot fetch channel %v: %v", channel, err)
	}
	defer resp.Body.Close()

	version, err := versionfile.ParseVersionID(resp.Body)
	if err != nil {
		return err
	}

	// make separate subdir for every channel, but share versions across them

	allPath := filepath.Join(dst, "all")
	if err := os.Mkdir(allPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	verPath := filepath.Join(allPath, version)
	if err := os.Mkdir(verPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	log.Printf("channel %v is at version %v", channel, version)
	verURL := chanURL.ResolveReference(&url.URL{Path: version + "/"})
	if err := mirrorVersion(ctx, verPath, verURL, filter, errFn); err != nil {
		return err
	}

	chanPath := filepath.Join(dst, channel.String())
	if err := os.Mkdir(chanPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	// create a symlink from current to the version dir
	if err := atomic.Symlink(path.Join("..", "all", version), path.Join(dst, channel.String(), "current")); err != nil {
		return err
	}
	return nil
}

// mirrorVersion downloads the CoreOS version at u into path, while
// checking signatures.
func mirrorVersion(ctx context.Context, dst string, u *url.URL, filter func(string) bool, errFn func(error) error) error {
	log.Printf("mirroring %v", u)

	// fetch directory listing
	resp, err := ctxhttp.Get(ctx, nil, u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	hrefs := href.New(resp.Body)
	for {
		link, err := hrefs.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if err := mirrorFile(ctx, dst, u, link, filter); err != nil {
			if err := errFn(err); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func mirrorFile(ctx context.Context, dst string, u *url.URL, link string, filter func(string) bool) error {
	rel, err := url.Parse(link)
	if err != nil {
		return err
	}
	if rel.Scheme != "" || rel.Opaque != "" || rel.Host != "" || strings.HasPrefix(rel.Path, "/") {
		// skip non-relative links
		return nil
	}
	if strings.Contains(rel.Path, "/") {
		// skip links to other directories
		return nil
	}
	if strings.HasPrefix(rel.Path, ".") {
		// skip links to hidden files (we use them for our own
		// purposes) or other directories (the ".." case, without
		// slash
		return nil
	}
	if rel.RawQuery != "" || rel.Fragment != "" {
		// skip things that don't look like links to static files
		return nil
	}

	const sigExt = ".sig"
	if !strings.HasSuffix(rel.Path, sigExt) {
		// we only download signed things, so filter out
		// everything that isn't a signature
		return nil
	}
	rel.Path = rel.Path[:len(rel.Path)-len(sigExt)]

	if !filter(rel.Path) {
		return nil
	}

	// see if we have it already; files are considered immutable
	if _, err := os.Stat(path.Join(dst, rel.Path)); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		// got it already
		return nil
	}

	log.Printf("downloading %v", rel.Path)
	u2 := u.ResolveReference(rel)
	if err := sig.Download(ctx, dst, u2); err != nil {
		return err
	}
	return nil
}
