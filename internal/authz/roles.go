// file: internal/authz/roles.go
package authz

// Permissoes define todas as permissões granulares possíveis no sistema.
// Usamos constantes para evitar erros de digitação e para ter autocompletar na IDE.
const (
	PermissaoObrasLer            = "obras:ler"
	PermissaoObrasEscrever       = "obras:escrever"
	PermissaoPessoalEscrever     = "pessoal:escrever"
	PermissaoSuprimentosLer      = "suprimentos:ler"
	PermissaoSuprimentosEscrever = "suprimentos:escrever"
)

// Papel define um nome de papel/função para um conjunto de permissões.
type Papel string

const (
	PapelAdmin          Papel = "ADMIN"
	PapelGerenteDeObras Papel = "GERENTE_OBRAS"
	PapelVisualizador   Papel = "VISUALIZADOR"
)

// mapaDePapeis associa cada Papel a uma lista de suas permissões.
var mapaDePapeis = map[Papel][]string{
	PapelGerenteDeObras: {
		PermissaoObrasLer,
		PermissaoObrasEscrever,
		PermissaoPessoalEscrever,
		PermissaoSuprimentosLer,
		PermissaoSuprimentosEscrever,
	},
	PapelVisualizador: {
		PermissaoObrasLer,
	},
	// O PapelAdmin é especial e terá todas as permissões.
}

// GetPermissoesParaPapel retorna a lista de permissões para um dado papel.
// Para o Admin, ele retorna todas as permissões existentes no sistema.
func GetPermissoesParaPapel(papel Papel) []string {
	if papel == PapelAdmin {
		// Junta todas as permissões de todos os papéis para o Admin.
		permissoesUnicas := make(map[string]struct{})
		for _, permissoesDoPapel := range mapaDePapeis {
			for _, p := range permissoesDoPapel {
				permissoesUnicas[p] = struct{}{}
			}
		}

		todasAsPermissoes := make([]string, 0, len(permissoesUnicas))
		for p := range permissoesUnicas {
			todasAsPermissoes = append(todasAsPermissoes, p)
		}
		return todasAsPermissoes
	}

	// Para outros papéis, apenas retorna do mapa.
	return mapaDePapeis[papel]
}
