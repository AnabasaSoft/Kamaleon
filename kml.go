package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const Version = "2.5.0"

// --- COLORES ---
const (
	Rojo     = "\033[0;31m"; Verde    = "\033[1;32m"; Amarillo = "\033[1;33m"
	Azul     = "\033[1;34m"; Cian     = "\033[1;36m"; Magenta  = "\033[1;35m"
	Blanco   = "\033[1;37m"; Negrita  = "\033[1m"; Reset    = "\033[0m"
)

type Translation struct {
	Header       string
	Detected     string
	Usage        string
	MainActions  string
	Options      string
	Executing    string
	ErrorRoot    string
	ErrorManager string
	ErrorAction  string
	NotaZypper   string
	Actions      map[string]string
}

var locales = map[string]Translation{
	"es": {
		Header: "kml - Envoltorio universal para gestores de paquetes", Detected: "Gestor actual detectado", Usage: "USO", MainActions: "ACCIONES PRINCIPALES", Options: "OPCIONES", Executing: "Ejecutando", ErrorRoot: "Error: No ejecutes kml con 'sudo'.", ErrorManager: "Error: No se detectó gestor.", ErrorAction: "Error: Acción '%s' no reconocida.", NotaZypper: "Zypper gestiona huérfanos con 'rm'.",
		Actions: map[string]string{"in":"Instala paquetes (repositorios oficiales y AUR).","rm":"Desinstala paquetes limpiando dependencias de forma segura.","up":"Actualiza el sistema y los repositorios (modo conservador).","dup":"Actualización profunda (full-upgrade, distro-sync o -Syyu).","re":"Refresca únicamente las bases de datos de los repositorios.","se":"Busca paquetes filtrando estrictamente por el NOMBRE.","sd":"Busca de forma amplia en el nombre y en la DESCRIPCIÓN.","info":"Muestra detalles técnicos, versión y peso de un paquete.","li":"Lista todos los paquetes instalados o busca entre ellos.","ar":"Rastrea y elimina dependencias huérfanas del sistema.","cl":"Vacía la caché de paquetes descargados para liberar espacio."},
	},
	"en": {
		Header: "kml - Universal package manager wrapper", Detected: "Current manager detected", Usage: "USAGE", MainActions: "MAIN ACTIONS", Options: "OPTIONS", Executing: "Executing", ErrorRoot: "Error: Do not run kml with 'sudo'.", ErrorManager: "Error: No manager detected.", ErrorAction: "Error: Action '%s' not recognized.", NotaZypper: "Note: Zypper manages orphans with 'rm'.",
		Actions: map[string]string{"in":"Install packages (official repos and AUR).","rm":"Remove packages and clean dependencies safely.","up":"Update system and repositories (conservative mode).","dup":"Deep update (full-upgrade, distro-sync, or -Syyu).","re":"Refresh repository databases only.","se":"Search packages strictly by NAME.","sd":"Broad search in both name and DESCRIPTION.","info":"Show technical details, version, and size of a package.","li":"List all installed packages or search through them.","ar":"Track and remove orphan dependencies from the system.","cl":"Clear downloaded package cache to free up space."},
	},
	"fr": {
		Header: "kml - Enveloppe universelle pour gestionnaires de paquets", Detected: "Gestionnaire actuel détecté", Usage: "UTILISATION", MainActions: "ACTIONS PRINCIPALES", Options: "OPTIONS", Executing: "Exécution", ErrorRoot: "Erreur : Ne lancez pas kml avec 'sudo'.", ErrorManager: "Erreur : Aucun gestionnaire détecté.", ErrorAction: "Erreur : Action '%s' non reconnue.", NotaZypper: "Note : Zypper gère les orphelins avec 'rm'.",
		Actions: map[string]string{"in":"Installe des paquets (dépôts officiels et AUR).","rm":"Supprime les paquets et nettoie les dépendances en toute sécurité.","up":"Met à jour le système et les dépôts (mode conservateur).","dup":"Mise à jour complète (full-upgrade ou -Syyu).","re":"Rafraîchit uniquement les bases de données des dépôts.","se":"Recherche des paquets strictement par NOM.","sd":"Recherche large dans le nom et la DESCRIPTION.","info":"Affiche les détails techniques, la version et la taille d'un paquet.","li":"Liste les paquets installés ou effectue une recherche locale.","ar":"Traque et supprime les dépendances orphelines du système.","cl":"Vide le cache des paquets téléchargés pour libérer de l'espace."},
	},
	"de": {
		Header: "kml - Universeller Paketmanager-Wrapper", Detected: "Aktueller Manager erkannt", Usage: "VERWENDUNG", MainActions: "HAUPTAKTIONEN", Options: "OPTIONEN", Executing: "Ausführen", ErrorRoot: "Fehler: kml nicht mit 'sudo' ausführen.", ErrorManager: "Fehler: Kein Manager erkannt.", ErrorAction: "Fehler: Aktion '%s' unbekannt.", NotaZypper: "Hinweis: Zypper verwaltet Waise mit 'rm'.",
		Actions: map[string]string{"in":"Installiert Pakete (offizielle Repos und AUR).","rm":"Entfernt Pakete und bereinigt Abhängigkeiten sicher.","up":"Aktualisiert das System und die Repos (konservativ).","dup":"Tiefgreifendes Update (full-upgrade oder -Syyu).","re":"Aktualisiert nur die Repository-Datenbanken.","se":"Sucht Pakete strikt nach NAMEN.","sd":"Breite Suche in Name und BESCHREIBUNG.","info":"Zeigt technische Details, Version und Größe eines Pakets.","li":"Listet installierte Pakete auf oder durchsucht diese.","ar":"Sucht und entfernt verwaiste Abhängigkeiten vom System.","cl":"Leert den Paket-Cache, um Speicherplatz freizugeben."},
	},
	"it": {
		Header: "kml - Wrapper universale per gestori di pacchetti", Detected: "Gestore attuale rilevato", Usage: "USO", MainActions: "AZIONI PRINCIPALI", Options: "OPZIONI", Executing: "Esecuzione", ErrorRoot: "Errore: Non eseguire kml con 'sudo'.", ErrorManager: "Errore: Nessun gestore rilevato.", ErrorAction: "Errore: Azione '%s' non riconosciuta.", NotaZypper: "Nota: Zypper gestisce gli orfani con 'rm'.",
		Actions: map[string]string{"in":"Installa pacchetti (repository ufficiali e AUR).","rm":"Rimuove i pacchetti e pulisce le dipendenze in sicurezza.","up":"Aggiorna il sistema e i repository (modalità conservativa).","dup":"Aggiornamento profondo (full-upgrade o -Syyu).","re":"Aggiorna solo i database dei repository.","se":"Cerca i pacchetti rigorosamente per NOME.","sd":"Ricerca ampia nel nome e nella DESCRIZIONE.","info":"Mostra dettagli tecnici, versione e dimensione di un pacchetto.","li":"Elenca tutti i pacchetti installati o cerca tra di essi.","ar":"Rintraccia e rimuove le dipendenze orfane dal sistema.","cl":"Svuota la cache dei pacchetti scaricati per liberare spazio."},
	},
	"pt": {
		Header: "kml - Wrapper universal para gestores de pacotes", Detected: "Gestor atual detectado", Usage: "USO", MainActions: "AÇÕES PRINCIPAIS", Options: "OPÇÕES", Executing: "Executando", ErrorRoot: "Erro: Não execute kml com 'sudo'.", ErrorManager: "Erro: Nenhum gestor detectado.", ErrorAction: "Erro: Ação '%s' não reconhecida.", NotaZypper: "Nota: Zypper gere órfãos com 'rm'.",
		Actions: map[string]string{"in":"Instala pacotes (repositórios oficiais e AUR).","rm":"Remove pacotes e limpa dependências com segurança.","up":"Atualiza o sistema e os repositórios (modo conservador).","dup":"Atualização profunda (full-upgrade ou -Syyu).","re":"Atualiza apenas as bases de dados dos repositórios.","se":"Procura pacotes estritamente pelo NOME.","sd":"Procura ampla no nome e na DESCRIÇÃO.","info":"Mostra detalhes técnicos, versão e tamanho de um pacote.","li":"Lista os pacotes instalados ou procura entre eles.","ar":"Rastreia e remove dependências órfãs do sistema.","cl":"Limpa a cache de pacotes baixados para libertar espaço."},
	},
	"ru": {
		Header: "kml - Универсальная оболочка для пакетных менеджеров", Detected: "Обнаружен менеджер", Usage: "ИСПОЛЬЗОВАНИЕ", MainActions: "ОСНОВНЫЕ ДЕЙСТВИЯ", Options: "ОПЦИИ", Executing: "Выполнение", ErrorRoot: "Ошибка: Не запускайте kml через 'sudo'.", ErrorManager: "Ошибка: Менеджер не найден.", ErrorAction: "Ошибка: Действие '%s' не опознано.", NotaZypper: "Примечание: Zypper удаляет сирот через 'rm'.",
		Actions: map[string]string{"in":"Установка пакетов (официальные репозитории и AUR).","rm":"Удаление пакетов и безопасная очистка зависимостей.","up":"Обновление системы и репозиториев (консервативный режим).","dup":"Полное обновление системы (full-upgrade или -Syyu).","re":"Обновление только баз данных репозиториев.","se":"Поиск пакетов строго по ИМЕНИ.","sd":"Широкий поиск по имени и ОПИСАНИЮ.","info":"Показать технические детали, версию и размер пакета.","li":"Список всех установленных пакетов или поиск по ним.","ar":"Поиск и удаление бесхозных зависимостей в системе.","cl":"Очистка кэша загруженных пакетов для освобождения места."},
	},
	"zh": {
		Header: "kml - 通用软件包管理器包装器", Detected: "已检测到管理器", Usage: "用法", MainActions: "主要动作", Options: "选项", Executing: "正在执行", ErrorRoot: "错误：请勿使用 'sudo' 运行 kml。", ErrorManager: "错误：未检测到包管理器。", ErrorAction: "错误：无法识别动作 '%s'。", NotaZypper: "注意：Zypper 使用 'rm' 管理孤立软件包。",
		Actions: map[string]string{"in":"安装软件包（官方仓库及 AUR）。","rm":"安全地卸载软件包并清理依赖关系。","up":"更新系统及仓库（保守模式）。","dup":"深度更新（full-upgrade 或 -Syyu）。","re":"仅刷新仓库数据库。","se":"仅根据名称精确搜索软件包。","sd":"在名称和描述中进行广泛搜索。","info":"显示软件包的技术详情、版本及大小。","li":"列出所有已安装的软件包或在其中搜索。","ar":"追踪并从系统中删除孤立依赖。","cl":"清理已下载的包缓存以释放空间。"},
	},
	"ja": {
		Header: "kml - ユニバーサルパッケージマネージャラッパー", Detected: "検出されたマネージャ", Usage: "使い方", MainActions: "主なアクション", Options: "オプション", Executing: "実行中", ErrorRoot: "エラー：kmlを'sudo'で実行しないでください。", ErrorManager: "エラー：マネージャが見つかりません。", ErrorAction: "エラー：アクション '%s' は無効です。", NotaZypper: "注意：Zypperは'rm'で孤立パッケージを管理します。",
		Actions: map[string]string{"in":"パッケージをインストール（公式リポジトリおよびAUR）。","rm":"パッケージを削除し、依存関係を安全にクリーンアップします。","up":"システムとリポジトリを更新します（保守モード）。","dup":"ディープアップデート（full-upgrade または -Syyu）。","re":"リポジトリのデータベースのみを更新します。","se":"名前でパッケージを厳密に検索します。","sd":"名前と説明の両方で幅広く検索します。","info":"パッケージの技術詳細、バージョン、サイズを表示します。","li":"インストール済みの全パッケージを表示または検索します。","ar":"システムから孤立した依存関係を追跡して削除します。","cl":"ダウンロードしたパッケージのキャッシュをクリアして空き容量を増やします。"},
	},
}

