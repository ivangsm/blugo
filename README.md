# GOB - Gestor de Bluetooth para Linux

Gestor de Bluetooth minimalista con interfaz de terminal (TUI) para Linux, construido con Go y BlueZ.

## Descripción

GOB es una aplicación de terminal interactiva que permite gestionar dispositivos Bluetooth de manera simple y eficiente. Utiliza el stack BlueZ a través de DBus para comunicarse con el adaptador Bluetooth del sistema, proporcionando una interfaz limpia y fácil de usar.

## Características Actuales

- **Escaneo de dispositivos**: Búsqueda automática de dispositivos Bluetooth cercanos
- **Pairing automático**: Emparejamiento de dispositivos con soporte para autenticación por passkey
- **Gestión de conexiones**: Conectar y desconectar dispositivos de forma sencilla
- **Información detallada**: Muestra nombre, dirección MAC, intensidad de señal (RSSI) y tipo de dispositivo
- **Interfaz dual**: Visualización separada de dispositivos disponibles y conectados
- **Iconos por tipo**: Identificación visual de dispositivos (auriculares, teléfonos, teclados, etc.)
- **Control de escaneo**: Pausar y reanudar el escaneo para ahorrar batería
- **Olvidar dispositivos**: Eliminar dispositivos pareados del sistema

## Requisitos

- Linux con BlueZ instalado
- Go 1.23 o superior
- Adaptador Bluetooth compatible
- Acceso a DBus del sistema

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

## Instalación

### Opción 1: Usando Make (Recomendado)

```bash
# Clonar el repositorio
git clone https://github.com/ivangsm/gob.git
cd gob

# Compilar
make build

# (Opcional) Instalar en el sistema
make install
```

### Opción 2: Manualmente con Go

```bash
# Clonar el repositorio
git clone https://github.com/ivangsm/gob.git
cd gob

# Descargar dependencias
go mod download

# Compilar desde cmd/gob
go build -o gob ./cmd/gob

# (Opcional) Instalar en el sistema
sudo mv gob /usr/local/bin/
```

## Uso

Simplemente ejecuta el binario:
```bash
./gob
```

O si lo instalaste en el sistema:
```bash
gob
```

### Controles

**Navegación:**
- `↑/↓` o `k/j`: Navegar entre dispositivos
- `Tab`: Cambiar entre secciones (disponibles/conectados)
- `r`: Refrescar lista de dispositivos manualmente

**Acciones:**
- `Enter`: Conectar a un dispositivo disponible / Desconectar un dispositivo conectado
- `d` o `x`: Olvidar dispositivo (desconectar y eliminar pairing)
- `s`: Pausar/reanudar escaneo de dispositivos
- `q` o `Ctrl+C`: Salir de la aplicación

**Durante el pairing:**
- `Enter` o `y`: Confirmar código de pairing
- `n` o `Esc`: Cancelar pairing

## Docker

La aplicación puede ejecutarse en un contenedor Docker, aunque requiere acceso privilegiado al DBus del sistema y al hardware Bluetooth:

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

**Nota:** El uso de Docker para esta aplicación es limitado debido a los requisitos de acceso al hardware Bluetooth. Se recomienda la instalación nativa.

## Roadmap

### Próximas características

- [ ] **Indicador de batería**: Mostrar el nivel de batería de dispositivos compatibles
- [ ] **Mejoras en la TUI**:
  - Diseño más moderno y personalizable
  - Temas de color
  - Mejor visualización de información
  - Animaciones y transiciones suaves
- [ ] **Control del adaptador Bluetooth**:
  - Encender/apagar el adaptador por completo
  - Reiniciar el servicio Bluetooth
- [ ] **Información del adaptador**:
  - Mostrar nombre/alias del adaptador Bluetooth
  - Ver y editar información del adaptador
- [ ] **Modos del adaptador**:
  - Activar/desactivar modo Pairable (emparejamiento)
  - Activar/desactivar modo Discoverable (detectable)
  - Configurar timeouts de visibilidad

### Futuras mejoras

- [ ] Soporte para múltiples adaptadores Bluetooth
- [ ] Perfiles Bluetooth específicos (A2DP, HFP, etc.)
- [ ] Historial de conexiones
- [ ] Configuración persistente
- [ ] Logs y debugging mejorado
- [ ] Tests unitarios e integración
- [ ] Configuración por archivo TOML/YAML

## Arquitectura

GOB sigue principios SOLID y una arquitectura modular profesional:

### Estructura del Proyecto

```
gob/
├── cmd/gob/              # Entry point de la aplicación
├── internal/
│   ├── models/           # Modelos de datos
│   ├── agent/            # Agente de pairing Bluetooth
│   ├── bluetooth/        # Gestión de Bluetooth/DBus
│   └── ui/               # Interfaz de usuario TUI
├── ARCHITECTURE.md       # Documentación detallada de arquitectura
├── Makefile             # Automatización de tareas
└── ...
```

### Tecnologías

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: Framework TUI basado en The Elm Architecture
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)**: Librería de estilos para terminales
- **[godbus](https://github.com/godbus/dbus)**: Cliente DBus para Go

### Principios de Diseño

- **Separation of Concerns**: Cada paquete tiene una responsabilidad única
- **Single Responsibility**: Cada archivo maneja un aspecto específico
- **Dependency Inversion**: Módulos de alto nivel no dependen de detalles de implementación
- **The Elm Architecture**: UI reactiva con Model-Update-View

Para más detalles sobre la arquitectura, patrones y mejores prácticas, consulta [ARCHITECTURE.md](ARCHITECTURE.md).

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

## Contribuir

Las contribuciones son bienvenidas. Por favor:
1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-caracteristica`)
3. Commit tus cambios (`git commit -am 'Agregar nueva característica'`)
4. Push a la rama (`git push origin feature/nueva-caracteristica`)
5. Crea un Pull Request

## Licencia

[Especificar licencia]

## Autor

Ivan - [@ivangsm](https://github.com/ivangsm)

## Agradecimientos

- Proyecto BlueZ por el stack Bluetooth de Linux
- Charm.sh por las excelentes herramientas de TUI
