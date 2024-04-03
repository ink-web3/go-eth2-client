// Copyright © 2023 Attestant Limited.
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

package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

// BeaconBlockBlobs fetches the blobs given a block ID.
func (s *Service) BeaconBlockBlobs(ctx context.Context, blockID string) ([]*deneb.BlobSidecar, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		beaconBlockBlobs, err := client.(consensusclient.BeaconBlockBlobsProvider).BeaconBlockBlobs(ctx, blockID)
		if err != nil {
			return nil, err
		}
		return beaconBlockBlobs, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.([]*deneb.BlobSidecar), nil
}