func main() {
	lang := "en"
	sysLang := os.Getenv("LANG")
	for k := range locales {
		if strings.HasPrefix(sysLang, k) {
			lang = k
			break
		}
	}
	t := locales[lang]

	usuarioActual, err := user.Current()
	if err == nil && usuarioActual.Uid == "0" {
		fmt.Printf("%s%s%s%s\n", Rojo, Negrita, t.ErrorRoot, Reset); os.Exit(1)
	}

	gestor := detectarGestor()
	if gestor == "" {
		fmt.Printf("%s%s%s\n", Rojo, t.ErrorManager, Reset); os.Exit(1)
	}

	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		mostrarAyuda(gestor, t); os.Exit(0)
	}

	accion := os.Args[1]
	paquetes := strings.Join(os.Args[2:], " ")

	comandoStr := mapearComando(gestor, accion, paquetes, t)
	if comandoStr == "" {
		fmt.Printf("%s%s%s\n", Rojo, fmt.Sprintf(t.ErrorAction, accion), Reset); os.Exit(1)
	}

	if !strings.HasPrefix(comandoStr, "echo") {
		fmt.Printf("%s%s==>%s %s: %s%s%s\n", Azul, Negrita, Reset, t.Executing, Negrita, comandoStr, Reset)
	}

	ejecutarComando(comandoStr)
}

