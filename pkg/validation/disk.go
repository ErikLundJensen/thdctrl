package validation

import (
	"fmt"
	"regexp"
)

var (
	// validDiskName matches common Linux disk device names like:
	// sda, sdb1, nvme0n1, nvme0n1p1, mmcblk0, mmcblk0p1
	validDiskName = regexp.MustCompile(`^(sd[a-z][1-9]?|nvme\d+n\d+p?\d*|mmcblk\d+p?\d*)$`)
)

// ValidateDiskName checks if the provided disk name is valid
func ValidateDiskName(disk string) error {
	if !validDiskName.MatchString(disk) {
		return fmt.Errorf("invalid disk name '%s'. Must be a valid Linux device name (e.g., sda, nvme0n1)", disk)
	}
	return nil
}
