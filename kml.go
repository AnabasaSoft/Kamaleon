package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const Version = "2.7.4"

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
	CmdUsage     string // Nueva: <acción> [paquetes...]
	MainActions  string
	Options      string
	Executing    string
	ErrorRoot    string
	ErrorManager string
	ErrorAction  string
	NotaZypper   string
	HelpDesc     string // Nueva: Muestra esta ayuda
	ResFlatpak   string // Nueva: Resultados en Flatpak
	InstFlatpak  string // Nueva: Flatpaks instalados
	EjemplosTit  string // Nueva: Título de ejemplos
	EjemplosTxt  string // Nueva: Texto de ejemplos
	Cmds         map[string]string // Nueva: Alias de comandos
	Actions      map[string]string
}

var locales = map[string]Translation{
	"es": {
		Header: "kml - Envoltorio universal para gestores de paquetes", Detected: "Gestor actual detectado", Usage: "USO", CmdUsage: "<acción> [paquetes...]", MainActions: "ACCIONES PRINCIPALES", Options: "OPCIONES", Executing: "Ejecutando", ErrorRoot: "Error: No ejecutes kml con 'sudo'.", ErrorManager: "Error: No se detectó gestor.", ErrorAction: "Error: Acción '%s' no reconocida.", NotaZypper: "Zypper gestiona huérfanos con 'rm'.", HelpDesc: "Muestra esta ayuda detallada.", ResFlatpak: "--- Resultados en Flatpak ---", InstFlatpak: "--- Flatpaks Instalados ---",
		EjemplosTit: "EJEMPLOS DE USO", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat paquetes.txt | kml rm",
		Cmds: map[string]string{"in":"in, instalar", "rm":"rm, quitar", "up":"up, actualizar", "dup":"dup, dist-upgrade", "re":"re, refrescar", "se":"se, buscar", "sd":"sd, sdesc", "info":"info, ver", "li":"li, lista", "is":"is, instalados", "ar":"ar, huerfanos", "cl":"cl, limpiar"},
		Actions: map[string]string{"in":"Instala paquetes (repositorios oficiales y AUR).","rm":"Desinstala paquetes limpiando dependencias de forma segura.","up":"Actualiza el sistema y los repositorios (modo conservador).","dup":"Actualización profunda (full-upgrade, distro-sync o -Syyu).","re":"Refresca únicamente las bases de datos de los repositorios.","se":"Busca paquetes filtrando estrictamente por el NOMBRE.","sd":"Busca de forma amplia en el nombre y en la DESCRIPCIÓN.","info":"Muestra detalles técnicos, versión y peso de un paquete.","li":"Lista todos los paquetes instalados o busca entre ellos.","is":"Busca un paquete solo entre los que ya tienes instalados.","ar":"Rastrea y elimina dependencias huérfanas del sistema.","cl":"Vacía la caché de paquetes descargados para liberar espacio."},
	},
	"en": {
		Header: "kml - Universal package manager wrapper", Detected: "Current manager detected", Usage: "USAGE", CmdUsage: "<action> [packages...]", MainActions: "MAIN ACTIONS", Options: "OPTIONS", Executing: "Executing", ErrorRoot: "Error: Do not run kml with 'sudo'.", ErrorManager: "Error: No manager detected.", ErrorAction: "Error: Action '%s' not recognized.", NotaZypper: "Note: Zypper manages orphans with 'rm'.", HelpDesc: "Show this detailed help.", ResFlatpak: "--- Flatpak Results ---", InstFlatpak: "--- Installed Flatpaks ---",
		EjemplosTit: "USAGE EXAMPLES", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat packages.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"Install packages (official repos and AUR).","rm":"Remove packages and clean dependencies safely.","up":"Update system and repositories (conservative mode).","dup":"Deep update (full-upgrade, distro-sync, or -Syyu).","re":"Refresh repository databases only.","se":"Search packages strictly by NAME.","sd":"Broad search in both name and DESCRIPTION.","info":"Show technical details, version, and size of a package.","li":"List all installed packages or search through them.","is":"Search for a package only among those already installed.","ar":"Track and remove orphan dependencies from the system.","cl":"Clear downloaded package cache to free up space."},
	},
	"fr": {
		Header: "kml - Enveloppe universelle pour gestionnaires de paquets", Detected: "Gestionnaire actuel détecté", Usage: "UTILISATION", CmdUsage: "<action> [paquets...]", MainActions: "ACTIONS PRINCIPALES", Options: "OPTIONS", Executing: "Exécution", ErrorRoot: "Erreur : Ne lancez pas kml avec 'sudo'.", ErrorManager: "Erreur : Aucun gestionnaire détecté.", ErrorAction: "Erreur : Action '%s' non reconnue.", NotaZypper: "Note : Zypper gère les orphelins avec 'rm'.", HelpDesc: "Affiche cette aide détaillée.", ResFlatpak: "--- Résultats Flatpak ---", InstFlatpak: "--- Flatpaks Installés ---",
		EjemplosTit: "EXEMPLES D'UTILISATION", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat paquets.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"Installe des paquets (dépôts officiels et AUR).","rm":"Supprime les paquets et nettoie les dépendances en toute sécurité.","up":"Met à jour le système et les dépôts (mode conservateur).","dup":"Mise à jour complète (full-upgrade ou -Syyu).","re":"Rafraîchit uniquement les bases de données des dépôts.","se":"Recherche des paquets strictement par NOM.","sd":"Recherche large dans le nom et la DESCRIPTION.","info":"Affiche les détails techniques, la version et la taille d'un paquet.","li":"Liste les paquets installés ou effectue une recherche locale.","is":"Recherche un paquet uniquement parmi ceux déjà installés.","ar":"Traque et supprime les dépendances orphelines du système.","cl":"Vide le cache des paquets téléchargés pour libérer de l'espace."},
	},
	"de": {
		Header: "kml - Universeller Paketmanager-Wrapper", Detected: "Aktueller Manager erkannt", Usage: "VERWENDUNG", CmdUsage: "<aktion> [pakete...]", MainActions: "HAUPTAKTIONEN", Options: "OPTIONEN", Executing: "Ausführen", ErrorRoot: "Fehler: kml nicht mit 'sudo' ausführen.", ErrorManager: "Fehler: Kein Manager erkannt.", ErrorAction: "Fehler: Aktion '%s' unbekannt.", NotaZypper: "Hinweis: Zypper verwaltet Waise mit 'rm'.", HelpDesc: "Zeigt diese detaillierte Hilfe an.", ResFlatpak: "--- Flatpak Ergebnisse ---", InstFlatpak: "--- Installierte Flatpaks ---",
		EjemplosTit: "ANWENDUNGSBEISPIELE", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat pakete.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"Installiert Pakete (offizielle Repos und AUR).","rm":"Entfernt Pakete und bereinigt Abhängigkeiten sicher.","up":"Aktualisiert das System und die Repos (konservativ).","dup":"Tiefgreifendes Update (full-upgrade oder -Syyu).","re":"Aktualisiert nur die Repository-Datenbanken.","se":"Sucht Pakete strikt nach NAMEN.","sd":"Breite Suche in Name und BESCHREIBUNG.","info":"Zeigt technische Details, Version und Größe eines Pakets.","li":"Listet installierte Pakete auf oder durchsucht diese.","is":"Sucht ein Paket nur unter den bereits installierten.","ar":"Sucht und entfernt verwaiste Abhängigkeiten vom System.","cl":"Leert den Paket-Cache, um Speicherplatz freizugeben."},
	},
	"it": {
		Header: "kml - Wrapper universale per gestori di pacchetti", Detected: "Gestore attuale rilevato", Usage: "USO", CmdUsage: "<azione> [pacchetti...]", MainActions: "AZIONI PRINCIPALI", Options: "OPZIONI", Executing: "Esecuzione", ErrorRoot: "Errore: Non eseguire kml con 'sudo'.", ErrorManager: "Errore: Nessun gestore rilevato.", ErrorAction: "Errore: Azione '%s' non riconosciuta.", NotaZypper: "Nota: Zypper gestisce gli orfani con 'rm'.", HelpDesc: "Mostra questo aiuto dettagliato.", ResFlatpak: "--- Risultati Flatpak ---", InstFlatpak: "--- Flatpak Installati ---",
		EjemplosTit: "ESEMPI DI UTILIZZO", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat pacchetti.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"Installa pacchetti (repository ufficiali e AUR).","rm":"Rimuove i pacchetti e pulisce le dipendenze in sicurezza.","up":"Aggiorna il sistema e i repository (modalità conservativa).","dup":"Aggiornamento profondo (full-upgrade o -Syyu).","re":"Aggiorna solo i database dei repository.","se":"Cerca i pacchetti rigorosamente per NOME.","sd":"Ricerca ampia nel nome e nella DESCRIZIONE.","info":"Mostra dettagli tecnici, versione e dimensione di un pacchetto.","li":"Elenca tutti i pacchetti installati o cerca tra di essi.","is":"Cerca un pacchetto solo tra quelli già installati.","ar":"Rintraccia e rimuove le dipendenze orfane dal sistema.","cl":"Svuota la cache dei pacchetti scaricati per liberare spazio."},
	},
	"pt": {
		Header: "kml - Wrapper universal para gestores de pacotes", Detected: "Gestor atual detectado", Usage: "USO", CmdUsage: "<ação> [pacotes...]", MainActions: "AÇÕES PRINCIPAIS", Options: "OPÇÕES", Executing: "Executando", ErrorRoot: "Erro: Não execute kml com 'sudo'.", ErrorManager: "Erro: Nenhum gestor detectado.", ErrorAction: "Erro: Ação '%s' não reconhecida.", NotaZypper: "Nota: Zypper gere órfãos com 'rm'.", HelpDesc: "Mostra esta ajuda detalhada.", ResFlatpak: "--- Resultados no Flatpak ---", InstFlatpak: "--- Flatpaks Instalados ---",
		EjemplosTit: "EXEMPLOS DE USO", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat pacotes.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"Instala pacotes (repositórios oficiais e AUR).","rm":"Remove pacotes e limpa dependências com segurança.","up":"Atualiza o sistema e os repositórios (modo conservador).","dup":"Atualização profunda (full-upgrade ou -Syyu).","re":"Atualiza apenas as bases de dados dos repositórios.","se":"Procura pacotes estritamente pelo NOME.","sd":"Procura ampla no nome e na DESCRIÇÃO.","info":"Mostra detalhes técnicos, versão e tamanho de um pacote.","li":"Lista os pacotes instalados ou procura entre eles.","is":"Procura um pacote apenas entre os já instalados.","ar":"Rastreia e remove dependências órfãs do sistema.","cl":"Limpa a cache de pacotes baixados para libertar espaço."},
	},
	"ru": {
		Header: "kml - Универсальная оболочка для пакетных менеджеров", Detected: "Обнаружен менеджер", Usage: "ИСПОЛЬЗОВАНИЕ", CmdUsage: "<действие> [пакеты...]", MainActions: "ОСНОВНЫЕ ДЕЙСТВИЯ", Options: "ОПЦИИ", Executing: "Выполнение", ErrorRoot: "Ошибка: Не запускайте kml через 'sudo'.", ErrorManager: "Ошибка: Менеджер не найден.", ErrorAction: "Ошибка: Действие '%s' не опознано.", NotaZypper: "Примечание: Zypper удаляет сирот через 'rm'.", HelpDesc: "Показать эту подробную справку.", ResFlatpak: "--- Результаты Flatpak ---", InstFlatpak: "--- Установленные Flatpak ---",
		EjemplosTit: "ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat packages.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"Установка пакетов (официальные репозитории и AUR).","rm":"Удаление пакетов и безопасная очистка зависимостей.","up":"Обновление системы и репозиториев (консервативный режим).","dup":"Полное обновление системы (full-upgrade или -Syyu).","re":"Обновление только баз данных репозиториев.","se":"Поиск пакетов строго по ИМЕНИ.","sd":"Широкий поиск по имени и ОПИСАНИЮ.","info":"Показать технические детали, версию и размер пакета.","li":"Список всех установленных пакетов или поиск по ним.","is":"Поиск пакета только среди уже установленных.","ar":"Поиск и удаление бесхозных зависимостей в системе.","cl":"Очистка кэша загруженных пакетов для освобождения места."},
	},
	"zh": {
		Header: "kml - 通用软件包管理器包装器", Detected: "已检测到管理器", Usage: "用法", CmdUsage: "<动作> [软件包...]", MainActions: "主要动作", Options: "选项", Executing: "正在执行", ErrorRoot: "错误：请勿使用 'sudo' 运行 kml。", ErrorManager: "错误：未检测到包管理器。", ErrorAction: "错误：无法识别动作 '%s'。", NotaZypper: "注意：Zypper 使用 'rm' 管理孤立软件包。", HelpDesc: "显示此详细帮助。", ResFlatpak: "--- Flatpak 搜索结果 ---", InstFlatpak: "--- 已安装的 Flatpak ---",
		EjemplosTit: "使用示例", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat packages.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"安装软件包（官方仓库及 AUR）。","rm":"安全地卸载软件包并清理依赖关系。","up":"更新系统及仓库（保守模式）。","dup":"深度更新（full-upgrade 或 -Syyu）。","re":"仅刷新仓库数据库。","se":"仅根据名称精确搜索软件包。","sd":"在名称和描述中进行广泛搜索。","info":"显示软件包的技术详情、版本及大小。","li":"列出所有已安装的软件包或在其中搜索。","is":"仅在已安装的软件包中搜索。","ar":"追踪并从系统中删除孤立依赖。","cl":"清理已下载的包缓存以释放空间。"},
	},
	"ja": {
		Header: "kml - ユニバーサルパッケージマネージャラッパー", Detected: "検出されたマネージャ", Usage: "使い方", CmdUsage: "<アクション> [パッケージ...]", MainActions: "主なアクション", Options: "オプション", Executing: "実行中", ErrorRoot: "エラー：kmlを'sudo'で実行しないでください。", ErrorManager: "エラー：マネージャが見つかりません。", ErrorAction: "エラー：アクション '%s' は無効です。", NotaZypper: "注意：Zypperは'rm'で孤立パッケージを管理します。", HelpDesc: "この詳細なヘルプを表示します。", ResFlatpak: "--- Flatpak の結果 ---", InstFlatpak: "--- インストール済みの Flatpak ---",
		EjemplosTit: "使用例", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat packages.txt | kml rm",
		Cmds: map[string]string{"in":"in, install", "rm":"rm, remove", "up":"up, update", "dup":"dup, dist-upgrade", "re":"re, refresh", "se":"se, search", "sd":"sd, sdesc", "info":"info, show", "li":"li, list", "is":"is, installed", "ar":"ar, autoremove", "cl":"cl, clean"},
		Actions: map[string]string{"in":"パッケージをインストール（公式リポジトリおよびAUR）。","rm":"パッケージを削除し、依存関係を安全にクリーンアップします。","up":"システムとリポジトリを更新します（保守モード）。","dup":"ディープアップデート（full-upgrade または -Syyu）。","re":"リポジトリのデータベースのみを更新します。","se":"名前でパッケージを厳密に検索します。","sd":"名前と説明の両方で幅広く検索します。","info":"パッケージの技術詳細、バージョン、サイズを表示します。","li":"インストール済みの全パッケージを表示または検索します。","is":"すでにインストールされているパッケージの中からのみ検索します。","ar":"システムから孤立した依存関係を追跡して削除します。","cl":"ダウンロードしたパッケージのキャッシュをクリアして空き容量を増やします。"},
	},
	"eu": {
		Header: "kml - Pakete-kudeatzaileentzako bilgarri unibertsala", Detected: "Detektatutako kudeatzailea", Usage: "ERABILERA", CmdUsage: "<ekintza> [paketeak...]", MainActions: "EKINTZA NAGUSIAK", Options: "AUKERAK", Executing: "Exekutatzen", ErrorRoot: "Errorea: Ez exekutatu kml 'sudo'-rekin.", ErrorManager: "Errorea: Ez da kudeatzailerik detektatu.", ErrorAction: "Errorea: '%s' ekintza ezezaguna.", NotaZypper: "Oharra: Zypper-ek umezurtzak 'rm'-rekin kudeatzen ditu.", HelpDesc: "Laguntza zehatz hau erakusten du.", ResFlatpak: "--- Emaitzak Flatpak-en ---", InstFlatpak: "--- Instalatutako Flatpak-ak ---",
		EjemplosTit: "ERABILERA ADIBIDEAK", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat paketeak.txt | kml rm",
		Cmds: map[string]string{"in":"in, instalar", "rm":"rm, quitar", "up":"up, actualizar", "dup":"dup, dist-upgrade", "re":"re, refrescar", "se":"se, buscar", "sd":"sd, sdesc", "info":"info, ver", "li":"li, lista", "is":"is, instalados", "ar":"ar, huerfanos", "cl":"cl, limpiar"},
		Actions: map[string]string{"in":"Paketeak instalatzen ditu (biltegi ofizialak eta AUR).","rm":"Paketeak desinstalatzen ditu, menpekotasunak modu seguruan garbituz.","up":"Sistema eta biltegiak eguneratzen ditu (modu kontserbadorea).","dup":"Eguneratze sakona (full-upgrade, distro-sync edo -Syyu).","re":"Biltegien datu-baseak soilik freskatzen ditu.","se":"Paketeak bilatzen ditu, IZENAREN arabera zorrotz iragaziz.","sd":"Izenaren eta DESKRIBAPENAREN arabera bilaketa zabala egiten du.","info":"Pakete baten xehetasun teknikoak, bertsioa eta pisua erakusten ditu.","li":"Instalatutako pakete guztiak zerrendatzen ditu edo haien artean bilatzen du.","is":"Dagoeneko instalatuta dauden paketeen artean bakarrik bilatzen du.","ar":"Sistemako menpekotasun umezurtzak arakatu eta ezabatzen ditu.","cl":"Deskargatutako paketeen katxea husten du lekua askatzeko."},
	},
	"ca": {
		Header: "kml - Embolcall universal per a gestors de paquets", Detected: "Gestor actual detectat", Usage: "ÚS", CmdUsage: "<acció> [paquets...]", MainActions: "ACCIONS PRINCIPALS", Options: "OPCIONS", Executing: "Executant", ErrorRoot: "Error: No executis kml amb 'sudo'.", ErrorManager: "Error: No s'ha detectat cap gestor.", ErrorAction: "Error: Acció '%s' no reconeguda.", NotaZypper: "Nota: Zypper gestiona els orfes amb 'rm'.", HelpDesc: "Mostra aquesta ajuda detallada.", ResFlatpak: "--- Resultats a Flatpak ---", InstFlatpak: "--- Flatpaks Instal·lats ---",
		EjemplosTit: "EXEMPLES D'ÚS", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat paquets.txt | kml rm",
		Cmds: map[string]string{"in":"in, instalar", "rm":"rm, quitar", "up":"up, actualizar", "dup":"dup, dist-upgrade", "re":"re, refrescar", "se":"se, buscar", "sd":"sd, sdesc", "info":"info, ver", "li":"li, lista", "is":"is, instalados", "ar":"ar, huerfanos", "cl":"cl, limpiar"},
		Actions: map[string]string{"in":"Instal·la paquets (repositoris oficials i AUR).","rm":"Desinstal·la paquets netejant dependències de manera segura.","up":"Actualitza el sistema i els repositoris (mode conservador).","dup":"Actualització profunda (full-upgrade, distro-sync o -Syyu).","re":"Refresca únicament les bases de dades dels repositoris.","se":"Cerca paquets filtrant estrictament pel NOM.","sd":"Cerca de manera àmplia en el nom i en la DESCRIPCIÓ.","info":"Mostra detalls tècnics, versió i pes d'un paquet.","li":"Llista tots els paquets instal·lats o cerca entre ells.","is":"Cerca un paquet només entre els que ja tens instal·lats.","ar":"Rastreja i elimina dependències òrfenes del sistema.","cl":"Buida la memòria cau de paquets descarregats per alliberar espai."},
	},
	"gl": {
		Header: "kml - Envoltorio universal para xestores de paquetes", Detected: "Xestor actual detectado", Usage: "USO", CmdUsage: "<acción> [paquetes...]", MainActions: "ACCIÓNS PRINCIPAIS", Options: "OPCIÓNS", Executing: "Executando", ErrorRoot: "Erro: Non executes kml con 'sudo'.", ErrorManager: "Erro: Non se detectou xestor.", ErrorAction: "Erro: Acción '%s' non recoñecida.", NotaZypper: "Nota: Zypper xestiona os orfos con 'rm'.", HelpDesc: "Mostra esta axuda detallada.", ResFlatpak: "--- Resultados en Flatpak ---", InstFlatpak: "--- Flatpaks Instalados ---",
		EjemplosTit: "EXEMPLOS DE USO", EjemplosTxt: "  kml in firefox\n  kml se \"gnome shell\"\n  echo \"htop vim\" | kml in\n  cat paquetes.txt | kml rm",
		Cmds: map[string]string{"in":"in, instalar", "rm":"rm, quitar", "up":"up, actualizar", "dup":"dup, dist-upgrade", "re":"re, refrescar", "se":"se, buscar", "sd":"sd, sdesc", "info":"info, ver", "li":"li, lista", "is":"is, instalados", "ar":"ar, huerfanos", "cl":"cl, limpiar"},
		Actions: map[string]string{"in":"Instala paquetes (repositorios oficiais e AUR).","rm":"Desinstala paquetes limpando dependencias de forma segura.","up":"Actualiza o sistema e os repositorios (modo conservador).","dup":"Actualización profunda (full-upgrade, distro-sync ou -Syyu).","re":"Refresca unicamente as bases de datos dos repositorios.","se":"Busca paquetes filtrando estritamente polo NOME.","sd":"Busca de forma ampla no nome e na DESCRICIÓN.","info":"Mostra detalles técnicos, versión e peso dun paquete.","li":"Lista todos os paquetes instalados ou busca entre eles.","is":"Busca un paquete só entre os que xa tes instalados.","ar":"Rastrexar e elimina dependencias orfas do sistema.","cl":"Baleira a caché de paquetes descargados para liberar espazo."},
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

	// Detectar si hay datos entrando por un pipe (ej: echo "firefox" | kml in)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		datosPipe, _ := io.ReadAll(os.Stdin)
		textoPipe := strings.TrimSpace(string(datosPipe))
		textoPipe = strings.ReplaceAll(textoPipe, "\n", " ") // Convertir multilinea en espacios

		if paquetes == "" {
			paquetes = textoPipe
		} else if textoPipe != "" {
			paquetes = paquetes + " " + textoPipe
		}
	}

	comandoStr := mapearComando(gestor, accion, paquetes, t)
	if comandoStr == "" {
		fmt.Printf("%s%s%s\n", Rojo, fmt.Sprintf(t.ErrorAction, accion), Reset); os.Exit(1)
	}

	if !strings.HasPrefix(comandoStr, "echo") && !strings.HasPrefix(comandoStr, "printf") {
		fmt.Printf("%s%s==>%s %s: %s%s%s\n", Azul, Negrita, Reset, t.Executing, Negrita, comandoStr, Reset)
	}

	ejecutarComando(comandoStr)
}