func detectarGestor() string {
	// Añadimos apk y xbps-install a la lista de candidatos
	candidatos := []string{"yay", "paru", "pacman", "apt", "dnf", "zypper", "apk", "xbps-install"}
	for _, g := range candidatos {
		if _, err := exec.LookPath(g); err == nil {
			return g
		}
	}
	return ""
}

func mostrarAyuda(gestor string, t Translation) {
	fmt.Printf("\n%s%s%s %s(v%s)%s\n", Cian, Negrita, t.Header, Blanco, Version, Reset)
	fmt.Printf("%s: %s%s%s%s\n\n", t.Detected, Verde, Negrita, gestor, Reset)
	fmt.Printf("%s%s:%s\n  kml <acción> [paquetes...]\n\n", Amarillo, t.Usage, Reset)
	fmt.Printf("%s%s:%s\n", Amarillo, t.MainActions, Reset)

	printCmd := func(cmd, color, desc string) {
		fmt.Printf("  %s%-22s%s %s\n", color, cmd, Reset, desc)
	}

	printCmd("in, instalar", Verde, t.Actions["in"])
	printCmd("rm, quitar", Rojo, t.Actions["rm"])
	printCmd("up, actualizar", Cian, t.Actions["up"])
	printCmd("dup, dist-upgrade", Blanco, t.Actions["dup"])
	printCmd("re, refrescar", Magenta, t.Actions["re"])
	printCmd("se, buscar", Azul, t.Actions["se"])
	printCmd("sd, sdesc", Azul, t.Actions["sd"])
	printCmd("info, ver", Amarillo, t.Actions["info"])
	printCmd("li, lista", Negrita, t.Actions["li"])
	printCmd("ar, huerfanos", Rojo, t.Actions["ar"])
	printCmd("cl, limpiar", Rojo, t.Actions["cl"])

	fmt.Printf("\n%s%s:%s\n  -h, --help             %s\n\n", Amarillo, t.Options, Reset, "Muestra esta ayuda detallada.")
}

