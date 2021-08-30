package v0

import (
	"errors"
	"github.com/epiphany-platform/e-structures/shared"
	"github.com/epiphany-platform/e-structures/utils/test"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestState_Init(t *testing.T) {
	tests := []struct {
		name          string
		moduleVersion string
		want          *State
	}{
		{
			name:          "happy path",
			moduleVersion: "v1.1.1",
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got := &State{}
			got.Init(tt.moduleVersion)
			a.Equal(tt.want, got)
		})
	}
}

func TestState_Backup(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr error
	}{
		{
			name:    "happy path",
			state:   &State{},
			wantErr: nil,
		},
		{
			name:    "file already exists",
			state:   &State{},
			wantErr: os.ErrExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p, err := createTempDirectory("azks-state-backup")
			if errors.Is(tt.wantErr, os.ErrExist) {
				err = ioutil.WriteFile(filepath.Join(p, "backup-file.json"), []byte("content"), 0644)
				t.Logf("path: %s", filepath.Join(p, "backup-file.json"))
				a.NoError(err)
			}
			err = tt.state.Backup(filepath.Join(p, "backup-file.json"))
			if tt.wantErr != nil {
				a.Error(err)
				a.Equal(tt.wantErr, err)
			} else {
				a.NoError(err)
			}
		})
	}
}

