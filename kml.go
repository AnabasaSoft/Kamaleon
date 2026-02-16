package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const Version = "2.0.0"

// --- DECLARACIÓN DE COLORES ---
const (
	Rojo    = "\033[0;31m"
	Verde   = "\033[1;32m"
	Amarillo= "\033[1;33m"
	Azul    = "\033[1;34m"
	Cian    = "\033[1;36m"
	Magenta = "\033[1;35m"
	Blanco  = "\033[1;37m"
	Negrita = "\033[1m"
	Reset   = "\033[0m"
)

func main() {
	// 0. Escudo Anti-Root
	usuarioActual, err := user.Current()
	if err == nil && usuarioActual.Uid == "0" {
		fmt.Printf("%s%sError:%s No ejecutes kml con 'sudo'.\n", Rojo, Negrita, Reset)
		fmt.Println("El programa ya solicita permisos automáticamente cuando es necesario.")
		os.Exit(1)
	}

	// 1. Autodetectar el gestor nativo
	gestor := detectarGestor()
	if gestor == "" {
		fmt.Printf("%sError:%s No se ha detectado ningún gestor de paquetes soportado.\n", Rojo, Reset)
		os.Exit(1)
	}

	// 2. Comprobar argumentos
	if len(os.Args) < 2 {
		mostrarAyuda(gestor)
		os.Exit(0)
	}

	accion := os.Args[1]
	paquetes := ""
	if len(os.Args) > 2 {
		paquetes = strings.Join(os.Args[2:], " ")
	}

	if accion == "-h" || accion == "--help" {
		mostrarAyuda(gestor)
		os.Exit(0)
	}

	// 3. Mapear la acción al comando real
	comandoStr := mapearComando(gestor, accion, paquetes)
	if comandoStr == "" {
		fmt.Printf("%sError:%s Acción '%s%s%s' no reconocida.\n", Rojo, Reset, Negrita, accion, Reset)
		fmt.Println("Ejecuta 'kml --help' para ver las opciones disponibles.")
		os.Exit(1)
	}

	// 4. Transparencia y ejecución
	if !strings.HasPrefix(comandoStr, "echo") {
		fmt.Printf("%s%s==>%s Ejecutando: %s%s%s\n", Azul, Negrita, Reset, Negrita, comandoStr, Reset)
	}

	ejecutarComando(comandoStr)
}

func detectarGestor() string {
	candidatos := []string{"yay", "paru", "pacman", "apt", "dnf", "zypper"}
	for _, gestor := range candidatos {
		if _, err := exec.LookPath(gestor); err == nil {
			return gestor
		}
	}
	return ""
}

func mostrarAyuda(gestor string) {
	fmt.Printf("%s%skml%s - Kamaleon - Envoltorio universal para gestores de paquetes (v%s)\n", Cian, Negrita, Reset, Version)
	fmt.Printf("Gestor actual detectado: %s%s%s%s\n\n", Verde, Negrita, gestor, Reset)
	fmt.Printf("%sUSO:%s\n  kml <acción> [paquetes...]\n\n", Amarillo, Reset)
	fmt.Printf("%sACCIONES PRINCIPALES:%s\n", Amarillo, Reset)
	fmt.Printf("  %sin, instalar, install%s    Instala paquetes (repos oficiales y AUR si aplica).\n", Verde, Reset)
	fmt.Printf("  %srm, quitar, remove%s       Desinstala paquetes (limpiando dependencias).\n", Rojo, Reset)
	fmt.Printf("  %sup, actualizar, update%s   Actualiza el sistema de forma conservadora.\n", Cian, Reset)
	fmt.Printf("  %sdup, dist-upgrade%s        Actualización profunda (-Syyu, full-upgrade...).\n", Blanco, Reset)
	fmt.Printf("  %sre, refrescar, refresh%s   Actualiza únicamente la lista de repositorios.\n", Magenta, Reset)
	fmt.Printf("  %sse, buscar, search%s       Busca estrictamente en el NOMBRE del paquete.\n", Azul, Reset)
	fmt.Printf("  %ssd, bdesc, sdesc%s         Busca de forma amplia (nombre y DESCRIPCIÓN).\n", Azul, Reset)
	fmt.Printf("  %sinfo, ver, show%s          Muestra los detalles completos de un paquete.\n", Amarillo, Reset)
	fmt.Printf("  %sli, lista, list%s          Muestra paquetes instalados o busca entre ellos.\n", Negrita, Reset)
	fmt.Printf("  %sar, autoremove, huerfanos%s Elimina dependencias huérfanas que ya no se usan.\n", Rojo, Reset)
	fmt.Printf("  %scl, limpiar, clean%s       Vacia la caché de paquetes descargados.\n\n", Rojo, Reset)
	fmt.Printf("%sOPCIONES:%s\n  -h, --help               Muestra este mensaje de ayuda y sale.\n\n", Amarillo, Reset)
}

