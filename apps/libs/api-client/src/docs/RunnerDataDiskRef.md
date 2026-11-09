# RunnerDataDiskRef

Cloud identity of the runner data disk. Disk type is intentionally omitted: the cloud API addresses a disk by name alone, the actual type (persistent disk vs hyperdisk) is derivable via compute.disks.get when monitoring needs it, and a future non-disk backend would need a different schema rather than an extra string field.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**provider** | **string** | Cloud provider of the data disk | [default to undefined]
**project** | **string** | Cloud project the data disk lives in | [default to undefined]
**zone** | **string** | Cloud zone the data disk lives in | [default to undefined]
**diskName** | **string** | Name of the data disk, as the cloud API addresses it | [default to undefined]

## Example

```typescript
import { RunnerDataDiskRef } from './api';

const instance: RunnerDataDiskRef = {
    provider,
    project,
    zone,
    diskName,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
