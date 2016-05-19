package sig

import (
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// Download fetches the URL and the corresponding *.sig signature, and
// creates files under dst with matching basenames if the signature is
// good.
func Download(ctx context.Context, dst string, u *url.URL) error {
	sigURL := new(url.URL)
	*sigURL = *u
	sigURL.Path += ".sig"

	sigFile, err := ioutil.TempFile(dst, "."+path.Base(sigURL.Path)+".tmp.")
	if err != nil {
		return err
	}
	defer func() {
		if sigFile != nil {
			if err := os.Remove(sigFile.Name()); err != nil {
				log.Printf("cannot clean up temp file: %v", err)
			}
			if err := sigFile.Close(); err != nil {
				log.Printf("cannot close temp file: %v", err)
			}
		}
	}()
	sigResp, err := ctxhttp.Get(ctx, nil, sigURL.String())
	if err != nil {
		return err
	}
	defer sigResp.Body.Close()
	signature := io.TeeReader(sigResp.Body, sigFile)

	mainFile, err := ioutil.TempFile(dst, "."+path.Base(u.Path)+".tmp.")
	if err != nil {
		return err
	}
	defer func() {
		if mainFile != nil {
			if err := os.Remove(mainFile.Name()); err != nil {
				log.Printf("cannot clean up temp file: %v", err)
			}
			if err := mainFile.Close(); err != nil {
				log.Printf("cannot close temp file: %v", err)
			}
		}
	}()
	mainResp, err := ctxhttp.Get(ctx, nil, u.String())
	if err != nil {
		return err
	}
	defer mainResp.Body.Close()
	signed := io.TeeReader(mainResp.Body, mainFile)

	if err := Check(signed, signature); err != nil {
		return err
	}

	if err := sigFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(sigFile.Name(), path.Join(dst, path.Base(sigURL.Path))); err != nil {
		return err
	}
	sigFile = nil

	if err := mainFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(mainFile.Name(), path.Join(dst, path.Base(u.Path))); err != nil {
		return err
	}
	mainFile = nil

	return nil
}
