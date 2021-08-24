package v0

import (
	"errors"
	"github.com/epiphany-platform/e-structures/shared"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Meta   *Meta    `json:"meta" validate:"required"`
	Params *Params  `json:"params" validate:"required"`
	Unused []string `json:"-"`
}

func (c *Config) Init(moduleVersion string) {
	*c = Config{
		Meta: &Meta{
			Kind:          to.StrPtr(configKind),
			Version:       to.StrPtr(configVersion),
			ModuleVersion: to.StrPtr(moduleVersion),
		},
		Params: &Params{
			Name:             to.StrPtr("epiphany"),
			Location:         to.StrPtr("northeurope"), //TODO possibly delete this value in future
			RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),

			RgName:     to.StrPtr("epiphany-rg"),
			VnetName:   to.StrPtr("epiphany-vnet"),
			SubnetName: to.StrPtr("azks"),

			KubernetesVersion:  to.StrPtr("1.18.14"), //TODO this is outdated version
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
	}
	// TODO consider if we should call Validate() here
}

func (c *Config) Backup(path string) error {
	return shared.Backup(c, path)
}

func (c *Config) Load(path string) error {
	i, err := shared.Load(c, path, configVersion)
	if err != nil {
		return err
	}
	config, ok := i.(*Config)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = config.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*c = *config
	return nil
}

func (c *Config) Save(path string) error {
	return shared.Save(c, path)
}

func (c *Config) Print() ([]byte, error) {
	return shared.Print(c)
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("expected config is nil")
	}
	validate := validator.New()

	err := validate.RegisterValidation("version", validators.HasVersion)
	if err != nil {
		return err
	}

	err = validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}

func (c *Config) Upgrade(path string) error {
	i, err := shared.Upgrade(c, path)
	if err != nil {
		return err
	}
	config, ok := i.(*Config)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = config.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*c = *config
	return nil
}

func (c *Config) UpgradeFunc(_ map[string]interface{}) error {
	return nil //TODO implement with first upgrade operation
}

func (c *Config) SetUnused(unused []string) {
	c.Unused = unused
}

type Meta struct {
	Kind          *string `json:"kind" validate:"required,eq=azksConfig|eq=azksState"`
	Version       *string `json:"version" validate:"required,version=~0"`
	ModuleVersion *string `json:"module_version" validate:"required"`
}

type Params struct {
	Name               *string            `json:"name" validate:"required,min=1"`
	Location           *string            `json:"location" validate:"required,min=1"`
	RsaPublicKeyPath   *string            `json:"rsa_pub_path" validate:"required,min=1"`
	RgName             *string            `json:"rg_name" validate:"required,min=1"`
	VnetName           *string            `json:"vnet_name" validate:"required,min=1"`
	SubnetName         *string            `json:"subnet_name" validate:"required,min=1"`
	KubernetesVersion  *string            `json:"kubernetes_version" validate:"required,min=1"`
	EnableNodePublicIp *bool              `json:"enable_node_public_ip" validate:"required"`
	EnableRbac         *bool              `json:"enable_rbac" validate:"required"`
	DefaultNodePool    *DefaultNodePool   `json:"default_node_pool" validate:"required,dive"`
	AutoScalerProfile  *AutoScalerProfile `json:"auto_scaler_profile" validate:"required,dive"`
	AzureAd            *AzureAd           `json:"azure_ad" validate:"omitempty"`
	IdentityType       *string            `json:"identity_type" validate:"required,min=1"`
	AdminUsername      *string            `json:"admin_username" validate:"required,min=1"`
}

type DefaultNodePool struct {
	Size        *int    `json:"size" validate:"required,min=0,gtefield=Min,ltefield=Max"`
	Min         *int    `json:"min" validate:"required,min=0"`
	Max         *int    `json:"max" validate:"required,min=0,gtefield=Min"`
	VmSize      *string `json:"vm_size" validate:"required,min=1"`
	DiskGbSize  *int    `json:"disk_gb_size" validate:"required,min=1"`
	AutoScaling *bool   `json:"auto_scaling" validate:"required"`
	Type        *string `json:"type" validate:"required,min=1"`
}

type AutoScalerProfile struct { //TODO consider changing types of string values here to make it more golang'ish
	BalanceSimilarNodeGroups      *bool   `json:"balance_similar_node_groups" validate:"required"`
	MaxGracefulTerminationSec     *string `json:"max_graceful_termination_sec" validate:"required,min=1"`
	ScaleDownDelayAfterAdd        *string `json:"scale_down_delay_after_add" validate:"required,min=1"`
	ScaleDownDelayAfterDelete     *string `json:"scale_down_delay_after_delete" validate:"required,min=1"`
	ScaleDownDelayAfterFailure    *string `json:"scale_down_delay_after_failure" validate:"required,min=1"`
	ScanInterval                  *string `json:"scan_interval" validate:"required,min=1"`
	ScaleDownUnneeded             *string `json:"scale_down_unneeded" validate:"required,min=1"`
	ScaleDownUnready              *string `json:"scale_down_unready" validate:"required,min=1"`
	ScaleDownUtilizationThreshold *string `json:"scale_down_utilization_threshold" validate:"required,min=1"`
}

type AzureAd struct {
	Managed             *bool    `json:"managed" validate:"required"`
	TenantId            *string  `json:"tenant_id" validate:"required,min=1"`
	AdminGroupObjectIds []string `json:"admin_group_object_ids" validate:"required,min=1,dive,required,min=1"`
}
