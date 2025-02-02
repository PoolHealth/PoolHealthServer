package google_oauth

import "net/http"

type Handler struct {
	adapter adapter
}

type adapter interface {
	ApplyAuthCode(authCode string) error
}

func NewHandler(adapter adapter) http.Handler {
	return &Handler{adapter: adapter}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	authCode := request.URL.Query().Get("code")
	if authCode == "" {
		http.Error(writer, "missing code", http.StatusBadRequest)
		return
	}

	if err := h.adapter.ApplyAuthCode(authCode); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte("<html><h1>Ok, you can close the page</h1></html>")); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
