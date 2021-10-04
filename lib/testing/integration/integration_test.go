/*
Copyright 2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	IntegrationSetup
}

func TestIntegration(t *testing.T) { suite.Run(t, &IntegrationSuite{}) }

func (s *IntegrationSuite) TestVersion() {
	t := s.T()

	versionMin, err := version.NewVersion("v6.2.7")
	require.NoError(t, err)
	versionMax, err := version.NewVersion("v8")
	require.NoError(t, err)

	assert.True(t, s.integration.Version().GreaterThanOrEqual(versionMin))
	assert.True(t, s.integration.Version().LessThan(versionMax))
}