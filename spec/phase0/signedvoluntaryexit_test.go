// Copyright © 2020, 2021 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package phase0_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestSignedVoluntaryExitJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.signedVoluntaryExitJSON",
		},
		{
			name:  "MessageMissing",
			input: []byte(`{"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "message missing",
		},
		{
			name:  "MessageWrongType",
			input: []byte(`{"message":true,"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.voluntaryExitJSON",
		},
		{
			name:  "MessageInvalid",
			input: []byte(`{"message":{"epoch":true},"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go struct field voluntaryExitJSON.epoch of type string",
		},
		{
			name:  "SignatureMissing",
			input: []byte(`{"message":{"epoch":"1","validator_index":"2"}}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureWrongType",
			input: []byte(`{"message":{"epoch":"1","validator_index":"2"},"signature":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field signedVoluntaryExitJSON.signature of type string",
		},
		{
			name:  "SignatureInvalid",
			input: []byte(`{"message":{"epoch":"1","validator_index":"2"},"signature":"invalid"}`),
			err:   "invalid value for signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureShort",
			input: []byte(`{"message":{"epoch":"1","validator_index":"2"},"signature":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureLong",
			input: []byte(`{"message":{"epoch":"1","validator_index":"2"},"signature":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "Good",
			input: []byte(`{"message":{"epoch":"1","validator_index":"2"},"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.SignedVoluntaryExit
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestSignedVoluntaryExitYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{message: {epoch: 1, validator_index: 2}, signature: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.SignedVoluntaryExit
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestSignedVoluntaryExitSpec(t *testing.T) {
	if os.Getenv("ETH2_SPEC_TESTS_DIR") == "" {
		t.Skip("ETH2_SPEC_TESTS_DIR not suppplied, not running spec tests")
	}
	baseDir := filepath.Join(os.Getenv("ETH2_SPEC_TESTS_DIR"), "tests", "mainnet", "phase0", "ssz_static", "SignedVoluntaryExit", "ssz_random")
	require.NoError(t, filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if path == baseDir {
			// Only interested in subdirectories.
			return nil
		}
		require.NoError(t, err)
		if info.IsDir() {
			t.Run(info.Name(), func(t *testing.T) {
				specYAML, err := os.ReadFile(filepath.Join(path, "value.yaml"))
				require.NoError(t, err)
				var res phase0.SignedVoluntaryExit
				require.NoError(t, yaml.Unmarshal(specYAML, &res))

				compressedSpecSSZ, err := os.ReadFile(filepath.Join(path, "serialized.ssz_snappy"))
				require.NoError(t, err)
				var specSSZ []byte
				specSSZ, err = snappy.Decode(specSSZ, compressedSpecSSZ)
				require.NoError(t, err)

				unmarshalled := &phase0.SignedVoluntaryExit{}
				require.NoError(t, unmarshalled.UnmarshalSSZ(specSSZ))
				remarshalled, err := unmarshalled.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, remarshalled)

				ssz, err := res.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, ssz)

				root, err := res.HashTreeRoot()
				require.NoError(t, err)
				rootsYAML, err := os.ReadFile(filepath.Join(path, "roots.yaml"))
				require.NoError(t, err)
				require.Equal(t, string(rootsYAML), fmt.Sprintf("{root: '%#x'}\n", root))
			})
		}
		return nil
	}))
}
