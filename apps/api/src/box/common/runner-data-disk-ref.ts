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
