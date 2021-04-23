package helmdecoder

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strings"

	rspb "helm.sh/helm/v3/pkg/release"
	v1 "k8s.io/api/core/v1"
)

// ConvertSecretToHelmRelease converts the given kubernetes Secret to a helm release or returns nil if its not a valid secret
func ConvertSecretToHelmRelease(s *v1.Secret) (*rspb.Release, error) {
	if !strings.HasPrefix(string(s.Type), "helm.sh") {
		return nil, nil
	}
	if s.Data == nil {
		return nil, nil
	}
	val, ok := s.Data["release"]
	if !ok {
		return nil, nil
	}
	return decodeRelease(string(val))
}

// Rest of the code found below copied from:
// https://github.com/helm/helm/blob/9b42702a4bced339ff424a78ad68dd6be6e1a80a/pkg/storage/driver/util.go

var b64 = base64.StdEncoding
var magicGzip = []byte{0x1f, 0x8b, 0x08}

func decodeRelease(data string) (*rspb.Release, error) {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls rspb.Release
	// unmarshal protobuf bytes
	if err := json.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}