func TestState_Load(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *State
		wantErr error
	}{
		{
			name: "happy path",
			json: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.1",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in multiple places",
			json: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.1",
		"module_version": "v1.1.1"
	},
	"status": "initialized",
	"unknown_key_1": "unknown_value_1",
	"config": {
		"meta": {
			"kind": "azksConfig",
			"version": "v0.3.0",
			"module_version": "v1.1.1"
		},
		"params": {
			"name": "epiphany",
			"unknown_key_2" : "unknown_value_2",
			"location": "northeurope",
			"rsa_pub_path": "/shared/vms_rsa.pub",
			"rg_name": "epiphany-rg",
			"vnet_name": "epiphany-vnet",
			"subnet_name": "azks",
			"kubernetes_version": "1.18.14",
			"enable_node_public_ip": false,
			"enable_rbac": false,
			"default_node_pool": {
				"size": 2,
				"min": 2,
				"max": 5,
				"vm_size": "Standard_DS2_v2",
				"disk_gb_size": 36,
				"auto_scaling": true,
				"type": "VirtualMachineScaleSets"
			},
			"auto_scaler_profile": {
				"unknown_key_3": "unknown_value_3",
				"balance_similar_node_groups": false,
				"max_graceful_termination_sec": "600",
				"scale_down_delay_after_add": "10m",
				"scale_down_delay_after_delete": "10s",
				"scale_down_delay_after_failure": "10m",
				"scan_interval": "10s",
				"scale_down_unneeded": "10m",
				"scale_down_unready": "10m",
				"scale_down_utilization_threshold": "0.5"
			},
			"azure_ad": null,
			"identity_type": "SystemAssigned",
			"admin_username": "operations"
		}
	},
	"output": {
		"unknown_key_4": "unknown_value_4"
	}
}`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{
					Meta: &Meta{
						Kind:          to.StrPtr("azksConfig"),
						Version:       to.StrPtr("v0.3.0"),
						ModuleVersion: to.StrPtr("v1.1.1"),
					},
					Params: &Params{
						Name:             to.StrPtr("epiphany"),
						Location:         to.StrPtr("northeurope"),
						RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),

						RgName:     to.StrPtr("epiphany-rg"),
						VnetName:   to.StrPtr("epiphany-vnet"),
						SubnetName: to.StrPtr("azks"),

						KubernetesVersion:  to.StrPtr("1.18.14"),
						EnableNodePublicIp: to.BoolPtr(false),
						EnableRbac:         to.BoolPtr(false),

						DefaultNodePool: &DefaultNodePool{
							Size:        to.IntPtr(2),
							Min:         to.IntPtr(2),
							Max:         to.IntPtr(5),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							DiskGbSize:  to.IntPtr(36),
							AutoScaling: to.BoolPtr(true),
							Type:        to.StrPtr("VirtualMachineScaleSets"),
						},
						AutoScalerProfile: &AutoScalerProfile{
							BalanceSimilarNodeGroups:      to.BoolPtr(false),
							MaxGracefulTerminationSec:     to.StrPtr("600"),
							ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
							ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
							ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
							ScanInterval:                  to.StrPtr("10s"),
							ScaleDownUnneeded:             to.StrPtr("10m"),
							ScaleDownUnready:              to.StrPtr("10m"),
							ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
						},
						IdentityType:  to.StrPtr("SystemAssigned"),
						AdminUsername: to.StrPtr("operations"),
					},
				},
				Output: &Output{},
				Unused: []string{
					"config.params.auto_scaler_profile.unknown_key_3",
					"config.params.unknown_key_2",
					"output.unknown_key_4",
					"unknown_key_1",
				},
			},
			wantErr: nil,
		},
		{
			name: "minimal state",
			json: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.1",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "future version mismatch",
			json: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v100.0.0",
		"module_version": "dev"
	},
	"status": "initialized"
}
`),
			want:    nil,
			wantErr: shared.NotCurrentVersionError{Version: "v100.0.0"},
		},
		{
			name: "old version",
			json: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.0",
		"module_version": "dev"
	},
	"status": "initialized"
}
`),
			want:    nil,
			wantErr: shared.NotCurrentVersionError{Version: "v0.0.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDocumentFile("azks-state-load", tt.json)
			r.NoError(err)
			got := &State{}
			err = got.Load(p)
			if tt.wantErr != nil {
				r.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
					a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))
				} else {
					a.Equal(tt.wantErr, err)
				}
			} else {
				a.NoError(err)
				wj, err2 := tt.want.Print()
				a.NoError(err2)
				gj, err2 := got.Print()
				a.NoError(err2)
				a.Equal(string(wj), string(gj))
				a.Equal(tt.want, got)
			}
		})
	}
}

func TestState_Save(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		want    []byte
		wantErr error
	}{
		{
			name: "happy path",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.1",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}`),
			wantErr: nil,
		},
		{
			name:  "invalid",
			state: &State{},
			want:  nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDirectory("azks-state-save")
			a.NoError(err)

			err = tt.state.Save(filepath.Join(p, "file.json"))
			if tt.wantErr != nil {
				a.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
					a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))
				} else {
					a.Equal(tt.wantErr, err)
				}
			} else {
				a.NoError(err)
				a.FileExists(filepath.Join(p, "file.json"))
				got, err2 := ioutil.ReadFile(filepath.Join(p, "file.json"))
				a.NoError(err2)
				a.Equal(string(tt.want), string(got))
			}
		})
	}
}

