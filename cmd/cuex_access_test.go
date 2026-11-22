/*
Copyright 2026 The KubeVela Authors.

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

package main

import (
	"testing"
	"time"

	cuexv1alpha1 "github.com/kubevela/pkg/apis/cue/v1alpha1"
	"github.com/kubevela/pkg/cue/cuex"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kubevela/workflow/pkg/providers"
)

// The compilers a WorkflowRun's steps run with check calls to restricted
// Packages, which the controller's scheme can read.
func TestInstallCallCheck(t *testing.T) {
	require.True(t, scheme.Recognizes(cuexv1alpha1.GroupVersion.WithKind("Package")))
	installCallCheck(fake.NewClientBuilder().WithScheme(scheme).Build(), time.Minute)
	require.NotNil(t, providers.DefaultCompiler.Get().CallCheck)
	require.NotNil(t, cuex.DefaultCompiler.Get().CallCheck)
}
