//go:build !linux

/*
 * Copyright 2025 BoxLite AI (originally Daytona Platforms Inc.
 * SPDX-License-Identifier: AGPL-3.0
 */

package diskresize

import "context"

// defaultReadDeviceSize always reports the device as missing off Linux: the
// GCP data disk and its sysfs entry only exist on the Linux runner hosts.
func defaultReadDeviceSize(ctx context.Context, devicePath string) (uint64, error) {
	return 0, ErrDeviceMissing
}

// defaultRunResize never runs off Linux for the same reason; the ticker stops
// at the missing-device skip before it can be reached.
func defaultRunResize(ctx context.Context, devicePath string) error {
	return ErrDeviceMissing
}
