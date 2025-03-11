package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDiskName(t *testing.T) {
	tests := []struct {
		name    string
		disk    string
		wantErr bool
	}{
		{"valid sda", "sda", false},
		{"valid sdb1", "sdb1", false},
		{"valid nvme0n1", "nvme0n1", false},
		{"valid nvme0n1p1", "nvme0n1p1", false},
		{"valid mmcblk0", "mmcblk0", false},
		{"valid mmcblk0p1", "mmcblk0p1", false},
		{"invalid path", "/dev/sda", true},
		{"invalid special chars", "sda;rm", true},
		{"invalid empty", "", true},
		{"invalid format", "disk1", true},
		{"invalid command injection", "$(rm -rf /)", true},
		{"invalid path traversal", "../etc/passwd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDiskName(tt.disk)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
