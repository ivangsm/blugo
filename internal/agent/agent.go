package agent

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const (
	bluezService = "org.bluez"
	agentIface   = "org.bluez.Agent1"
)

var agentPath = dbus.ObjectPath("/org/bluez/agent_gob")

const agentIntrospection = `
<node>
	<interface name="org.bluez.Agent1">
		<method name="Release"></method>
		<method name="RequestPinCode">
			<arg type="o" direction="in"/>
			<arg type="s" direction="out"/>
		</method>
		<method name="DisplayPinCode">
			<arg type="o" direction="in"/>
			<arg type="s" direction="in"/>
		</method>
		<method name="RequestPasskey">
			<arg type="o" direction="in"/>
			<arg type="u" direction="out"/>
		</method>
		<method name="DisplayPasskey">
			<arg type="o" direction="in"/>
			<arg type="u" direction="in"/>
			<arg type="q" direction="in"/>
		</method>
		<method name="RequestConfirmation">
			<arg type="o" direction="in"/>
			<arg type="u" direction="in"/>
		</method>
		<method name="RequestAuthorization">
			<arg type="o" direction="in"/>
		</method>
		<method name="AuthorizeService">
			<arg type="o" direction="in"/>
			<arg type="s" direction="in"/>
		</method>
		<method name="Cancel"></method>
	</interface>
	<interface name="org.freedesktop.DBus.Introspectable">
		<method name="Introspect">
			<arg name="data" type="s" direction="out"/>
		</method>
	</interface>
</node>
`

// Agent handles BlueZ pairing requests.
type Agent struct {
	program        *tea.Program
	passkeyChannel chan uint32
	confirmChannel chan bool
}

// NewAgent creates a new agent instance.
func NewAgent(program *tea.Program) *Agent {
	return &Agent{
		program:        program,
		passkeyChannel: make(chan uint32, 1),
		confirmChannel: make(chan bool, 1),
	}
}

// GetPasskeyChannel returns the passkey channel.
func (a *Agent) GetPasskeyChannel() <-chan uint32 {
	return a.passkeyChannel
}

// GetConfirmChannel returns the confirmation channel.
func (a *Agent) GetConfirmChannel() chan<- bool {
	return a.confirmChannel
}

// Release is called when the agent is unregistered.
func (a *Agent) Release() *dbus.Error {
	return nil
}

// RequestPinCode requests an old PIN (4 digits).
func (a *Agent) RequestPinCode(device dbus.ObjectPath) (string, *dbus.Error) {
	return "0000", nil
}

// DisplayPinCode displays a PIN in the interface.
func (a *Agent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	if a.program != nil {
		// TODO: Send message to UI
		fmt.Printf("PIN: %s\n", pincode)
	}
	return nil
}

// RequestPasskey requests a passkey from the user.
func (a *Agent) RequestPasskey(device dbus.ObjectPath) (uint32, *dbus.Error) {
	return 0, dbus.MakeFailedError(fmt.Errorf("no se puede solicitar passkey en TUI"))
}

// DisplayPasskey shows a 6-digit passkey.
func (a *Agent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	if a.program != nil {
		a.passkeyChannel <- passkey
	}

	// Wait for user confirmation
	confirmed := <-a.confirmChannel
	if !confirmed {
		return dbus.MakeFailedError(fmt.Errorf("pairing cancelado por el usuario"))
	}

	return nil
}

// RequestConfirmation requests confirmation of a passkey.
func (a *Agent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	if a.program != nil {
		a.passkeyChannel <- passkey
	}

	confirmed := <-a.confirmChannel
	if !confirmed {
		return dbus.MakeFailedError(fmt.Errorf("confirmación rechazada"))
	}

	return nil
}

// RequestAuthorization requests authorization for a device.
func (a *Agent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	return nil
}

// AuthorizeService requests authorization for a service.
func (a *Agent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	return nil
}

// Cancel cancels the current pairing process.
func (a *Agent) Cancel() *dbus.Error {
	if a.confirmChannel != nil {
		select {
		case a.confirmChannel <- false:
		default:
		}
	}
	return nil
}

// Register registers the agent with BlueZ.
func (a *Agent) Register(conn *dbus.Conn) error {
	// Unregister any existing agent
	a.Unregister(conn)

	// Export the agent on DBus
	err := conn.Export(a, agentPath, agentIface)
	if err != nil {
		return fmt.Errorf("no se pudo exportar agente: %w", err)
	}

	// Export introspection
	err = conn.Export(introspect.Introspectable(agentIntrospection), agentPath,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return fmt.Errorf("no se pudo exportar introspección: %w", err)
	}

	// Register the agent with BlueZ
	obj := conn.Object(bluezService, dbus.ObjectPath("/org/bluez"))

	// Try KeyboardDisplay first
	err = obj.Call("org.bluez.AgentManager1.RegisterAgent", 0, agentPath, "KeyboardDisplay").Err
	if err != nil {
		// If it fails, try NoInputNoOutput
		err = obj.Call("org.bluez.AgentManager1.RegisterAgent", 0, agentPath, "NoInputNoOutput").Err
		if err != nil {
			return fmt.Errorf("no se pudo registrar agente: %w", err)
		}
	}

	// Set as default agent
	_ = obj.Call("org.bluez.AgentManager1.RequestDefaultAgent", 0, agentPath).Err

	return nil
}

// Unregister unregisters the agent from BlueZ.
func (a *Agent) Unregister(conn *dbus.Conn) {
	obj := conn.Object(bluezService, dbus.ObjectPath("/org/bluez"))
	_ = obj.Call("org.bluez.AgentManager1.UnregisterAgent", 0, agentPath).Err
}
