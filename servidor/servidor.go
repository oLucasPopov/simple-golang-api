package servidor

import (
	"encoding/json"
	"go_crud/banco"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    *uint32 `json:"id"`
	Nome  *string `json:"nome"`
	Email *string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body)

	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var usuario usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter usuario para struct"))
		return
	}

	db, erro := banco.Conectar()

	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("insert into usuarios(nome, email) values($1, $2) returning id")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement: " + erro.Error()))
		return
	}
	defer statement.Close()

	insersao, erro := statement.Query(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte("Erro ao executar o statement: " + erro.Error()))
		return
	}

	if insersao.Next() {
		insersao.Scan(&usuario.ID)

		if erro != nil {
			w.Write([]byte("Erro ao recuperar o id inserido: " + erro.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	usuarioCriado, _ := json.Marshal(usuario)
	w.Write([]byte(usuarioCriado))
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.Write([]byte("Erro ao converter para Uint: " + erro.Error()))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("select * from usuarios where id = $1")

	if err != nil {
		w.Write([]byte("Erro ao preparar query: " + err.Error()))
		return
	}

	res, err := statement.Query(ID)

	if err != nil {
		w.Write([]byte("Erro ao executar query: " + err.Error()))
		return
	}

	var usuario usuario

	if res.Next() {
		res.Scan(&usuario.ID, &usuario.Nome, &usuario.Email)
	}

	if usuario.ID == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.Write([]byte("Erro ao converter usuário para JSON"))
		return
	}
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()

	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	res, erro := db.Query("select * from usuarios")

	if erro != nil {
		w.Write([]byte("Erro ao buscar os usuários no banco de dados!"))
		return
	}
	defer res.Close()

	var usuarios []usuario

	for res.Next() {
		var usuario usuario
		if erro := res.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear o usuário!"))
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.Write([]byte("Erro ao converter usuário para JSON"))
		return
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	Id, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.Write([]byte("Erro ao pegar ID"))
		return
	}

	corpoRequisicao, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Erro ao pegar dados da requisicao"))
		return
	}

	var usuario usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter usuario para struct:" + erro.Error()))
		return
	}

	db, err := banco.Conectar()

	if err != nil {
		w.Write([]byte("Erro ao conectar no banco: " + err.Error()))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update usuarios set nome = $1, email = $2 where id = $3 returning id, nome, email")

	if err != nil {
		w.Write([]byte("Erro ao preparar query: " + err.Error()))
		return
	}

	rows, err := statement.Query(usuario.Nome, usuario.Email, Id)

	if err != nil {
		w.Write([]byte("Erro ao executar query: " + err.Error()))
		return
	}

	if rows.Next() {
		w.WriteHeader(http.StatusOK)
		rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email)
	}
	if erro := json.NewEncoder(w).Encode(&usuario); erro != nil {
		w.Write([]byte("Erro ao converter usuário para JSON"))
		return
	}
}

func ApagarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao dar parse no ID: " + err.Error()))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Erro ao conectar com o banco: " + err.Error()))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from usuarios where id = $1")
	if err != nil {
		w.Write([]byte("Erro ao preparar query: " + err.Error()))
		return
	}
	defer statement.Close()

	_, err = statement.Exec(id)

	if err != nil {
		w.Write([]byte("Erro ao deletar usuário: " + err.Error()))
		return
	}

	w.Write([]byte("Usuário excluído com sucesso!"))
}