func TestState_Print(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		want    []byte
		wantErr error
	}{
		{
			name: "happy path",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{
					Meta: &Meta{
						Kind:          to.StrPtr("azksConfig"),
						Version:       to.StrPtr("v0.3.0"),
						ModuleVersion: to.StrPtr("v1.1.1"),
					},
					Params: &Params{
						Name:             to.StrPtr("epiphany"),
						Location:         to.StrPtr("northeurope"),
						RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),

						RgName:     to.StrPtr("epiphany-rg"),
						VnetName:   to.StrPtr("epiphany-vnet"),
						SubnetName: to.StrPtr("azks"),

						KubernetesVersion:  to.StrPtr("1.18.14"),
						EnableNodePublicIp: to.BoolPtr(false),
						EnableRbac:         to.BoolPtr(false),

						DefaultNodePool: &DefaultNodePool{
							Size:        to.IntPtr(2),
							Min:         to.IntPtr(2),
							Max:         to.IntPtr(5),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							DiskGbSize:  to.IntPtr(36),
							AutoScaling: to.BoolPtr(true),
							Type:        to.StrPtr("VirtualMachineScaleSets"),
						},
						AutoScalerProfile: &AutoScalerProfile{
							BalanceSimilarNodeGroups:      to.BoolPtr(false),
							MaxGracefulTerminationSec:     to.StrPtr("600"),
							ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
							ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
							ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
							ScanInterval:                  to.StrPtr("10s"),
							ScaleDownUnneeded:             to.StrPtr("10m"),
							ScaleDownUnready:              to.StrPtr("10m"),
							ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
						},
						IdentityType:  to.StrPtr("SystemAssigned"),
						AdminUsername: to.StrPtr("operations"),
					},
					Unused: []string{},
				},
				Output: &Output{
					KubeConfig: to.StrPtr("some kube config value"),
				},
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.1",
		"module_version": "v1.1.1"
	},
	"status": "initialized",
	"config": {
		"meta": {
			"kind": "azksConfig",
			"version": "v0.3.0",
			"module_version": "v1.1.1"
		},
		"params": {
			"name": "epiphany",
			"location": "northeurope",
			"rsa_pub_path": "/shared/vms_rsa.pub",
			"rg_name": "epiphany-rg",
			"vnet_name": "epiphany-vnet",
			"subnet_name": "azks",
			"kubernetes_version": "1.18.14",
			"enable_node_public_ip": false,
			"enable_rbac": false,
			"default_node_pool": {
				"size": 2,
				"min": 2,
				"max": 5,
				"vm_size": "Standard_DS2_v2",
				"disk_gb_size": 36,
				"auto_scaling": true,
				"type": "VirtualMachineScaleSets"
			},
			"auto_scaler_profile": {
				"balance_similar_node_groups": false,
				"max_graceful_termination_sec": "600",
				"scale_down_delay_after_add": "10m",
				"scale_down_delay_after_delete": "10s",
				"scale_down_delay_after_failure": "10m",
				"scan_interval": "10s",
				"scale_down_unneeded": "10m",
				"scale_down_unready": "10m",
				"scale_down_utilization_threshold": "0.5"
			},
			"azure_ad": null,
			"identity_type": "SystemAssigned",
			"admin_username": "operations"
		}
	},
	"output": {
		"kubeconfig": "some kube config value"
	}
}`),
			wantErr: nil,
		},
		{
			name:  "invalid",
			state: &State{},
			want:  nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			got, err := tt.state.Print()
			if tt.wantErr != nil {
				a.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
					a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))
				} else {
					a.Equal(tt.wantErr, err)
				}
			} else {
				a.NoError(err)
				a.Equal(string(tt.want), string(got))
			}
		})
	}
}

func TestState_Valid(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr error
	}{
		{
			name: "minimal correct",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name:  "empty struct",
			state: &State{},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "required",
				},
			},
		},
		{
			name: "meta missing",
			state: &State{
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
			},
		},
		{
			name: "major version mismatch",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v100.0.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		{
			name: "minor version mismatch",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.100.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.100"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "incorrect status",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: "incorrect",
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "eq=initialized|eq=applied|eq=destroyed",
				},
			},
		},
		{
			name: "empty config and output",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{},
				Output: &Output{},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			err := tt.state.Validate()
			if tt.wantErr != nil {
				r.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				_, ok = err.(validator.ValidationErrors)
				r.Equal(true, ok)
				errs := err.(validator.ValidationErrors)
				a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))

				for _, e := range errs {
					found := false
					for _, we := range tt.wantErr.(test.TestValidationErrors) {
						if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
					}
				}
			} else {
				a.NoError(err)
			}
		})
	}
}

func TestState_Upgrade(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *State
		wantErr error
	}{
		{
			name: "happy path nothing to upgrade state without config",
			json: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.1",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}
`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "happy path nothing to upgrade state with config and output",
			json: []byte(`{
	"meta": {
		"kind": "azksState",
		"version": "v0.0.1",
		"module_version": "v1.1.1"
	},
	"status": "initialized",
	"config": {
		"meta": {
			"kind": "azksConfig",
			"version": "v0.3.0",
			"module_version": "v1.1.1"
		},
		"params": {
			"name": "epiphany",
			"location": "northeurope",
			"rsa_pub_path": "/shared/vms_rsa.pub",
			"rg_name": "epiphany-rg",
			"vnet_name": "epiphany-vnet",
			"subnet_name": "azks",
			"kubernetes_version": "1.18.14",
			"enable_node_public_ip": false,
			"enable_rbac": false,
			"default_node_pool": {
				"size": 2,
				"min": 2,
				"max": 5,
				"vm_size": "Standard_DS2_v2",
				"disk_gb_size": 36,
				"auto_scaling": true,
				"type": "VirtualMachineScaleSets"
			},
			"auto_scaler_profile": {
				"balance_similar_node_groups": false,
				"max_graceful_termination_sec": "600",
				"scale_down_delay_after_add": "10m",
				"scale_down_delay_after_delete": "10s",
				"scale_down_delay_after_failure": "10m",
				"scan_interval": "10s",
				"scale_down_unneeded": "10m",
				"scale_down_unready": "10m",
				"scale_down_utilization_threshold": "0.5"
			},
			"azure_ad": null,
			"identity_type": "SystemAssigned",
			"admin_username": "operations"
		}
	},
	"output": {
		"kubeconfig": "some kube config value"
	}
}
`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azksState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{
					Meta: &Meta{
						Kind:          to.StrPtr("azksConfig"),
						Version:       to.StrPtr("v0.3.0"),
						ModuleVersion: to.StrPtr("v1.1.1"),
					},
					Params: &Params{
						Name:             to.StrPtr("epiphany"),
						Location:         to.StrPtr("northeurope"),
						RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),

						RgName:     to.StrPtr("epiphany-rg"),
						VnetName:   to.StrPtr("epiphany-vnet"),
						SubnetName: to.StrPtr("azks"),

						KubernetesVersion:  to.StrPtr("1.18.14"),
						EnableNodePublicIp: to.BoolPtr(false),
						EnableRbac:         to.BoolPtr(false),

						DefaultNodePool: &DefaultNodePool{
							Size:        to.IntPtr(2),
							Min:         to.IntPtr(2),
							Max:         to.IntPtr(5),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							DiskGbSize:  to.IntPtr(36),
							AutoScaling: to.BoolPtr(true),
							Type:        to.StrPtr("VirtualMachineScaleSets"),
						},
						AutoScalerProfile: &AutoScalerProfile{
							BalanceSimilarNodeGroups:      to.BoolPtr(false),
							MaxGracefulTerminationSec:     to.StrPtr("600"),
							ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
							ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
							ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
							ScanInterval:                  to.StrPtr("10s"),
							ScaleDownUnneeded:             to.StrPtr("10m"),
							ScaleDownUnready:              to.StrPtr("10m"),
							ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
						},
						IdentityType:  to.StrPtr("SystemAssigned"),
						AdminUsername: to.StrPtr("operations"),
					},
					Unused: []string{},
				},
				Output: &Output{
					KubeConfig: to.StrPtr("some kube config value"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "ensure that validation is also performed in upgrade",
			json: []byte(`{
	"meta": {
		"version": "v0.0.1",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}
`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta.Kind",
					Field: "Kind",
					Tag:   "required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDocumentFile("azbi-state-load", tt.json)
			r.NoError(err)
			got := &State{}
			err = got.Upgrade(p)
			if tt.wantErr != nil {
				r.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
					a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))
				} else {
					a.Equal(tt.wantErr, err)
				}
			} else {
				a.NoError(err)
				wj, err2 := tt.want.Print()
				a.NoError(err2)
				gj, err2 := got.Print()
				a.NoError(err2)
				a.Equal(string(wj), string(gj))
			}
		})
	}
}
