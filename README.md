# GOB - Gestor de Bluetooth para Linux

Gestor de Bluetooth minimalista con interfaz de terminal (TUI) para Linux, construido con Go y BlueZ.

## Descripción

GOB es una aplicación de terminal interactiva que permite gestionar dispositivos Bluetooth de manera simple y eficiente. Utiliza el stack BlueZ a través de DBus para comunicarse con el adaptador Bluetooth del sistema, proporcionando una interfaz limpia y fácil de usar.

## Características Actuales

- **Escaneo de dispositivos**: Búsqueda automática de dispositivos Bluetooth cercanos
- **Pairing automático**: Emparejamiento de dispositivos con soporte para autenticación por passkey
- **Gestión de conexiones**: Conectar y desconectar dispositivos de forma sencilla
- **Información detallada**: Muestra nombre, dirección MAC, intensidad de señal (RSSI) y tipo de dispositivo
- **Indicador de batería**: Visualización del nivel de batería de dispositivos compatibles con colores dinámicos
- **Interfaz moderna y responsiva**:
  - Layout adaptable que cambia entre una y dos columnas según el ancho de la terminal
  - Diseño con bordes redondeados, badges y separadores elegantes
  - Paneles con resaltado visual para la sección activa
  - Componentes reutilizables y modulares
- **Iconos y badges**: Identificación visual de dispositivos y estados (pareado, conectado, confiable)
- **Control de escaneo**: Pausar y reanudar el escaneo con indicador visual en tiempo real
- **Olvidar dispositivos**: Eliminar dispositivos pareados del sistema
- **Control del adaptador Bluetooth**:
  - Ver información detallada del adaptador (nombre, dirección, estado)
  - Encender/apagar el adaptador Bluetooth (tecla `P`)
  - Activar/desactivar modo Discoverable (tecla `V`) - Hacer el adaptador visible para otros dispositivos
  - Activar/desactivar modo Pairable (tecla `B`) - Permitir emparejamiento con nuevos dispositivos
  - Panel de información con tecla `I`

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

**Acciones de dispositivos:**
- `Enter`: Conectar a un dispositivo disponible / Desconectar un dispositivo conectado
- `d` o `x`: Olvidar dispositivo (desconectar y eliminar pairing)
- `s`: Pausar/reanudar escaneo de dispositivos

**Control del adaptador:**
- `i`: Mostrar/ocultar información del adaptador
- `p`: Encender/apagar el adaptador Bluetooth
- `v`: Activar/desactivar modo Discoverable
- `b`: Activar/desactivar modo Pairable

**General:**
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

- [x] **Indicador de batería**: Mostrar el nivel de batería de dispositivos compatibles
- [x] **Mejoras en la TUI**:
  - Layout responsivo (1 o 2 columnas según ancho de terminal)
  - Diseño moderno con bordes y paneles
  - Badges y estilos de colores mejorados
  - Componentes UI reutilizables
  - [ ] Temas de color personalizables
  - [ ] Animaciones y transiciones suaves
- [x] **Control del adaptador Bluetooth**:
  - [x] Encender/apagar el adaptador por completo
  - [x] Mostrar nombre/alias del adaptador Bluetooth
  - [x] Ver información del adaptador
  - [x] Activar/desactivar modo Pairable (emparejamiento)
  - [x] Activar/desactivar modo Discoverable (detectable)
  - [ ] Reiniciar el servicio Bluetooth
  - [ ] Editar alias del adaptador
  - [ ] Configurar timeouts de visibilidad

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

Para más detalles sobre la arquitectura, patrones y mejores prácticas, consulta:
- [ARCHITECTURE.md](ARCHITECTURE.md) - Arquitectura y patrones de diseño
- [docs/UI_IMPROVEMENTS.md](docs/UI_IMPROVEMENTS.md) - Mejoras de la interfaz de usuario

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
