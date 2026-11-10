/*
 * Copyright 2026 BoxLite AI
 * SPDX-License-Identifier: AGPL-3.0
 */

/**
 * Identifies the cloud data disk backing a runner's `BOXLITE_HOME` filesystem,
 * as reported by the runner itself in each healthcheck. Only runners whose
 * infrastructure provisions such a disk (GCP Hyperdisk today) send one; the
 * column is refreshed on every heartbeat that carries it and never nulled, so
 * an older runner binary must not wipe a newer one's disk identity.
 */
export interface RunnerDataDiskRef {
  /** Cloud provider owning the disk; `gcp` is the only writer today. */
  provider: string
  project: string
  zone: string
  /** Disk name inside the project/zone, not a full resource path. */
  diskName: string
}

export function isGcpDataDiskRef(ref: RunnerDataDiskRef | null | undefined): ref is RunnerDataDiskRef {
  return !!ref && ref.provider === 'gcp' && !!ref.project && !!ref.zone && !!ref.diskName
}

/**
 * Whether the runner's data disk supports online capacity growth. The
 * capability rides on the provider: a complete GCP ref is the one expandable
 * shape today, and other providers would join here rather than in the
 * expansion controller. This gate runs ahead of the gray-release switch —
 * when it answers false, no amount of enabling matters, the disk simply
 * cannot grow.
 */
export function isDiskCapacityExpandable(ref: RunnerDataDiskRef | null | undefined): ref is RunnerDataDiskRef {
  return isGcpDataDiskRef(ref)
}
