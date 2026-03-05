<p align="center">
  <img src="logo.png" alt="Kamaleon Logo" width="200"/>
</p>

# 🦎 Kamaleon (`kml`)

Kamaleon es un envoltorio (*wrapper*) universal, ultrarrápido y sin dependencias escrito en Go para los gestores de paquetes más populares de Linux.

Olvídate de cambiar el chip y la sintaxis cuando saltas entre distribuciones o servidores. Kamaleon detecta automáticamente tu sistema y traduce comandos genéricos al gestor nativo subyacente de forma transparente.

**Soporta detección automática de:** `yay`, `paru`, `pacman`, `apt`, `dnf`, `zypper`, `apk`, `flatpak`  y `xbps`.

## 🚀 Características

- **Un único binario:** Escrito en Go, sin scripts pesados ni intérpretes adicionales.
- **Inteligente:** Prioriza *helpers* de AUR (`yay`/`paru`) si estás en el ecosistema Arch/Manjaro.
- **Gestión de privilegios:** Añade `sudo` automáticamente solo cuando es necesario (y lo evita al compilar en AUR).
- **Protección Anti-Root:** Bloquea la ejecución directa con `sudo` para evitar romper los permisos del sistema.
- **Autocompletado nativo:** Soporte total para Zsh.

## 📦 Instalación

Clona el repositorio y compila el binario utilizando Go:

```bash
git clone https://github.com/anabasasoft/kamaleon.git
cd kamaleon
go mod init kml
go build -ldflags="-s -w" -o kml kml.go
sudo cp kml /usr/local/bin/
```

### Autocompletado

Para que `kml` se sienta como una herramienta nativa, habilita el autocompletado según tu shell:

#### Para Zsh 
Copia el archivo de funciones a tu directorio de sistema:
```bash
sudo cp _kml /usr/share/zsh/site-functions/
rm -f ~/.zcompdump; compinit
```
#### Para Bash

Copia el script de completado al directorio estándar de Bash:
```bash
sudo cp kml-completion.bash /usr/share/bash-completion/completions/kml
```

Si no tienes permisos de root, puedes añadir `source ~/ruta/a/kml-completion.bash` en tu archivo `~/.bashrc`.

## 🛠️ Uso y Comandos

La sintaxis siempre es: `kml <acción> [paquetes...]`

| Acción | Alias | Descripción | Ejemplo en Manjaro/Arch |
| --- | --- | --- | --- |
| **Instalar** | `in`, `install` | Instala paquetes desde repos oficiales o AUR. | `yay -S htop` |
| **Quitar** | `rm`, `remove` | Desinstala paquetes y limpia sus dependencias. | `yay -Rs htop` |
| **Actualizar** | `up`, `update` | Actualiza todo el sistema (modo conservador). | `yay -Syu` |
| **Full Upgrade** | `dup`, `dist-upgrade` | Actualización profunda forzando bases de datos. | `yay -Syyu` |
| **Refrescar** | `re`, `refresh` | Descarga las listas de paquetes sin instalar nada. | `yay -Sy` |
| **Buscar** | `se`, `search` | Busca filtrando **solo** por el nombre del paquete. | *Búsqueda estricta + grep* |
| **Buscar (Amplio)** | `sd`, `bdesc` | Busca en el nombre y en la descripción. | `yay -Ss editor` |
| **Información** | `info`, `show` | Muestra los detalles técnicos de un paquete. | `yay -Si zed` |
| **Listar** | `li`, `list` | Muestra los paquetes instalados o busca entre ellos. | `yay -Qe` |
| **Huérfanos** | `ar`, `autoremove` | Elimina las dependencias que ya no se utilizan. | `yay -Yc` |
| **Limpiar** | `cl`, `clean` | Vacía la caché de paquetes descargados. | `yay -Sc` |

## 🤝 Contribuir

¡Las *Pull Requests* son bienvenidas! Si usas un gestor de paquetes diferente siéntete libre de añadirlo al bloque de traducciones.

## 📧 Contacto

Para preguntas, sugerencias o reportar problemas: **anabasasoft@gmail.com**
