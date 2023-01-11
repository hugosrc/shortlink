package port

import "net/http"

type Auth interface {
	Authenticate(r *http.Request, w http.ResponseWriter) (string, error)
}