func detectarGestor() string {
	candidatos := []string{"yay", "paru", "pacman", "apt", "dnf", "zypper", "apk", "xbps-install"}
	for _, g := range candidatos {
		if _, err := exec.LookPath(g); err == nil {
			if g == "xbps-install" {
				return "xbps"
			}
			return g
		}
	}
	return ""
}

func mostrarAyuda(gestor string, t Translation) {
	fmt.Printf("\n%s%s%s %s(v%s)%s\n", Cian, Negrita, t.Header, Blanco, Version, Reset)
	fmt.Printf("%s: %s%s%s%s\n\n", t.Detected, Verde, Negrita, gestor, Reset)
	fmt.Printf("%s%s:%s\n  kml %s\n\n", Amarillo, t.Usage, Reset, t.CmdUsage)
	fmt.Printf("%s%s:%s\n", Amarillo, t.MainActions, Reset)

	printCmd := func(cmd, color, desc string) {
		fmt.Printf("  %s%-22s%s %s\n", color, cmd, Reset, desc)
	}

	printCmd(t.Cmds["in"], Verde, t.Actions["in"])
	printCmd(t.Cmds["rm"], Rojo, t.Actions["rm"])
	printCmd(t.Cmds["up"], Cian, t.Actions["up"])
	printCmd(t.Cmds["dup"], Blanco, t.Actions["dup"])
	printCmd(t.Cmds["re"], Magenta, t.Actions["re"])
	printCmd(t.Cmds["se"], Azul, t.Actions["se"])
	printCmd(t.Cmds["sd"], Azul, t.Actions["sd"])
	printCmd(t.Cmds["info"], Amarillo, t.Actions["info"])
	printCmd(t.Cmds["li"], Negrita, t.Actions["li"])
	printCmd(t.Cmds["is"], Verde, t.Actions["is"])
	printCmd(t.Cmds["ar"], Rojo, t.Actions["ar"])
	printCmd(t.Cmds["cl"], Rojo, t.Actions["cl"])

	fmt.Printf("\n%s%s:%s\n  -h, --help             %s\n", Amarillo, t.Options, Reset, t.HelpDesc)
	fmt.Printf("\n%s%s:%s\n%s\n\n", Amarillo, t.EjemplosTit, Reset, t.EjemplosTxt)
}

