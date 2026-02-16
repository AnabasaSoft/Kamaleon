#!/bin/bash

_kml_completions() {
    local curr_arg
    curr_arg="${COMP_WORDS[COMP_CWORD]}"

    # Lista de todas las acciones y alias soportados por kml
    local comandos="in instalar install rm quitar remove up actualizar update dup dist-upgrade re refrescar refresh se buscar search sd bdesc sdesc info ver show li lista list ar autoremove huerfanos cl limpiar clean -h --help"

    # Solo autocompletamos si el usuario está escribiendo el primer argumento (después de 'kml')
    if [[ $COMP_CWORD -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "${comandos}" -- "${curr_arg}") )
    fi
    
    # Nota: El autocompletado de paquetes (el segundo argumento) se delega al comportamiento
    # por defecto del sistema para no ralentizar la consola haciendo consultas complejas.
}

# Le decimos a bash que use esta función para el comando 'kml'
complete -F _kml_completions kml