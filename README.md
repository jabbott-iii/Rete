## Salus

Salus is an environment health checker CLI. It verifies disk space, memory, CPU
load, Docker status, Kubernetes status, service uptime, and common
misconfigurations, then reports a concise PASS/WARN/FAIL summary.

## Features:

- **System Resources**
  - Disk space usage on a configurable mount point
  - Memory and swap usage
  - CPU load average relative to available CPUs

- **Container Runtime**
  - Docker daemon reachability

- **Orchestration**
  - Kubernetes cluster reachability via `kubectl`

- **Services**
  - systemd service uptime (or host uptime when no service is specified)

- **Configuration**
  - Common environment misconfigurations (missing env vars, unsafe file permissions)

- **Job Tracking**
  - Every `check run` is recorded as a job in a local SQLite database
  - Browse past runs and their individual results

## Core CLI capabilities

Salus is organized into focused command groups:

- `salus check` — list and run health checks
- `salus jobs` — view past health check runs

### check

- `salus check list` — list available health checks
- `salus check run` — run health checks and report the results

Examples:
```
salus check list
salus check run
salus check run --only disk-space,memory,cpu-load
salus check run --service nginx
salus check run --json
salus check run --fail-only
salus check run --disk-path /data
```

Flags for `check run`:
- `--only` — comma-separated list of checks to run (default: all)
- `--service` — systemd service name to check uptime for (defaults to host uptime)
- `--disk-path` — mount path to check for free disk space (default `/`)
- `--json` — output results as JSON
- `--fail-only` — only show WARN and FAIL results in text output
- `--quiet` — suppress output (still sets the exit code)
- `--no-save` — do not persist this run to the database

`check run` exits with code `0` when every check passes, `1` if any check
reports WARN, and `2` if any check reports FAIL — making it suitable for use
in scripts and CI pipelines.

### jobs

- `salus jobs list` — list recent health check runs
- `salus jobs show [job-id]` — show details for a specific run

Examples:
```
salus jobs list
salus jobs list --limit 50
salus jobs show 7
```

## Install:

Download the appropriate binary for your platform below and make it executable:

Linux:
```
chmod +x salus_linux_amd64
```
```
sudo mv salus_linux_amd64 /usr/local/bin/salus
```
 or
```
chmod +x salus_linux_arm64
```
```
sudo mv salus_linux_arm64 /usr/local/bin/salus
```
macOS:
```
chmod +x salus_darwin_arm64
```
```
sudo mv salus_darwin_arm64 /usr/local/bin/salus
```
  or
```
chmod +x salus_darwin_amd64
```
```
sudo mv salus_darwin_amd64 /usr/local/bin/salus
```
Windows:
```
Download salus_windows_amd64.exe and add it to your PATH as salus.
```

## Docker

Build the image:
```bash
docker build -t salus .
```

Run a health check:
```bash
docker run -it --rm \
  -v salus-data:/app/data \
  -e SALUS_DB_PATH=/app/data/salus.db \
  salus check run
```

Note:
 - The container uses `SALUS_DB_PATH=/app/data/salus.db` by default.
 - Database state is persisted in `/app/data`.
 - Checking Docker or Kubernetes status from inside the container requires
   mounting the Docker socket or a kubeconfig, respectively.