func mapearComando(gestor, accion, paq string) string {
	esAUR := (gestor == "yay" || gestor == "paru")
	esPacman := (gestor == "pacman" || esAUR)

	switch accion {
	case "in", "instalar", "install":
		if esAUR { return gestor + " -S " + paq }
		if gestor == "pacman" { return "sudo pacman -S " + paq }
		if gestor == "apt" { return "sudo apt install " + paq }
		if gestor == "dnf" { return "sudo dnf install " + paq }
		if gestor == "zypper" { return "sudo zypper in " + paq }

	case "rm", "quitar", "remove":
		if esAUR { return gestor + " -Rs " + paq }
		if gestor == "pacman" { return "sudo pacman -Rs " + paq }
		if gestor == "apt" { return "sudo apt remove " + paq }
		if gestor == "dnf" { return "sudo dnf remove " + paq }
		if gestor == "zypper" { return "sudo zypper rm " + paq }

	case "up", "actualizar", "update":
		if esAUR { return gestor + " -Syu" }
		if gestor == "pacman" { return "sudo pacman -Syu" }
		if gestor == "apt" { return "sudo apt update && sudo apt upgrade" }
		if gestor == "dnf" { return "sudo dnf upgrade" }
		if gestor == "zypper" { return "sudo zypper up" }

	case "dup", "dist-upgrade":
		if esAUR { return gestor + " -Syyu" }
		if gestor == "pacman" { return "sudo pacman -Syyu" }
		if gestor == "apt" { return "sudo apt update && sudo apt full-upgrade" }
		if gestor == "dnf" { return "sudo dnf distro-sync" }
		if gestor == "zypper" { return "sudo zypper dup" }

	case "re", "refrescar", "refresh":
		if esAUR { return gestor + " -Sy" }
		if gestor == "pacman" { return "sudo pacman -Sy" }
		if gestor == "apt" { return "sudo apt update" }
		if gestor == "dnf" { return "sudo dnf makecache" }
		if gestor == "zypper" { return "sudo zypper ref" }

	case "se", "buscar", "search":
		if esPacman { return gestor + " -Ss " + paq + ` | grep --color=auto -i -E -A 1 "^[^/]+/[^ ]*` + paq + `"` }
		if gestor == "apt" { return "apt search --names-only " + paq }
		if gestor == "dnf" { return "dnf list \"*" + paq + "*\" 2>/dev/null" }
		if gestor == "zypper" { return "zypper se -n " + paq }

	case "sd", "bdesc", "sdesc":
		if esPacman { return gestor + " -Ss " + paq }
		if gestor == "apt" { return "apt search " + paq }
		if gestor == "dnf" { return "dnf search " + paq }
		if gestor == "zypper" { return "zypper se -d " + paq }

	case "info", "ver", "show":
		if esPacman { return gestor + " -Si " + paq }
		if gestor == "apt" { return "apt show " + paq }
		if gestor == "dnf" { return "dnf info " + paq }
		if gestor == "zypper" { return "zypper info " + paq }

	case "li", "lista", "list":
		if paq == "" {
			if esPacman { return gestor + " -Qe" }
			if gestor == "apt" { return "apt list --installed 2>/dev/null" }
			if gestor == "dnf" { return "dnf list installed 2>/dev/null" }
			if gestor == "zypper" { return "zypper se -i" }
		} else {
			if esPacman { return gestor + " -Q | grep --color=auto -i \"" + paq + "\"" }
			if gestor == "apt" { return "apt list --installed 2>/dev/null | grep --color=auto -i \"" + paq + "\"" }
			if gestor == "dnf" { return "dnf list installed 2>/dev/null | grep --color=auto -i \"" + paq + "\"" }
			if gestor == "zypper" { return "zypper se -i \"" + paq + "\"" }
		}

	case "ar", "autoremove", "huerfanos":
		if gestor == "yay" { return "yay -Yc" }
		if gestor == "paru" { return "paru -c" }
		if gestor == "pacman" { return "pacman -Qdtq | sudo xargs -r pacman -Rns" }
		if gestor == "apt" { return "sudo apt autoremove" }
		if gestor == "dnf" { return "sudo dnf autoremove" }
		if gestor == "zypper" { return "echo -e \"" + Amarillo + "Nota:" + Reset + " Zypper gestiona los huérfanos con 'rm'. Para verlos: zypper packages --unneeded\"" }

	case "cl", "limpiar", "clean":
		if esAUR { return gestor + " -Sc" }
		if gestor == "pacman" { return "sudo pacman -Sc" }
		if gestor == "apt" { return "sudo apt clean" }
		if gestor == "dnf" { return "sudo dnf clean all" }
		if gestor == "zypper" { return "sudo zypper clean" }
	}
	return ""
}

// ejecutarComando lanza el shell por debajo para que las tuberías (|) y variables de entorno funcionen bien.
func ejecutarComando(comandoStr string) {
	cmd := exec.Command("sh", "-c", comandoStr)
	
	// Redirigimos la entrada, salida y errores estándar para interactuar (ej: meter la contraseña)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		// Salir silenciosamente porque el propio gestor de paquetes ya habrá impreso su error en pantalla
		os.Exit(1)
	}
}
