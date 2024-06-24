package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Avaliacao struct {
	IDProduto   string       `json:"id_produto"`
	IDUsuario   string       `json:"id_usuario"`
	Nota        int          `json:"nota"`
	Comentario  string       `json:"comentario"`
	Comentarios []Comentario `json:"comentarios,omitempty"`
	Marcado     bool         `json:"marcado"`
	Moderado    bool         `json:"moderado"`
}

type Comentario struct {
	IDUsuario string `json:"id_usuario"`
	Texto     string `json:"texto"`
	Timestamp int64  `json:"timestamp"`
	Marcado   bool   `json:"marcado"`
	Moderado  bool   `json:"moderado"`
}

type Pergunta struct {
	IDProduto  string     `json:"id_produto"`
	IDUsuario  string     `json:"id_usuario"`
	IDVendedor string     `json:"id_vendedor"`
	Duvida     string     `json:"duvida"`
	Timestamp  int64      `json:"timestamp"`
	Respostas  []Resposta `json:"respostas,omitempty"`
	Marcado    bool       `json:"marcado"`
	Moderado   bool       `json:"moderado"`
}

type Resposta struct {
	IDVendedor string `json:"id_vendedor"`
	Resposta   string `json:"resposta"`
	Timestamp  int64  `json:"timestamp"`
	Marcado    bool   `json:"marcado"`
	Moderado   bool   `json:"moderado"`
}

type AvaliacoesProdutos struct {
	Avaliacoes []Avaliacao
	lock       sync.Mutex
}

type PerguntasProdutos struct {
	Perguntas []Pergunta
	lock      sync.Mutex
}

var palavrasBanidas []string
var lock sync.Mutex

func carregarPalavrasBanidas(caminhoArquivo string) error {
	lock.Lock()
	defer lock.Unlock()

	arquivo, err := os.Open(caminhoArquivo)
	if err != nil {
		return err
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		palavrasBanidas = append(palavrasBanidas, strings.ToLower(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func contemPalavrasBanidas(texto string) bool {
	lock.Lock()
	defer lock.Unlock()

	textoMinusculo := strings.ToLower(texto)
	for _, palavra := range palavrasBanidas {
		if strings.Contains(textoMinusculo, palavra) {
			return true
		}
	}
	return false
}

func (ap *AvaliacoesProdutos) AdicionarAvaliacao(avaliacao Avaliacao) error {
	if contemPalavrasBanidas(avaliacao.Comentario) {
		return fmt.Errorf("a avaliação contém palavras banidas")
	}
	ap.lock.Lock()
	defer ap.lock.Unlock()
	ap.Avaliacoes = append(ap.Avaliacoes, avaliacao)
	return nil
}

func (ap *AvaliacoesProdutos) AdicionarComentario(idProduto, idUsuario, texto string) error {
	if contemPalavrasBanidas(texto) {
		return fmt.Errorf("o comentário contém palavras banidas")
	}
	ap.lock.Lock()
	defer ap.lock.Unlock()
	for i, avaliacao := range ap.Avaliacoes {
		if avaliacao.IDProduto == idProduto {
			avaliacao.Comentarios = append(avaliacao.Comentarios, Comentario{
				IDUsuario: idUsuario,
				Texto:     texto,
				Timestamp: time.Now().Unix(),
			})
			ap.Avaliacoes[i] = avaliacao
			break
		}
	}
	return nil
}

func (pp *PerguntasProdutos) AdicionarPergunta(pergunta Pergunta) error {
	if contemPalavrasBanidas(pergunta.Duvida) {
		return fmt.Errorf("a pergunta contém palavras banidas")
	}
	pp.lock.Lock()
	defer pp.lock.Unlock()
	pp.Perguntas = append(pp.Perguntas, pergunta)
	return nil
}

func (pp *PerguntasProdutos) AdicionarResposta(idProduto, idUsuario, duvida, idVendedor, resposta string) error {
	if contemPalavrasBanidas(resposta) {
		return fmt.Errorf("a resposta contém palavras banidas")
	}
	pp.lock.Lock()
	defer pp.lock.Unlock()
	for i, pergunta := range pp.Perguntas {
		if pergunta.IDProduto == idProduto && pergunta.IDUsuario == idUsuario && pergunta.Duvida == duvida && pergunta.IDVendedor == idVendedor {
			pergunta.Respostas = append(pergunta.Respostas, Resposta{
				IDVendedor: idVendedor,
				Resposta:   resposta,
				Timestamp:  time.Now().Unix(),
			})
			pp.Perguntas[i] = pergunta
			break
		}
	}
	return nil
}

func (ap *AvaliacoesProdutos) MarcarAvaliacao(idProduto, idUsuario string) error {
	ap.lock.Lock()
	defer ap.lock.Unlock()
	for i, avaliacao := range ap.Avaliacoes {
		if avaliacao.IDProduto == idProduto && avaliacao.IDUsuario == idUsuario {
			ap.Avaliacoes[i].Marcado = true
			break
		}
	}
	return nil
}

func (pp *PerguntasProdutos) MarcarPergunta(idProduto, idUsuario string) error {
	pp.lock.Lock()
	defer pp.lock.Unlock()
	for i, pergunta := range pp.Perguntas {
		if pergunta.IDProduto == idProduto && pergunta.IDUsuario == idUsuario {
			pp.Perguntas[i].Marcado = true
			break
		}
	}
	return nil
}

func (ap *AvaliacoesProdutos) ModerarAvaliacao(idProduto, idUsuario, acao string) error {
	ap.lock.Lock()
	defer ap.lock.Unlock()
	for i, avaliacao := range ap.Avaliacoes {
		if avaliacao.IDProduto == idProduto && avaliacao.IDUsuario == idUsuario {
			if acao == "remover" {
				ap.Avaliacoes = append(ap.Avaliacoes[:i], ap.Avaliacoes[i+1:]...)
			} else if acao == "aprovar" {
				ap.Avaliacoes[i].Moderado = true
			}
			break
		}
	}
	return nil
}

func (pp *PerguntasProdutos) ModerarPergunta(idProduto, idUsuario, acao string) error {
	pp.lock.Lock()
	defer pp.lock.Unlock()
	for i, pergunta := range pp.Perguntas {
		if pergunta.IDProduto == idProduto && pergunta.IDUsuario == idUsuario {
			if acao == "remover" {
				pp.Perguntas = append(pp.Perguntas[:i], pp.Perguntas[i+1:]...)
			} else if acao == "aprovar" {
				pp.Perguntas[i].Moderado = true
			}
			break
		}
	}
	return nil
}

func (ap *AvaliacoesProdutos) ServirHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var avaliacao Avaliacao
		if err := json.NewDecoder(r.Body).Decode(&avaliacao); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ap.AdicionarAvaliacao(avaliacao); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Avaliação adicionada")
	case "GET":
		ap.lock.Lock()
		defer ap.lock.Unlock()
		avaliacoesJSON, err := json.Marshal(ap.Avaliacoes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(avaliacoesJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Método não permitido")
	}
}

func (ap *AvaliacoesProdutos) ServirModeracaoHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ap.lock.Lock()
		defer ap.lock.Unlock()
		avaliacoesMarcadas := []Avaliacao{}
		for _, avaliacao := range ap.Avaliacoes {
			if avaliacao.Marcado && !avaliacao.Moderado {
				avaliacoesMarcadas = append(avaliacoesMarcadas, avaliacao)
			}
		}
		avaliacoesJSON, err := json.Marshal(avaliacoesMarcadas)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(avaliacoesJSON)
	case "POST":
		var dados struct {
			IDProduto string `json:"id_produto"`
			IDUsuario string `json:"id_usuario"`
			Acao      string `json:"acao"` // "aprovar" ou "remover"
		}
		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		if err := ap.ModerarAvaliacao(dados.IDProduto, dados.IDUsuario, dados.Acao); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Avaliação moderada")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Método não permitido")
	}
}

func (pp *PerguntasProdutos) ServirModeracaoHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		pp.lock.Lock()
		defer pp.lock.Unlock()
		perguntasMarcadas := []Pergunta{}
		for _, pergunta := range pp.Perguntas {
			if pergunta.Marcado && !pergunta.Moderado {
				perguntasMarcadas = append(perguntasMarcadas, pergunta)
			}
		}
		perguntasJSON, err := json.Marshal(perguntasMarcadas)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(perguntasJSON)
	case "POST":
		var dados struct {
			IDProduto string `json:"id_produto"`
			IDUsuario string `json:"id_usuario"`
			Acao      string `json:"acao"` // "aprovar" ou "remover"
		}
		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		if err := pp.ModerarPergunta(dados.IDProduto, dados.IDUsuario, dados.Acao); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Pergunta moderada")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Método não permitido")
	}
}

