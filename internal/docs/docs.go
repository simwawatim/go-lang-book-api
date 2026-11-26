// Package docs serves the OpenAPI spec and a Swagger UI page.
package docs

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.yaml
var spec []byte

// Swagger UI is loaded from a CDN, so there is nothing extra to install.
const page = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Book API docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/openapi.yaml",
      dom_id: "#swagger-ui",
      persistAuthorization: true
    });
  </script>
</body>
</html>`

// Register adds GET /docs (the page) and GET /openapi.yaml (the spec).
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	})
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(spec)
	})
}
