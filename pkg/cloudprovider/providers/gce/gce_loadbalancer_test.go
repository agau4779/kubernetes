/*
Copyright 2017 The Kubernetes Authors.

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

package gce

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureLoadBalancerCreatesExternalLb(t *testing.T) {
	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLbApiService()
	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	assert.NoError(t, err)
	assert.NotEmpty(t, status.Ingress)
	assertExternalLbResources(t, gce, apiService, vals, nodeNames)
}

func TestEnsureLoadBalancerCreatesInternalLb(t *testing.T) {
	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLbApiService()
	annotations := make(map[string]string)
	annotations[ServiceAnnotationLoadBalancerType] = string(LBTypeInternal)
	apiService.Annotations = annotations

	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	assert.NoError(t, err)
	assert.NotEmpty(t, status.Ingress)
	assertInternalLbResources(t, gce, apiService, vals, nodeNames)
}
