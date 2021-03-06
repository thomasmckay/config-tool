package distributedstorage

import (
	"github.com/quay/config-tool/pkg/lib/shared"
)

// Validate checks the configuration settings for this field group
func (fg *DistributedStorageFieldGroup) Validate(opts shared.Options) []shared.ValidationError {

	// Make empty errors
	errors := []shared.ValidationError{}

	// If no storage locations
	if len(fg.DistributedStorageConfig) == 0 {
		newError := shared.ValidationError{
			Tags:    []string{"DISTRIBUTED_STORAGE_CONFIG"},
			Policy:  "A is empty",
			Message: "DISTRIBUTED_STORAGE_CONFIG must contain at least one storage location.",
		}
		errors = append(errors, newError)
		return errors
	}

	for _, storageConf := range fg.DistributedStorageConfig {

		if storageConf.Name == "LocalStorage" && fg.FeatureStorageReplication {
			newError := shared.ValidationError{
				Tags:    []string{"FEATURE_STORAGE_REPLICATION"},
				Policy:  "",
				Message: "FEATURE_STORAGE_REPLICATION is not supported by LocalStorage.",
			}
			errors = append(errors, newError)
		}

		if ok, err := shared.ValidateMinioStorage(&storageConf.Args, "DistributedStorage"); !ok {
			errors = append(errors, err)
		}
	}

	// Return errors
	return errors

}