func main() {
	if err := carregarPalavrasBanidas("palavras.txt"); err != nil {
		log.Fatalf("Falha ao carregar palavras banidas: %v", err)
	}

	var avaliacoes AvaliacoesProdutos
	var perguntas PerguntasProdutos

	http.HandleFunc("/avaliacoes", avaliacoes.ServirHTTP)
	http.HandleFunc("/avaliacoes/comentar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		var dados struct {
			IDProduto string `json:"id_produto"`
			IDUsuario string `json:"id_usuario"`
			Texto     string `json:"texto"`
		}
		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		if err := avaliacoes.AdicionarComentario(dados.IDProduto, dados.IDUsuario, dados.Texto); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Comentário adicionado")
	})

	http.HandleFunc("/perguntas/adicionar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		var pergunta Pergunta
		if err := json.NewDecoder(r.Body).Decode(&pergunta); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		if err := perguntas.AdicionarPergunta(pergunta); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Pergunta adicionada")
	})

	http.HandleFunc("/perguntas/responder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		var dados struct {
			IDProduto  string `json:"id_produto"`
			IDUsuario  string `json:"id_usuario"`
			Duvida     string `json:"duvida"`
			IDVendedor string `json:"id_vendedor"`
			Resposta   string `json:"resposta"`
		}
		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		if err := perguntas.AdicionarResposta(dados.IDProduto, dados.IDUsuario, dados.Duvida, dados.IDVendedor, dados.Resposta); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Resposta adicionada")
	})

	http.HandleFunc("/avaliacoes/moderar", avaliacoes.ServirModeracaoHTTP)
	http.HandleFunc("/perguntas/moderar", perguntas.ServirModeracaoHTTP)

	http.HandleFunc("/avaliacoes/marcar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		var dados struct {
			IDProduto string `json:"id_produto"`
			IDUsuario string `json:"id_usuario"`
		}
		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		if err := avaliacoes.MarcarAvaliacao(dados.IDProduto, dados.IDUsuario); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Avaliação marcada")
	})

	http.HandleFunc("/perguntas/marcar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		var dados struct {
			IDProduto string `json:"id_produto"`
			IDUsuario string `json:"id_usuario"`
		}
		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		if err := perguntas.MarcarPergunta(dados.IDProduto, dados.IDUsuario); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Pergunta marcada")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
