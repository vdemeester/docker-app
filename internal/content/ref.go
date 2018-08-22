package content

import (
	"crypto/sha256"
	"fmt"

	"github.com/docker/distribution/reference"
)

type appRef struct {
	reference reference.Named
}

func (r appRef) Sum() string {
	sh := sha256.New()
	sh.Write([]byte(r.String()))
	return fmt.Sprintf("%x", sh.Sum(nil))
}

func (r appRef) String() string {
	return r.reference.String()
}

func referenceSum(ref string) (appRef, error) {
	r, err := reference.ParseNormalizedNamed(ref)
	if err != nil {
		return appRef{}, err
	}
	return appRef{
		reference: r,
	}, nil
}
