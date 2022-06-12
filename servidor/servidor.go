package servidor

import (
	"encoding/json"
	"fmt"
	"gocrud/banco"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    int32  `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// CriarUsuario insere um usuario no banco
func CriarUsuario(w http.ResponseWriter, r *http.Request) {

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisicao"))
		return
	}

	var usuario usuario
	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuario para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco"))
		return
	}

	// Prepare statement. sql injection.
	statement, erro := db.Prepare("insert into usuarios(nome, email) values (?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}
	defer statement.Close()

	insercao, erro := statement.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte("Erro ao criar o usuario"))
		return
	}
	db.Close()

	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao buscar o id do usuario inserido no banco"))
		return
	}

	// Status Code
	w.WriteHeader(http.StatusCreated)
	w.Write(
		[]byte(fmt.Sprintf("Usuario Inserido com Sucesso: id: %d", idInserido)),
	)
}

// RetornaUsuario busca um usuario no banco
func RetornaUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao converter ID em Int"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao Conectar no banco"))
		return
	}

	linha, erro := db.Query("select * from usuarios where id=?", ID)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao Buscar usuario no banco"))
		return
	}

	var usuario usuario
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro ao escanear usuario"))
			return
		}
	}

	if usuario.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro converter usuario para json"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RetornaUsuarios tras todos os usuarios do banco
func RetornaUsuarios(w http.ResponseWriter, r *http.Request) {

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco"))
		return
	}

	linhas, erro := db.Query("select * from usuarios")
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao realizar select no banco"))
		return
	}
	defer linhas.Close()
	db.Close()

	var usuarios []usuario
	for linhas.Next() {
		var usuario usuario

		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao scanear o usuario"))
			return
		}
		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)

	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.Write([]byte("Erro ao converter usuarios para JSON"))
	}
}

// Atualiza um usuario no banco
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter parametro para inteiro"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler o body"))
		return
	}

	var usuario usuario
	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o body em json"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco"))
		return
	}

	statement, erro := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement"))
		fmt.Println(erro.Error())
		return
	}

	if _, erro := statement.Exec(usuario.Nome, usuario.Email, ID); erro != nil {
		w.Write([]byte("Erro ao realizar o update no usuario"))
		return
	}
	defer statement.Close()
	defer db.Close()

	w.WriteHeader(http.StatusNoContent)
}

// DeletarUsuario remove um usuario do banco de dados
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o id para int"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco"))
		return
	}

	statement, erro := db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao deletar usuario"))
		return
	}
	defer statement.Close()
	defer db.Close()

	w.WriteHeader(http.StatusNoContent)
}
