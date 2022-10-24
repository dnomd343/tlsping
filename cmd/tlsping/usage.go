package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"text/template"
)

type usageType uint32

const (
	usageShort usageType = iota
	usageLong
)

func printUsage(f *os.File, kind usageType) {
	const template = `
USAGE:
{{.Tab1}}{{.AppName}} [-c count] [-tcponly] [-json] [-ca=<file>]
{{.Tab1}}{{.AppNameFiller}} [-insecure] [-host=<sni>] <server address>
{{.Tab1}}{{.AppName}} -help
{{.Tab1}}{{.AppName}} -version
{{if eq .UsageVersion "short"}}
Use '{{.AppName}} -help' to get detailed information about options and examples
of usage.{{else}}

DESCRIPTION:
{{.Tab1}}{{.AppName}} is a basic tool to measure the time required to establish a
{{.Tab1}}TCP connection and perform the TLS handshake with a remote server.
{{.Tab1}}It reports summary statistics of the measurements obtained over a number
{{.Tab1}}of successful connections.

{{.Tab1}}The address of the remote server, i.e. <server address>, is of the form
{{.Tab1}}'host:port', for instance 'mail.google.com:443', '216.58.215.37:443' or
{{.Tab1}}'[2a00:1450:400a:800::2005]:443'.

OPTIONS:
{{.Tab1}}-c count
{{.Tab2}}Perform count concurrent measurements.
{{.Tab2}}Default: {{.DefaultCount}}

{{.Tab1}}-tcponly
{{.Tab2}}Establish the TCP connection with the remote server but do not perform
{{.Tab2}}the TLS handshake.

{{.Tab1}}-insecure
{{.Tab2}}Don't verify the validity of the server certificate. Only relevant when
{{.Tab2}}TLS handshake is performed (see '-tcponly' option).
{{.Tab2}}This option is intended to be used for measuring times for connecting
{{.Tab2}}to servers which use custom not widely trusted certificates, for
{{.Tab2}}instance, development servers using self-signed certificates.

{{.Tab1}}-host <sni>
{{.Tab2}}Specifies the SNI parameter of the tls connection.

{{.Tab1}}-ca <file>
{{.Tab2}}PEM-formatted file containing the CA certificates this program trusts.
{{.Tab2}}Default: use this host's CA certificate store.

{{.Tab1}}-json
{{.Tab2}}Format the result in JSON format and print to standard output. Reported
{{.Tab2}}times are understood in seconds.

{{.Tab1}}-raw
{{.Tab2}}Show all raw latency result to standard output.

{{.Tab1}}-help
{{.Tab2}}Prints this help

{{.Tab1}}-version
{{.Tab2}}Show detailed version information about this application

EXAMPLES:
{{.Tab1}}To measure the time to establish a TCP connection and perform TLS
{{.Tab1}}handshaking with host 'mail.google.com' port 443 use:

{{.Tab3}}{{.AppName}} mail.google.com:443

{{.Tab1}}Initiate a tls request to an ip address and specify sni information
{{.Tab1}}at the same time:

{{.Tab3}}{{.AppName}} -host ip.343.re 8.210.148.24:443

{{.Tab1}}To measure the time to establishing a TCP connection (i.e. without
{{.Tab1}}performing TLS handshaking) to host at IPv6 address
{{.Tab1}}'2606:4700::6811:d209' port 443 use:

{{.Tab3}}{{.AppName}} -tcponly [2a00:1450:400a:800::2005]:443
{{end}}
`
	if kind == usageLong {
		tmplFields["UsageVersion"] = "long"
	}
	tmplFields["DefaultCount"] = fmt.Sprintf("%d", defaultIterations)
	render(template, tmplFields, f)
}

func render(tpl string, fields map[string]string, out io.Writer) {
	minWidth, tabWidth, padding := 4, 4, 0
	tabwriter := tabwriter.NewWriter(out, minWidth, tabWidth, padding, byte(' '), 0)
	templ := template.Must(template.New("").Parse(tpl))
	templ.Execute(tabwriter, fields)
	tabwriter.Flush()
}
