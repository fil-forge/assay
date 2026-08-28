package cmdutil

import (
	"bytes"
	"fmt"

	"github.com/fil-forge/ucantone/ipld/datamodel"
	"github.com/fil-forge/ucantool/pkg/ipldfmt"
)

// FormatCBOR decodes DAG-CBOR bytes as a dynamic IPLD value and formats it as
// indented DAG-JSON (syntax highlighted when stdout is a terminal).
func FormatCBOR(buf []byte) (string, error) {
	var value datamodel.Any
	if err := value.UnmarshalCBOR(bytes.NewReader(buf)); err != nil {
		return "", fmt.Errorf("decoding DAG-CBOR: %w", err)
	}
	return ipldfmt.FormatDagJSON(&value), nil
}
