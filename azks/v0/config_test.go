package v0

import (
	"errors"
	"fmt"
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

func TestConfig_Init(t *testing.T) {
	tests := []struct {
		name          string
		moduleVersion string
		want          *Config
	}{
		{
			name:          "happy path",
			moduleVersion: "v1.1.1",
			want: &Config{
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
					AzureAd: nil,

					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got := &Config{}
			got.Init(tt.moduleVersion)
			a.Equal(tt.want, got)
		})
	}
}

func TestConfig_Backup(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr error
	}{
		{
			name:    "happy path",
			config:  &Config{},
			wantErr: nil,
		},
		{
			name:    "file already exists",
			config:  &Config{},
			wantErr: os.ErrExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p, err := createTempDirectory("azks-config-backup")
			if errors.Is(tt.wantErr, os.ErrExist) {
				err = ioutil.WriteFile(filepath.Join(p, "backup-file.json"), []byte("content"), 0644)
				t.Logf("path: %s", filepath.Join(p, "backup-file.json"))
				a.NoError(err)
			}
			err = tt.config.Backup(filepath.Join(p, "backup-file.json"))
			if tt.wantErr != nil {
				a.Error(err)
				a.Equal(tt.wantErr, err)
			} else {
				a.NoError(err)
			}
		})
	}
}

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "happy path",
			json: []byte(`{
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
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		},
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
					AzureAd: &AzureAd{
						Managed:             to.BoolPtr(true),
						TenantId:            to.StrPtr("123123123123"),
						AdminGroupObjectIds: []string{"123123123123"},
					},
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in multiple places",
			json: []byte(`{
	"meta": {
		"kind": "azksConfig",
		"version": "v0.3.0",
		"module_version": "v1.1.1"
	},
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"name": "epiphany",
		"extra_inner_field" : "extra_inner_value",
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
}`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
					AzureAd:       nil,
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{"params.extra_inner_field", "extra_outer_field"},
			},
			wantErr: nil,
		},
		{
			name: "ensure load is performing validation",
			json: []byte(`{
	"meta": {
		"kind": "azksConfig",
		"version": "v0.3.0",
		"module_version": "v0.0.1"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing azure_ad",
			json: []byte(`{
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
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
			wantErr: nil,
		},
		{
			name: "null azure_ad",
			json: []byte(`{
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
}`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDocumentFile("azks-config-load", tt.json)
			r.NoError(err)
			got := &Config{}
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
			}
		})
	}
}

func TestConfig_Save(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    []byte
		wantErr error
	}{
		{
			name: "happy path",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
					AzureAd: nil,

					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azksConfig",
		"version": "v0.3.0",
		"module_version": "v0.0.1"
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
}`),
			wantErr: nil,
		},
		{
			name:   "invalid",
			config: &Config{},
			want:   nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params",
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
			p, err := createTempDirectory("azbi-config-save")
			a.NoError(err)

			err = tt.config.Save(filepath.Join(p, "file.json"))
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

func TestConfig_Print(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    []byte
		wantErr bool
	}{
		{
			name: "happy path",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
					AzureAd: nil,

					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azksConfig",
		"version": "v0.3.0",
		"module_version": "v0.0.1"
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
}`),
			wantErr: false,
		},
		{
			name:    "invalid",
			config:  &Config{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got, err := tt.config.Print()
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
				a.Equal(string(tt.want), string(got))
			}
		})
	}
}

func TestConfig_Valid(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr error
	}{
		{
			name: "minimal correct",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
			wantErr: nil,
		},
		{
			name:   "empty struct",
			config: &Config{},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "meta missing",
			config: &Config{
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
					AzureAd: nil,

					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Name",
					Field: "Name",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.Location",
					Field: "Location",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.RsaPublicKeyPath",
					Field: "RsaPublicKeyPath",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.RgName",
					Field: "RgName",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VnetName",
					Field: "VnetName",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.SubnetName",
					Field: "SubnetName",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.KubernetesVersion",
					Field: "KubernetesVersion",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.EnableNodePublicIp",
					Field: "EnableNodePublicIp",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.EnableRbac",
					Field: "EnableRbac",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.IdentityType",
					Field: "IdentityType",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AdminUsername",
					Field: "AdminUsername",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool",
					Field: "DefaultNodePool",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile",
					Field: "AutoScalerProfile",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Name:             to.StrPtr(""),
					Location:         to.StrPtr(""),
					RsaPublicKeyPath: to.StrPtr(""),

					RgName:     to.StrPtr(""),
					VnetName:   to.StrPtr(""),
					SubnetName: to.StrPtr(""),

					KubernetesVersion:  to.StrPtr(""),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),

					DefaultNodePool: &DefaultNodePool{
						Size:        to.IntPtr(2),
						Min:         to.IntPtr(2),
						Max:         to.IntPtr(5),
						VmSize:      to.StrPtr(""),
						DiskGbSize:  to.IntPtr(36),
						AutoScaling: to.BoolPtr(true),
						Type:        to.StrPtr(""),
					},
					AutoScalerProfile: &AutoScalerProfile{
						BalanceSimilarNodeGroups:      to.BoolPtr(false),
						MaxGracefulTerminationSec:     to.StrPtr(""),
						ScaleDownDelayAfterAdd:        to.StrPtr(""),
						ScaleDownDelayAfterDelete:     to.StrPtr(""),
						ScaleDownDelayAfterFailure:    to.StrPtr(""),
						ScanInterval:                  to.StrPtr(""),
						ScaleDownUnneeded:             to.StrPtr(""),
						ScaleDownUnready:              to.StrPtr(""),
						ScaleDownUtilizationThreshold: to.StrPtr(""),
					},

					IdentityType:  to.StrPtr(""),
					AdminUsername: to.StrPtr(""),
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Name",
					Field: "Name",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.Location",
					Field: "Location",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.RsaPublicKeyPath",
					Field: "RsaPublicKeyPath",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.RgName",
					Field: "RgName",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VnetName",
					Field: "VnetName",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.SubnetName",
					Field: "SubnetName",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.KubernetesVersion",
					Field: "KubernetesVersion",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.IdentityType",
					Field: "IdentityType",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AdminUsername",
					Field: "AdminUsername",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.VmSize",
					Field: "VmSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Type",
					Field: "Type",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.MaxGracefulTerminationSec",
					Field: "MaxGracefulTerminationSec",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterAdd",
					Field: "ScaleDownDelayAfterAdd",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterDelete",
					Field: "ScaleDownDelayAfterDelete",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterFailure",
					Field: "ScaleDownDelayAfterFailure",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScanInterval",
					Field: "ScanInterval",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnneeded",
					Field: "ScaleDownUnneeded",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnready",
					Field: "ScaleDownUnready",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUtilizationThreshold",
					Field: "ScaleDownUtilizationThreshold",
					Tag:   "min",
				},
			},
		},
		{
			name: "major version mismatch",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v100.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
					AzureAd: nil,

					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		{
			name: "minor version mismatch",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.100.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
					AzureAd: nil,

					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.1.100"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
					AzureAd: nil,

					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "empty azure ad",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
					AzureAd:       &AzureAd{},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.Managed",
					Field: "Managed",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.TenantId",
					Field: "TenantId",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.AdminGroupObjectIds",
					Field: "AdminGroupObjectIds",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty azure ad params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
					AzureAd: &AzureAd{
						Managed:             to.BoolPtr(true),
						TenantId:            to.StrPtr(""),
						AdminGroupObjectIds: []string{},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.TenantId",
					Field: "TenantId",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.AdminGroupObjectIds",
					Field: "AdminGroupObjectIds",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty azure_ad.admin_group_object_ids element",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
					AzureAd: &AzureAd{
						Managed:             to.BoolPtr(true),
						TenantId:            to.StrPtr("123123123123"),
						AdminGroupObjectIds: []string{""},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.AdminGroupObjectIds[0]",
					Field: "AdminGroupObjectIds[0]",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing auto_scaler_profile",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile",
					Field: "AutoScalerProfile",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty auto_scaler_profile aka missing params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
					AutoScalerProfile: &AutoScalerProfile{},
					IdentityType:      to.StrPtr("SystemAssigned"),
					AdminUsername:     to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.BalanceSimilarNodeGroups",
					Field: "BalanceSimilarNodeGroups",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.MaxGracefulTerminationSec",
					Field: "MaxGracefulTerminationSec",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterAdd",
					Field: "ScaleDownDelayAfterAdd",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterDelete",
					Field: "ScaleDownDelayAfterDelete",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterFailure",
					Field: "ScaleDownDelayAfterFailure",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScanInterval",
					Field: "ScanInterval",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnneeded",
					Field: "ScaleDownUnneeded",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnready",
					Field: "ScaleDownUnready",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUtilizationThreshold",
					Field: "ScaleDownUtilizationThreshold",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty auto_scaler_profile params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
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
						MaxGracefulTerminationSec:     to.StrPtr(""),
						ScaleDownDelayAfterAdd:        to.StrPtr(""),
						ScaleDownDelayAfterDelete:     to.StrPtr(""),
						ScaleDownDelayAfterFailure:    to.StrPtr(""),
						ScanInterval:                  to.StrPtr(""),
						ScaleDownUnneeded:             to.StrPtr(""),
						ScaleDownUnready:              to.StrPtr(""),
						ScaleDownUtilizationThreshold: to.StrPtr(""),
					},
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.MaxGracefulTerminationSec",
					Field: "MaxGracefulTerminationSec",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterAdd",
					Field: "ScaleDownDelayAfterAdd",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterDelete",
					Field: "ScaleDownDelayAfterDelete",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterFailure",
					Field: "ScaleDownDelayAfterFailure",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScanInterval",
					Field: "ScanInterval",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnneeded",
					Field: "ScaleDownUnneeded",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnready",
					Field: "ScaleDownUnready",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUtilizationThreshold",
					Field: "ScaleDownUtilizationThreshold",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing default_node_pool",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool",
					Field: "DefaultNodePool",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty default_node_pool aka missing params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
					DefaultNodePool:    &DefaultNodePool{},
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Min",
					Field: "Min",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.VmSize",
					Field: "VmSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.DiskGbSize",
					Field: "DiskGbSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.AutoScaling",
					Field: "AutoScaling",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Type",
					Field: "Type",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty default_node_pool params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
					DefaultNodePool: &DefaultNodePool{
						Size:        to.IntPtr(2),
						Min:         to.IntPtr(2),
						Max:         to.IntPtr(5),
						VmSize:      to.StrPtr(""),
						DiskGbSize:  to.IntPtr(0),
						AutoScaling: to.BoolPtr(true),
						Type:        to.StrPtr(""),
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.VmSize",
					Field: "VmSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.DiskGbSize",
					Field: "DiskGbSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Type",
					Field: "Type",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing default_node_pool.min",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
						VmSize: to.StrPtr("Standard_DS2_v2"),
						Type:   to.StrPtr("VirtualMachineScaleSets"),
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Min",
					Field: "Min",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.DiskGbSize",
					Field: "DiskGbSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.AutoScaling",
					Field: "AutoScaling",
					Tag:   "required",
				},
			},
		},
		{
			name: "default_node_pool min > max",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
						Size:        to.IntPtr(3),
						Min:         to.IntPtr(3),
						Max:         to.IntPtr(2),
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "gtefield",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "ltefield",
				},
			},
		},
		{
			name: "default_node_pool size < min",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
						Size:        to.IntPtr(1),
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "gtefield",
				},
			},
		},
		{
			name: "default_node_pool size > max",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
						Size:        to.IntPtr(6),
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "ltefield",
				},
			},
		},
		{
			name: "default_node_pool negative sizes",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azksConfig"),
					Version:       to.StrPtr("v0.3.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
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
						Size:        to.IntPtr(-1),
						Min:         to.IntPtr(-1),
						Max:         to.IntPtr(-1),
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
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Min",
					Field: "Min",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "min",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			err := tt.config.Validate()
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

func createTempDocumentFile(name string, document []byte) (string, error) {
	p, err := ioutil.TempDir("", fmt.Sprintf("e-structures-%s-*", name))
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(p, "file.json"), document, 0644)
	return filepath.Join(p, "file.json"), err
}

func createTempDirectory(name string) (string, error) {
	return ioutil.TempDir("", fmt.Sprintf("e-structures-%s-*", name))
}
