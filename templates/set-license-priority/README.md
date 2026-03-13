# set-license-priority

The **set-license-priority** command does not use file templates. It takes two arguments:

- **license-name** – The name of the license.
- **license-priority** – The priority value (e.g. `1`, `2`, `3`).

Optional flag: **--server-id** for a non-default JFrog CLI server.

## Usage

```bash
jf xray-utils set-license-priority "MIT" 1
jf xray-utils slp "Apache-2.0" 2 --server-id my-server
```
