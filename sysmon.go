package sysmon

const (
	// WellKnownBusName is the D-bus "well know bus name" that the
	// server will use for its connection.
	WellKnownBusName = "com.github.alcortesm.sysmon1"
	// InterfaceName is the D-bus interface name the sysmon server implements.
	InterfaceName = WellKnownBusName
	// Path is the single D-bus path the server will use.
	Path = "/com/github/alcortesm/sysmon1"
)

// The server interface provides a higher-level API suitable for applications
// to run, access and shutdown sysmon servers.
type Server interface {
	// Connect connects the server to the D-bus system.
	Connect() error
	// Disconnect disconnects the server from the D-bus system.
	Disconnect() error
}
