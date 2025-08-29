package main


import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var authToken = getAuthToken()

func getAuthToken() string {
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		return "seu_token_fixo_aqui"
	}
	return token
}

// Estrutura para gerenciar clientes por canal


var channels = make(map[string][]chan string)
var messagesPublished int


// Middleware simples de autenticação
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer "+authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Handler SSE
func sseHandler(w http.ResponseWriter, r *http.Request) {
	canal := r.URL.Query().Get("canal")
	if canal == "" {
		http.Error(w, "Canal não especificado", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
		return
	}

	msgChan := make(chan string)
	channels[canal] = append(channels[canal], msgChan)
	defer func() {
		// Remove canal do slice ao desconectar
		for i, c := range channels[canal] {
			if c == msgChan {
				channels[canal] = append(channels[canal][:i], channels[canal][i+1:]...)
				break
			}
		}
		close(msgChan)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-msgChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// Handler de publicação
func publishHandler(w http.ResponseWriter, r *http.Request) {
	canal := r.URL.Query().Get("canal")
	msg := r.URL.Query().Get("msg")
	if canal == "" || msg == "" {
		http.Error(w, "Canal ou mensagem não especificados", http.StatusBadRequest)
		return
	}
	for _, client := range channels[canal] {
		select {
		case client <- fmt.Sprintf("%s [%s]", msg, time.Now().Format(time.RFC3339)):
		default:
		}
	}
	messagesPublished++
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mensagem publicada"))
}

// Handler de estatísticas em HTML
func statsHandler(w http.ResponseWriter, r *http.Request) {
	activeChannels := len(channels)
	connectedClients := 0
	for _, clients := range channels {
		connectedClients += len(clients)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<html><head><title>Stats SSE Server</title></head><body style='font-family:sans-serif;padding:32px;'>`)
	fmt.Fprintf(w, `<h1>Estatísticas do SSE Server</h1>`)
	fmt.Fprintf(w, `<p><strong>Canais ativos:</strong> %d</p>`, activeChannels)
	fmt.Fprintf(w, `<p><strong>Clientes conectados:</strong> %d</p>`, connectedClients)
	fmt.Fprintf(w, `<p><strong>Mensagens publicadas:</strong> %d</p>`, messagesPublished)
	fmt.Fprintf(w, `</body></html>`)
}

func main() {
	http.HandleFunc("/sse", sseHandler)
	http.HandleFunc("/publish", authMiddleware(publishHandler))
	http.HandleFunc("/stats", statsHandler)
	log.Println("Servidor SSE rodando em :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
