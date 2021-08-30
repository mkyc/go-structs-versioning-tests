package load

import (
	"io/ioutil"
	"os"

	awsbi "github.com/epiphany-platform/e-structures/awsbi/v0"
	hi "github.com/epiphany-platform/e-structures/hi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
)

func State(path string) (*st.State, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return st.NewState(), nil
	} else {
		state := &st.State{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		// TODO after issue https://github.com/epiphany-platform/e-structures/issues/10 is solved
		// TODO this should be changed back to err = state.Unmarshal(bytes)
		err = state.UnmarshalDoNotUse(bytes)
		if err != nil {
			return nil, err
		}

		// TODO temporary code because of before mentioned issue
		if state.GetAzKSState() != nil && state.GetAzKSState().Status == "" {
			state.AzKS = nil
		}
		if state.GetHiState() != nil && state.GetHiState().Status == "" {
			state.Hi = nil
		}
		err = state.IsValidDoNotUse()
		if err != nil {
			return nil, err
		}
		// TODO end of temporary code

		return state, nil
	}
}

func HiConfig(path string) (*hi.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return hi.NewConfig(), nil
	} else {
		config := &hi.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}

func AwsBIConfig(path string) (*awsbi.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return awsbi.NewConfig(), nil
	} else {
		config := &awsbi.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}
