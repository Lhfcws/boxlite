/*
 * Copyright 2026 BoxLite AI
 * SPDX-License-Identifier: AGPL-3.0
 */

/** Who set `runner.unschedulable`; only the API's own marks are lifted automatically. */
export enum RunnerUnschedulableReason {
  /** A manual mark: a person set it through the scheduling endpoint, and only a person lifts it. */
  OPERATOR = 'operator',
  DISK_PRESSURE = 'disk_pressure',
  DISK_EXPAND_EXHAUSTED = 'disk_expand_exhausted',
}
