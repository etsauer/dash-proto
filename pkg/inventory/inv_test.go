package inventory

import (
	"path/filepath"
	"testing"
)

var testInv = `
version: 3.0
args:
  - --dry-run
resource_groups:
  - name: Default Resources
    resources:
      - name: Raw Manifests
        file:
          path: manifests/
      - name: Helm Charts
        helm:
          chart: stable/redis
          valueFiles:
          - redis-vars.yaml
          values:
            replicas: 3
            labels:
              something: foo
              else: bar
  - name: OpenShift Templates
    resources:
    - name: OpenShift Templates
      openshiftTemplate:
        template: templates/app-stack.yaml
        params:
          APP_NAME: single-file-tpl
    - name: Directory of Templates
      openshiftTemplate:
        template: templates/
        params:
          APP_NAME: tpl-dir-test
        paramFiles:
        - template-params/app-stack
    - name: Directory of templates and params
      openshiftTemplate:
        template: templates/
        params:
          SOME_PARAM: ignore
        paramDir: template-params/
    - name: One template directory of params
      openshiftTemplate:
        template: templates/app-stack.yaml
        params:
          SOME_PARAM: ignore
        paramDir: template-params/
`

var testStruct = Inventory{
	Version: 3.0,
	Args:    []string{"--dry-run"},
	ResourceGroups: []ResourceGroup{
		ResourceGroup{
			Name: "Default Resources",
			Resources: []Resource{
				Resource{
					Name: "Raw Manifests",
					File: FileTemplate{
						Path: "manifests/",
					},
				},
				Resource{
					Name: "Helm Charts",
					Helm: HelmChart{
						Chart: "stable/redis",
					},
				},
			},
		},
	},
}

var i Inventory
var invPath = "../../examples/default"

func init() {
	i.Load([]byte(testInv), invPath)
}

func TestValues(t *testing.T) {
	if i.Version != 3.0 {
		t.Errorf("Wrong Version. Want \"%v\", got %v", 3.0, i.Version)
	}
	if i.Output == "" {
		t.Error("Output dir should not be empty")
	}
}

func TestPrefix(t *testing.T) {
	iPre, err := filepath.Abs(i.Prefix)
	if err != nil {
		t.Error(err)
	}
	tstPre, err := filepath.Abs(invPath)
	if err != nil {
		t.Error(err)
	}

	if filepath.Clean(iPre) != filepath.Clean(tstPre) {
		t.Errorf("Inventory prefix test failed. Want \"%s\", got %s", tstPre, iPre)
	}

	for _, rg := range i.ResourceGroups {
		rgPre, err := filepath.Abs(rg.Prefix)
		if err != nil {
			t.Error(err)
		}
		tstPre, err := filepath.Abs(invPath)
		if err != nil {
			t.Error(err)
		}

		if rgPre != tstPre {
			t.Errorf("ResourceGroup/%s prefix test failed. Want \"%s\", got %s", rg.Name, invPath, rg.Prefix)
		}
	}
}

func TestAction(t *testing.T) {
	if i.Action != "apply" {
		t.Errorf("Inventory action test failed. Want \"apply\", got %s", i.Action)
	}
}

func TestProcess(t *testing.T) {
	namespace := ""
	err := i.Process(&namespace)
	if err != nil {
		t.Errorf("Process failed, %v", err)
	}

	namespace = "dash-inv-test"
	err = i.Process(&namespace)
	if err != nil {
		t.Errorf("Process failed, %v", err)
	}
}
