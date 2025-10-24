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

// Agent maneja las solicitudes de pairing de BlueZ.
type Agent struct {
	program        *tea.Program
	passkeyChannel chan uint32
	confirmChannel chan bool
}

// NewAgent crea una nueva instancia del agente.
func NewAgent(program *tea.Program) *Agent {
	return &Agent{
		program:        program,
		passkeyChannel: make(chan uint32, 1),
		confirmChannel: make(chan bool, 1),
	}
}

// GetPasskeyChannel devuelve el canal de passkeys.
func (a *Agent) GetPasskeyChannel() <-chan uint32 {
	return a.passkeyChannel
}

// GetConfirmChannel devuelve el canal de confirmación.
func (a *Agent) GetConfirmChannel() chan<- bool {
	return a.confirmChannel
}

// Release se llama cuando el agente se desregistra.
func (a *Agent) Release() *dbus.Error {
	return nil
}

// RequestPinCode solicita un PIN antiguo (4 dígitos).
func (a *Agent) RequestPinCode(device dbus.ObjectPath) (string, *dbus.Error) {
	return "0000", nil
}

// DisplayPinCode muestra un PIN en la interfaz.
func (a *Agent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	if a.program != nil {
		// TODO: Enviar mensaje a la UI
		fmt.Printf("PIN: %s\n", pincode)
	}
	return nil
}

// RequestPasskey solicita un passkey del usuario.
func (a *Agent) RequestPasskey(device dbus.ObjectPath) (uint32, *dbus.Error) {
	return 0, dbus.MakeFailedError(fmt.Errorf("no se puede solicitar passkey en TUI"))
}

// DisplayPasskey muestra un passkey de 6 dígitos.
func (a *Agent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	if a.program != nil {
		a.passkeyChannel <- passkey
	}

	// Esperar confirmación del usuario
	confirmed := <-a.confirmChannel
	if !confirmed {
		return dbus.MakeFailedError(fmt.Errorf("pairing cancelado por el usuario"))
	}

	return nil
}

// RequestConfirmation solicita confirmación de un passkey.
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

// RequestAuthorization solicita autorización para un dispositivo.
func (a *Agent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	return nil
}

// AuthorizeService solicita autorización para un servicio.
func (a *Agent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	return nil
}

// Cancel cancela el proceso de pairing actual.
func (a *Agent) Cancel() *dbus.Error {
	if a.confirmChannel != nil {
		select {
		case a.confirmChannel <- false:
		default:
		}
	}
	return nil
}

// Register registra el agente en BlueZ.
func (a *Agent) Register(conn *dbus.Conn) error {
	// Desregistrar cualquier agente existente
	a.Unregister(conn)

	// Exportar el agente en DBus
	err := conn.Export(a, agentPath, agentIface)
	if err != nil {
		return fmt.Errorf("no se pudo exportar agente: %w", err)
	}

	// Exportar introspección
	err = conn.Export(introspect.Introspectable(agentIntrospection), agentPath,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return fmt.Errorf("no se pudo exportar introspección: %w", err)
	}

	// Registrar el agente con BlueZ
	obj := conn.Object(bluezService, dbus.ObjectPath("/org/bluez"))

	// Intentar con KeyboardDisplay primero
	err = obj.Call("org.bluez.AgentManager1.RegisterAgent", 0, agentPath, "KeyboardDisplay").Err
	if err != nil {
		// Si falla, intentar con NoInputNoOutput
		err = obj.Call("org.bluez.AgentManager1.RegisterAgent", 0, agentPath, "NoInputNoOutput").Err
		if err != nil {
			return fmt.Errorf("no se pudo registrar agente: %w", err)
		}
	}

	// Establecer como agente por defecto
	_ = obj.Call("org.bluez.AgentManager1.RequestDefaultAgent", 0, agentPath).Err

	return nil
}

// Unregister desregistra el agente de BlueZ.
func (a *Agent) Unregister(conn *dbus.Conn) {
	obj := conn.Object(bluezService, dbus.ObjectPath("/org/bluez"))
	_ = obj.Call("org.bluez.AgentManager1.UnregisterAgent", 0, agentPath).Err
}
