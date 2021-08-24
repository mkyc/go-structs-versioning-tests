package v0

import (
	"errors"
	"github.com/epiphany-platform/e-structures/shared"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
)

type State struct {
	Meta   *Meta         `json:"meta" validate:"required"`
	Status shared.Status `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *Config       `json:"config" validate:"omitempty"`
	Output *Output       `json:"output" validate:"omitempty"`
	Unused []string      `json:"-"`
}

func (s *State) Init(moduleVersion string) {
	*s = State{
		Meta: &Meta{
			Kind:          to.StrPtr(stateKind),
			Version:       to.StrPtr(stateVersion),
			ModuleVersion: to.StrPtr(moduleVersion),
		},
		Status: shared.Initialized,
		Config: nil, // TODO should it be nil?
		Output: nil, // TODO should it be nil?
		Unused: []string{},
	}
}

func (s *State) Backup(path string) error {
	return shared.Backup(s, path)
}

func (s *State) Load(path string) error {
	i, err := shared.Load(s, path, stateVersion)
	if err != nil {
		return err
	}
	state, ok := i.(*State)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = state.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*s = *state
	return nil
}

func (s *State) Save(path string) error {
	return shared.Save(s, path)
}

func (s *State) Print() ([]byte, error) {
	return shared.Print(s)
}

func (s *State) Validate() error {
	if s == nil {
		return errors.New("expected state is nil")
	}
	validate := validator.New()
	err := validate.RegisterValidation("version", validators.HasVersion)
	if err != nil {
		return err
	}
	err = validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}

func (s *State) Upgrade(path string) error {
	i, err := shared.Upgrade(s, path)
	if err != nil {
		return err
	}
	state, ok := i.(*State)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = state.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*s = *state
	return nil
}

func (s *State) UpgradeFunc(_ map[string]interface{}) error {
	return nil
}

func (s *State) SetUnused(unused []string) {
	s.Unused = unused
}

// TODO consider validation in output ... but really think about it hard. It might not be desired.

type Output struct {
	KubeConfig *string `json:"kubeconfig"`
}
