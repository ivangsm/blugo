# BLUGO - Gestor de Bluetooth para Linux

> Gestor de Bluetooth minimalista con interfaz de terminal (TUI) moderna para Linux, construido con Go y BlueZ.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![BlueZ](https://img.shields.io/badge/BlueZ-5.0+-blue.svg)](http://www.bluez.org/)

Español | [English](README.md)

---

## Características

### Gestión de Dispositivos
- **Escaneo automático** de dispositivos Bluetooth cercanos
- **Pairing automático** con soporte para autenticación por passkey
- **Conectar/desconectar** dispositivos fácilmente
- **Olvidar dispositivos** para eliminar el pairing del sistema
- **Información detallada**: nombre, dirección MAC, intensidad de señal (RSSI) y tipo de dispositivo
- **Indicador de batería** con colores dinámicos para dispositivos compatibles

### Control del Adaptador
- **Control de energía**: Encender/apagar el adaptador Bluetooth (tecla `P`)
- **Modo Discoverable**: Hacer el adaptador visible para otros dispositivos (tecla `V`)
- **Modo Pairable**: Permitir emparejamiento con nuevos dispositivos (tecla `B`)
- **Información del adaptador**: Ver estado detallado y configuración del adaptador

### Interfaz Moderna
- **Diseño responsivo**: Se adapta a cualquier tamaño de terminal
- **Layout limpio**: Diseño de una columna con espaciado apropiado
- **Feedback visual**: Iconos, badges y estados con códigos de color
- **Actualizaciones en tiempo real**: Refresco automático cada 2 segundos
- **Control de escaneo**: Pausar/reanudar el escaneo con indicador visual
- **Modo pantalla alternativa**: Terminal limpia al salir y redimensionar
- **Internacionalización**: Soporte completo de i18n (Inglés/Español), fácilmente extensible para más idiomas

---

## Capturas de Pantalla

```
╭────────────────────────────────────────────────────────────────╮
│ 🔵 BLUGO - Gestor Bluetooth               🔍 Escaneando         │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│ ⚡ Adaptador Bluetooth                                          │
│ ────────────────────────────────────────────────────────────   │
│ Name            Alias           Power    Pairable  Discoverable│
│ bluez-5.0       Mi Laptop       ON       OFF       OFF         │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│                                                                │
│ 📡 DISPOSITIVOS DISPONIBLES (2) ◀                              │
│ ────────────────────────────────────────────────────────────   │
│                                                                │
│ ▶ 🎧 Sony WH-1000XM4 (AA:BB:CC:DD:EE:FF) | -45 dBm | 🔋 85%   │
│   ⌨️  Keychron K3 (11:22:33:44:55:66) | -38 dBm PAREADO        │
│                                                                │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│                                                                │
│ 🔗 DISPOSITIVOS CONECTADOS (1)                                 │
│ ────────────────────────────────────────────────────────────   │
│                                                                │
│   🖱️  Logitech MX Master 3 (FF:EE:DD:CC:BB:AA) | 🪫 12%       │
│                                                                │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│ ↑/↓: navegar | Tab: cambiar | Enter: conectar | D/X: olvidar  │
│ S: escaneo | P: encendido | V: descubrible | B: pairable      │
│ Q: salir                                                       │
╰────────────────────────────────────────────────────────────────╯
```

---

## Requisitos

- **Linux** con BlueZ instalado
- **Go 1.23** o superior
- **Adaptador Bluetooth** compatible con BlueZ
- Acceso a **DBus** del sistema

### Instalación de BlueZ

**Arch Linux / Manjaro:**
```bash
sudo pacman -S bluez bluez-utils
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```

**Ubuntu / Debian:**
```bash
sudo apt install bluez bluetooth
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```

**Fedora:**
```bash
sudo dnf install bluez bluez-tools
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```

---

## Instalación

### Opción 1: Usando Make (Recomendado)

```bash
# Clonar el repositorio
git clone https://github.com/ivangsm/blugo.git
cd gob

# Compilar
make build

# (Opcional) Instalar en el sistema
make install
```

### Opción 2: Manualmente con Go

```bash
# Clonar el repositorio
git clone https://github.com/ivangsm/blugo.git
cd gob

# Descargar dependencias
go mod download

# Compilar desde cmd/blugo
go build -o gob ./cmd/blugo

# (Opcional) Instalar en el sistema
sudo mv gob /usr/local/bin/
```

---

## Uso

Simplemente ejecuta el binario:
```bash
./blugo
```

O si lo instalaste en el sistema:
```bash
gob
```

### Controles de Teclado

**Navegación:**
- `↑/↓` o `k/j`: Navegar entre dispositivos
- `Tab`: Cambiar entre secciones (disponibles/conectados)
- `r`: Refrescar lista de dispositivos manualmente

**Acciones de Dispositivos:**
- `Enter`: Conectar a un dispositivo disponible / Desconectar un dispositivo conectado
- `d` o `x`: Olvidar dispositivo (desconectar y eliminar pairing)
- `s`: Pausar/reanudar escaneo de dispositivos

**Control del Adaptador:**
- `p`: Encender/apagar el adaptador Bluetooth
- `v`: Activar/desactivar modo Discoverable
- `b`: Activar/desactivar modo Pairable
- `l`: Cambiar idioma (Inglés/Español)

**General:**
- `q` o `Ctrl+C`: Salir de la aplicación

**Durante el Pairing:**
- `Enter` o `y`: Confirmar código de pairing
- `n` o `Esc`: Cancelar pairing

---

## Arquitectura

BLUGO sigue principios SOLID y patrones de arquitectura limpia:

### Estructura del Proyecto

```
blugo/
├── cmd/blugo/              # Entry point de la aplicación
├── internal/
│   ├── models/           # Modelos de datos
│   ├── agent/            # Agente de pairing Bluetooth
│   ├── bluetooth/        # Gestión de Bluetooth/DBus
│   └── ui/               # Interfaz de Usuario de Terminal
│       ├── styles.go     # Estilos de Lipgloss
│       ├── components.go # Componentes UI reutilizables
│       ├── model.go      # Estado de la aplicación
│       ├── update.go     # Lógica de actualización (TEA)
│       ├── view.go       # Lógica de renderizado (TEA)
│       ├── messages.go   # Tipos de mensajes
│       └── commands.go   # Comandos asíncronos
├── Makefile
├── Dockerfile
├── go.mod
└── README.md
```

### Principios de Diseño

- **Separación de Responsabilidades**: Cada paquete tiene una responsabilidad única y clara
- **Responsabilidad Única**: Cada archivo maneja un aspecto específico
- **Inversión de Dependencias**: Los módulos de alto nivel no dependen de detalles de implementación de bajo nivel
- **The Elm Architecture**: UI reactiva con patrón Model-Update-View

### Tecnologías

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: Framework TUI basado en The Elm Architecture
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)**: Librería de estilos para terminales
- **[godbus](https://github.com/godbus/dbus)**: Cliente DBus para Go

---

## Desarrollo

### Comandos Make Disponibles

```bash
make build        # Compilar la aplicación
make run          # Compilar y ejecutar
make install      # Instalar en /usr/local/bin
make clean        # Limpiar archivos compilados
make test         # Ejecutar tests
make fmt          # Formatear código
make help         # Ver todos los comandos
```

### Agregar Nuevas Características

**Nueva funcionalidad de Bluetooth:**
1. Agregar método en `bluetooth/adapter.go` o `bluetooth/device.go`
2. Crear comando en `ui/commands.go`
3. Agregar handler en `ui/update.go`
4. Actualizar vista en `ui/view.go` si es necesario

**Nueva sección de UI:**
1. Agregar estado al modelo en `ui/model.go`
2. Crear mensaje en `ui/messages.go`
3. Implementar handler en `ui/update.go`
4. Crear función de renderizado en `ui/view.go`

---

## Docker

La aplicación puede ejecutarse en un contenedor Docker, aunque requiere acceso privilegiado al DBus y hardware Bluetooth:

**Construir la imagen:**
```bash
docker build -t gob .
```

**Ejecutar:**
```bash
docker run --rm -it --privileged --net=host \
  -v /var/run/dbus:/var/run/dbus \
  gob
```

**Nota:** El uso de Docker es limitado debido a los requisitos de acceso al hardware Bluetooth. Se recomienda la instalación nativa.

---

## Roadmap

### Versión Actual
- ✅ Indicador de batería para dispositivos compatibles
- ✅ Control completo del adaptador (energía, discoverable, pairable)
- ✅ TUI moderna y responsiva con layouts apropiados
- ✅ Badges e indicadores de estado con códigos de color
- ✅ Escaneo en tiempo real con pausar/reanudar
- ✅ Manejo limpio de terminal (modo pantalla alternativa)

### Características Planeadas
- [ ] Soporte para múltiples adaptadores Bluetooth
- [ ] Perfiles Bluetooth específicos (A2DP, HFP, etc.)
- [ ] Historial de conexiones
- [ ] Configuración persistente
- [ ] Logging y debugging mejorado
- [ ] Tests unitarios e integración
- [ ] Soporte de archivos de configuración (TOML/YAML)
- [ ] Temas de color personalizables

---

## Contribuir

¡Las contribuciones son bienvenidas! Por favor:

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/caracteristica-increible`)
3. Commit tus cambios (`git commit -am 'Agregar característica increíble'`)
4. Push a la rama (`git push origin feature/caracteristica-increible`)
5. Crea un Pull Request

---

## Licencia

MIT

---

## Autor

Ivan - [@ivangsm](https://github.com/ivangsm)

---

## Agradecimientos

- Proyecto BlueZ por el stack Bluetooth de Linux
- Charm.sh por las excelentes herramientas de TUI
- La comunidad de Go por las increíbles librerías y soporte