func mapearComando(gestor, accion, paq string, t Translation) string {
	esAUR := (gestor == "yay" || gestor == "paru")

	tieneFlatpak := false
	if _, err := exec.LookPath("flatpak"); err == nil {
		tieneFlatpak = true
	}

	var cmd string

	switch accion {
		case "in", "instalar", "install":
			if esAUR { cmd = gestor + " -S " + paq } else if gestor == "pacman" { cmd = "sudo pacman -S " + paq } else if gestor == "apk" { cmd = "sudo apk add " + paq } else if gestor == "xbps" { cmd = "sudo xbps-install -S " + paq } else { cmd = "sudo " + gestor + " install " + paq }
			if tieneFlatpak { return cmd + " || flatpak install " + paq }
			return cmd

		case "rm", "quitar", "remove":
			if esAUR { cmd = gestor + " -Rs " + paq } else if gestor == "pacman" { cmd = "sudo pacman -Rs " + paq } else if gestor == "apk" { cmd = "sudo apk del " + paq } else if gestor == "xbps" { cmd = "sudo xbps-remove -R " + paq } else { cmd = "sudo " + gestor + " remove " + paq }
			if tieneFlatpak { return cmd + " ; flatpak uninstall " + paq + " 2>/dev/null" }
			return cmd

		case "up", "actualizar", "update":
			extra := ""
			if paq != "" { extra = " " + paq }

			switch gestor {
				case "yay", "paru": cmd = gestor + " -Syu" + extra
				case "pacman": cmd = "sudo pacman -Syu" + extra
				case "apt": cmd = "sudo apt update && sudo apt upgrade" + extra
				case "apk": cmd = "sudo apk update && sudo apk upgrade" + extra
				case "xbps": cmd = "sudo xbps-install -Su" + extra
				default: cmd = "sudo " + gestor + " upgrade" + extra
			}
			if tieneFlatpak { return cmd + " && flatpak update" }
			return cmd

				case "dup", "dist-upgrade":
					switch gestor {
						case "yay", "paru": cmd = gestor + " -Syyu"
						case "pacman": cmd = "sudo pacman -Syyu"
						case "apt": cmd = "sudo apt update && sudo apt full-upgrade"
						case "dnf": cmd = "sudo dnf distro-sync"
						case "apk": cmd = "sudo apk update && sudo apk upgrade -a"
						case "xbps": cmd = "sudo xbps-install -Su"
						case "zypper": cmd = "sudo zypper dup"
						default: return ""
					}
					if tieneFlatpak { return cmd + " && flatpak update" }
					return cmd

						case "re", "refrescar", "refresh":
							if esAUR { cmd = gestor + " -Sy" } else if gestor == "pacman" { cmd = "sudo pacman -Sy" } else if gestor == "dnf" { cmd = "sudo dnf makecache" } else if gestor == "zypper" { cmd = "sudo zypper ref" } else if gestor == "apk" { cmd = "sudo apk update" } else if gestor == "xbps" { cmd = "sudo xbps-install -S" } else { cmd = "sudo apt update" }
							return cmd

						case "se", "buscar", "search":
							if esAUR || gestor == "pacman" { cmd = gestor + " -Ss " + paq + ` | grep --color=auto -i -E -A 1 "^[^/]+/[^ ]*` + paq + `"` } else if gestor == "apt" { cmd = "apt search --names-only " + paq } else if gestor == "dnf" { cmd = "dnf list \"*" + paq + "*\" 2>/dev/null" } else if gestor == "zypper" { cmd = "zypper se -n " + paq } else if gestor == "apk" { cmd = "apk search " + paq } else if gestor == "xbps" { cmd = "xbps-query -Rs " + paq }
							if tieneFlatpak && cmd != "" { return cmd + ` ; printf "\n\033[36m` + t.ResFlatpak + `\033[0m\n" ; flatpak search ` + paq }
							return cmd

						case "sd", "bdesc", "sdesc":
							if esAUR || gestor == "pacman" { cmd = gestor + " -Ss " + paq } else if gestor == "apk" { cmd = "apk search -v -d " + paq } else if gestor == "xbps" { cmd = "xbps-query -Rs " + paq } else { cmd = gestor + " search " + paq }
							if tieneFlatpak { return cmd + ` ; printf "\n\033[36m` + t.ResFlatpak + `\033[0m\n" ; flatpak search ` + paq }
							return cmd

						case "info", "ver", "show":
							if esAUR || gestor == "pacman" { cmd = gestor + " -Si " + paq } else if gestor == "apt" { cmd = "apt show " + paq } else if gestor == "apk" { cmd = "apk info -d -s " + paq } else if gestor == "xbps" { cmd = "xbps-query -RS " + paq } else { cmd = gestor + " info " + paq }
							if tieneFlatpak { return cmd + " ; flatpak info " + paq + " 2>/dev/null" }
							return cmd

						case "li", "lista", "list", "is", "instalados", "installed":
							if paq == "" {
								if esAUR || gestor == "pacman" { cmd = gestor + " -Qe" } else if gestor == "apt" { cmd = "apt list --installed 2>/dev/null" } else if gestor == "apk" { cmd = "apk info" } else if gestor == "xbps" { cmd = "xbps-query -l" } else if gestor == "zypper" { cmd = "zypper se -i" } else { cmd = gestor + " list installed" }
								if tieneFlatpak { return cmd + ` ; printf "\n\033[36m` + t.InstFlatpak + `\033[0m\n" ; flatpak list` }
								return cmd
							}
							if esAUR || gestor == "pacman" { cmd = gestor + " -Q | grep --color=auto -i \"" + paq + "\"" } else if gestor == "apk" { cmd = "apk info | grep -i " + paq } else if gestor == "xbps" { cmd = "xbps-query -l | grep -i " + paq } else { cmd = gestor + " list installed | grep -i " + paq }
							if tieneFlatpak { return cmd + " ; flatpak list | grep -i " + paq }
							return cmd

						case "ar", "autoremove", "huerfanos":
							if gestor == "yay" { cmd = "yay -Yc" } else if gestor == "paru" { cmd = "paru -c" } else if gestor == "pacman" { cmd = "pacman -Qdtq | sudo xargs -r pacman -Rns" } else if gestor == "zypper" { cmd = "echo -e \"" + t.NotaZypper + "\"" } else if gestor == "apk" { cmd = "echo \"En Alpine se usa 'clean' para mantenimiento.\"" } else if gestor == "xbps" { cmd = "sudo xbps-remove -O" } else { cmd = "sudo " + gestor + " autoremove" }
							if tieneFlatpak { return cmd + " ; flatpak uninstall --unused 2>/dev/null" }
							return cmd

						case "cl", "limpiar", "clean":
							if esAUR { cmd = gestor + " -Sc" } else if gestor == "pacman" { cmd = "sudo pacman -Sc" } else if gestor == "dnf" { cmd = "sudo dnf clean all" } else if gestor == "apk" { cmd = "sudo apk cache clean" } else if gestor == "xbps" { cmd = "sudo xbps-remove -O" } else { cmd = "sudo " + gestor + " clean" }
							if tieneFlatpak { return cmd + " ; flatpak uninstall --unused" }
							return cmd
	}
	return ""
}

func ejecutarComando(comandoStr string) {
	cmd := exec.Command("sh", "-c", comandoStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Si venimos de un pipe, el stdin original está agotado.
	// Reconectamos el comando hijo directamente al teclado del usuario (/dev/tty)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		tty, err := os.Open("/dev/tty")
		if err == nil {
			cmd.Stdin = tty
			defer tty.Close()
		} else {
			cmd.Stdin = os.Stdin // Fallback por si acaso
		}
	} else {
		cmd.Stdin = os.Stdin // Comportamiento normal
	}

	_ = cmd.Run()
}
