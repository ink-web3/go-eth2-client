// Copyright © 2021 Attestant Limited.
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

package altair

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// BeaconBlock represents a beacon block.
type BeaconBlock struct {
	Slot          phase0.Slot
	ProposerIndex phase0.ValidatorIndex
	ParentRoot    phase0.Root `ssz-size:"32"`
	StateRoot     phase0.Root `ssz-size:"32"`
	Body          *BeaconBlockBody
}

// beaconBlockJSON is the spec representation of the struct.
type beaconBlockJSON struct {
	Slot          string           `json:"slot"`
	ProposerIndex string           `json:"proposer_index"`
	ParentRoot    string           `json:"parent_root"`
	StateRoot     string           `json:"state_root"`
	Body          *BeaconBlockBody `json:"body"`
}

// beaconBlockYAML is the spec representation of the struct.
type beaconBlockYAML struct {
	Slot          uint64           `yaml:"slot"`
	ProposerIndex uint64           `yaml:"proposer_index"`
	ParentRoot    string           `yaml:"parent_root"`
	StateRoot     string           `yaml:"state_root"`
	Body          *BeaconBlockBody `yaml:"body"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockJSON{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    fmt.Sprintf("%#x", b.ParentRoot),
		StateRoot:     fmt.Sprintf("%#x", b.StateRoot),
		Body:          b.Body,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlock) UnmarshalJSON(input []byte) error {
	var beaconBlockJSON beaconBlockJSON
	if err := json.Unmarshal(input, &beaconBlockJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return b.unpack(&beaconBlockJSON)
}

func (b *BeaconBlock) unpack(beaconBlockJSON *beaconBlockJSON) error {
	if beaconBlockJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(beaconBlockJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = phase0.Slot(slot)
	if beaconBlockJSON.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(beaconBlockJSON.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	b.ProposerIndex = phase0.ValidatorIndex(proposerIndex)
	if beaconBlockJSON.ParentRoot == "" {
		return errors.New("parent root missing")
	}
	parentRoot, err := hex.DecodeString(strings.TrimPrefix(beaconBlockJSON.ParentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent root")
	}
	if len(parentRoot) != phase0.RootLength {
		return errors.New("incorrect length for parent root")
	}
	copy(b.ParentRoot[:], parentRoot)
	if beaconBlockJSON.StateRoot == "" {
		return errors.New("state root missing")
	}
	stateRoot, err := hex.DecodeString(strings.TrimPrefix(beaconBlockJSON.StateRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for state root")
	}
	if len(stateRoot) != phase0.RootLength {
		return errors.New("incorrect length for state root")
	}
	copy(b.StateRoot[:], stateRoot)
	if beaconBlockJSON.Body == nil {
		return errors.New("body missing")
	}
	b.Body = beaconBlockJSON.Body

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (b *BeaconBlock) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&beaconBlockYAML{
		Slot:          uint64(b.Slot),
		ProposerIndex: uint64(b.ProposerIndex),
		ParentRoot:    fmt.Sprintf("%#x", b.ParentRoot),
		StateRoot:     fmt.Sprintf("%#x", b.StateRoot),
		Body:          b.Body,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BeaconBlock) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var beaconBlockJSON beaconBlockJSON
	if err := yaml.Unmarshal(input, &beaconBlockJSON); err != nil {
		return err
	}
	return b.unpack(&beaconBlockJSON)
}

// String returns a string version of the structure.
func (b *BeaconBlock) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
