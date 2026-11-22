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

package custom

import (
	"context"
	"testing"

	"cuelang.org/go/cue"
	"github.com/stretchr/testify/require"

	oamv1alpha1 "github.com/kubevela/pkg/apis/oam/v1alpha1"
	"github.com/kubevela/pkg/cue/cuex"
	cuexruntime "github.com/kubevela/pkg/cue/cuex/runtime"
	pkgruntime "github.com/kubevela/pkg/util/runtime"

	"github.com/kubevela/workflow/api/v1alpha1"
	"github.com/kubevela/workflow/pkg/cue/process"
	"github.com/kubevela/workflow/pkg/types"
)

// A step's template is compiled as the definition it was loaded from.
func TestStepCompilesAsItsDefinition(t *testing.T) {
	var seen []cuex.Owner
	compiler := cuex.NewCompilerWithInternalPackages(
		pkgruntime.Must(cuexruntime.NewInternalPackage("whoami", `
#Ask: {
	#do:       "ask"
	#provider: "whoami"
	$returns?: {...}
}`, map[string]cuexruntime.ProviderFn{
			"ask": cuexruntime.NativeProviderFn(func(ctx context.Context, v cue.Value) (cue.Value, error) {
				owner, _ := cuex.OwnerFrom(ctx)
				seen = append(seen, owner)
				return v.FillPath(cue.ParsePath("$returns.ok"), true), nil
			}),
		})))
	owner := cuex.Owner{Kind: "WorkflowStepDefinition", Namespace: "team-a", Name: "ask",
		Writer: &cuex.Writer{User: "alice"}}
	loader := NewTaskLoader(nil, 0, process.NewContext(process.ContextData{Name: "app", Namespace: "team-a"}), compiler).
		WithOwners(func(context.Context, string) (string, *cuex.Owner, error) {
			return "import \"vela/whoami\"\nask: whoami.#Ask\n", &owner, nil
		})
	gen, err := loader.GetTaskGenerator(context.Background(), "ask")
	require.NoError(t, err)
	run, err := gen(oamv1alpha1.WorkflowStep{WorkflowStepBase: oamv1alpha1.WorkflowStepBase{Name: "ask", Type: "ask"}}, &types.TaskGeneratorOptions{})
	require.NoError(t, err)
	status, _, err := run.Run(newWorkflowContextForTest(t), &types.TaskRunOptions{})
	require.NoError(t, err)
	require.Equal(t, v1alpha1.WorkflowStepPhaseSucceeded, status.Phase, status.Message)
	require.NotEmpty(t, seen)
	require.Equal(t, owner, seen[0])
}
