package shared

type Initializer interface {

	// Init is responsible for building default version of structure to start from.
	Init(moduleVersion string)
}

type Backupper interface {

	// Backup is responsible for storing current version of structure in new file. It MUST fail if there is
	// file already in place pointed by path.
	Backup(path string) error // TODO this should not argument, but backup location should be unified
	// TODO add BackupRaw() method
}

type Loader interface {

	// Load is responsible for loading file pointed with path argument. It should fail in following cases:
	// - if file provided in path argument is not present
	// - if version in file is not version expected by structure
	// - if validation of loaded structure failed
	//
	// In case of incorrect version fallback to Upgrader.Upgrade method should be applied by module.
	// In case of failed validation (provided by Validator.Validate method) it should be considered
	// panic situation and usually user is forced to fix file.
	Load(path string) error
}

type Saver interface {

	// Save is responsible for storing structure into file. It must always use Validator.Validate method
	// to ensure that saved file is not corrupted. It should (but not must) use Printer.Print method
	// to produce structure JSON.
	Save(path string) error
}

type Printer interface {

	// Print is responsible for producing JSON form of structure.
	Print() ([]byte, error)
}

type Validator interface {

	// Validate is responsible for checking that structure is correct with all custom validation rules.
	Validate() error
}

type Upgrader interface {

	// Upgrade is responsible for upgrading structure to current version. It is designed to be a fallback
	// method after Loader.Load wasn't able to load structure from file and returned with NotCurrentVersionError.
	Upgrade(path string) error

	// UpgradeFunc method is responsible for delivery of structure upgrading function.
	UpgradeFunc(input map[string]interface{}) error
}

type WithUnused interface {

	// SetUnused is responsible for setting list of strings indicating that some found fields are unknown to
	// structure.
	SetUnused(unused []string)
}
