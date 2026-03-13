# Ignore rule JSON templates

Sample request bodies for the **create-ignore-rule** command. The structure matches the [Create Ignore Rule](https://jfrog.com/help/r/xray-rest-apis/create-ignore-rule) API (see "Sample Request" on that page).

## Request body structure

The request must have:

- **`notes`** (string, required) – Description of the ignore rule.
- **`expires_at`** (string, optional) – When the rule expires (RFC3339 / ISO 8601, e.g. `"2025-12-31T23:59:59Z"`).
- **`ignore_filters`** (object, required) – All filter criteria go inside this object. At least one of the following must be present:

| Key | Type | Description |
|-----|------|-------------|
| `cves` | array of string | CVE IDs to ignore. |
| `vulnerabilities` | array of string | Vulnerability IDs to ignore. |
| `licenses` | array of string | License names to ignore. |
| `components` | array of `{ "name", "version" }` | Components to ignore. |
| `builds` | array of `{ "name", "version" }` | Builds to ignore. |
| `artifacts` | array of `{ "name", "version", "path" }` | Artifacts to ignore (`path` must end with `/`). |
| `docker_layers` | array of string | Docker layer SHA256 hashes. |
| `release_bundles` | array of `{ "name", "version" }` | Release bundles to ignore. |
| `operational_risk` | array of string | Operational risk to ignore (e.g. `["any"]`). |
| `exposures` | object | Exposures to ignore; see below. |
| `watches` | array of string | Scope the rule to these watches. |
| `policies` | array of string | Scope the rule to these policies. |

**`exposures` object** (all keys optional, use as needed):

| Key | Type | Description |
|-----|------|-------------|
| `categories` | array of string | Exposure categories (e.g. `"secrets"`, `"services"`). |
| `file_path` | array of string | File paths to ignore. |
| `scanners` | array of string | Scanner/exposure IDs (e.g. `"EXP-12345"`). |

## Templates

| File | Description |
|------|-------------|
| `01-cve-ignore.json` | Ignore specific CVEs with expiration. |
| `02-vulnerabilities-ignore.json` | Ignore by vulnerability IDs. |
| `03-licenses-ignore.json` | Ignore specific licenses. |
| `04-component-ignore.json` | Ignore by component name/version (array). |
| `05-artifact-ignore.json` | Ignore by artifact name, version, path (array). |
| `06-build-ignore.json` | Ignore by build name/version (array). |
| `07-watches-scoped.json` | Ignore CVE only for specific watches. |
| `08-docker-layers.json` | Ignore by Docker layer SHAs. |
| `09-operational-risk.json` | Ignore by operational risk (e.g. any). |
| `10-exposures.json` | Ignore exposures by categories and file path. |
| `11-exposures-scanners.json` | Ignore exposures by scanner ID and file path. |

To **curate a combined ignore rule** (e.g. exposures + artifacts + watches), use the **`test/`** directory. It is gitignored so your drafts and test data are not committed. See `test/README.md`.

## Usage

```bash
jf xray-utils create-ignore-rule templates/create-ignore-rule/01-cve-ignore.json
jf xray-utils cir path/to/my-rule.json --server-id my-server
```