func mapearComando(gestor, accion, paq string, t Translation) string {
	esAUR := (gestor == "yay" || gestor == "paru")

	switch accion {
		case "in", "instalar", "install":
			if esAUR { return gestor + " -S " + paq }
			if gestor == "pacman" { return "sudo pacman -S " + paq }
			if gestor == "apk" { return "sudo apk add " + paq }
			if gestor == "xbps" { return "sudo xbps-install -S " + paq }
			return "sudo " + gestor + " install " + paq

		case "rm", "quitar", "remove":
			if esAUR { return gestor + " -Rs " + paq }
			if gestor == "pacman" { return "sudo pacman -Rs " + paq }
			if gestor == "apk" { return "sudo apk del " + paq }
			if gestor == "xbps" { return "sudo xbps-remove -R " + paq }
			return "sudo " + gestor + " remove " + paq

		case "up", "actualizar", "update":
			if esAUR { return gestor + " -Syu" }
			if gestor == "pacman" { return "sudo pacman -Syu" }
			if gestor == "apt" { return "sudo apt update && sudo apt upgrade" }
			if gestor == "apk" { return "sudo apk update && sudo apk upgrade" }
			if gestor == "xbps" { return "sudo xbps-install -Su" }
			return "sudo " + gestor + " upgrade" // dnf y zypper usan 'upgrade' o 'up'

		case "dup", "dist-upgrade":
			if esAUR { return gestor + " -Syyu" }
			if gestor == "pacman" { return "sudo pacman -Syyu" }
			if gestor == "apt" { return "sudo apt update && sudo apt full-upgrade" }
			if gestor == "dnf" { return "sudo dnf distro-sync" }
			if gestor == "apk" { return "sudo apk update && sudo apk upgrade -a" }
			if gestor == "xbps" { return "sudo xbps-install -Su" } // Void es rolling puro
			if gestor == "zypper" { return "sudo zypper dup" }
			return ""

		case "re", "refrescar", "refresh":
			if esAUR { return gestor + " -Sy" }
			if gestor == "pacman" { return "sudo pacman -Sy" }
			if gestor == "dnf" { return "sudo dnf makecache" }
			if gestor == "zypper" { return "sudo zypper ref" }
			if gestor == "apk" { return "sudo apk update" }
			if gestor == "xbps" { return "sudo xbps-install -S" }
			return "sudo apt update"

		case "se", "buscar", "search":
			// Búsqueda estricta (solo nombre)
			if esAUR || gestor == "pacman" { return gestor + " -Ss " + paq + ` | grep --color=auto -i -E -A 1 "^[^/]+/[^ ]*` + paq + `"` }
			if gestor == "apt" { return "apt search --names-only " + paq }
			if gestor == "dnf" { return "dnf list \"*" + paq + "*\" 2>/dev/null" }
			if gestor == "zypper" { return "zypper se -n " + paq }
			if gestor == "apk" { return "apk search " + paq }
			if gestor == "xbps" { return "xbps-query -Rs " + paq }
			return ""

		case "sd", "bdesc", "sdesc":
			// Búsqueda amplia
			if esAUR || gestor == "pacman" { return gestor + " -Ss " + paq }
			if gestor == "apk" { return "apk search -v -d " + paq }
			if gestor == "xbps" { return "xbps-query -Rs " + paq }
			return gestor + " search " + paq

		case "info", "ver", "show":
			if esAUR || gestor == "pacman" { return gestor + " -Si " + paq }
			if gestor == "apt" { return "apt show " + paq }
			if gestor == "apk" { return "apk info -d -s " + paq }
			if gestor == "xbps" { return "xbps-query -RS " + paq }
			return gestor + " info " + paq

		case "li", "lista", "list":
			if paq == "" {
				if esAUR || gestor == "pacman" { return gestor + " -Qe" }
				if gestor == "apt" { return "apt list --installed 2>/dev/null" }
				if gestor == "apk" { return "apk info" }
				if gestor == "xbps" { return "xbps-query -l" }
				if gestor == "zypper" { return "zypper se -i" }
				return gestor + " list installed"
			}
			if esAUR || gestor == "pacman" { return gestor + " -Q | grep --color=auto -i \"" + paq + "\"" }
			if gestor == "apk" { return "apk info | grep -i " + paq }
			if gestor == "xbps" { return "xbps-query -l | grep -i " + paq }
			return gestor + " list installed | grep -i " + paq

		case "ar", "autoremove", "huerfanos":
			if gestor == "yay" { return "yay -Yc" }
			if gestor == "paru" { return "paru -c" }
			if gestor == "pacman" { return "pacman -Qdtq | sudo xargs -r pacman -Rns" }
			if gestor == "zypper" { return "echo -e \"" + t.NotaZypper + "\"" }
			if gestor == "apk" { return "echo \"En Alpine se usa 'clean' para mantenimiento. Ejecuta: kml cl\"" }
			if gestor == "xbps" { return "sudo xbps-remove -O" }
			return "sudo " + gestor + " autoremove"

		case "cl", "limpiar", "clean":
			if esAUR { return gestor + " -Sc" }
			if gestor == "pacman" { return "sudo pacman -Sc" }
			if gestor == "dnf" { return "sudo dnf clean all" }
			if gestor == "apk" { return "sudo apk cache clean" }
			if gestor == "xbps" { return "sudo xbps-remove -O" }
			return "sudo " + gestor + " clean"
	}
	return ""
}

func ejecutarComando(comandoStr string) {
	cmd := exec.Command("sh", "-c", comandoStr)
	cmd.Stdin = os.Stdin; cmd.Stdout = os.Stdout; cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
