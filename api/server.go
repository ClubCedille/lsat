package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	Macaroon string `json:"macaroon" uri:"macaroon"`
	Invoice  string `json:"invoice" uri:"invoice"`
}

type Server struct {
	database map[string]int
}

func HttpServer() {
	fmt.Println("Server launched!")

	// Retourné dans le cas où la platforme reçoit un LSAT invalide
	http.HandleFunc("/", HandlePaymentRequest)

	// Ca devrait etre lié à la platforme?
	http.ListenAndServe(":8080", nil)
}

func RunServer() {
	engine := gin.New()

	// adding path params to router
	engine.GET("/auth/:macaroon/:invoice", func(context *gin.Context) {

		uri := PaymentRequest{}

		// binding to URI
		if err := context.BindUri(&uri); err != nil {

			context.AbortWithError(http.StatusBadRequest, err)

			return
		}

		fmt.Println(uri)

		context.JSON(http.StatusAccepted, &uri)
	})

	engine.Run(":8080")
}

func HandlePaymentRequest(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusPaymentRequired)
	w.Header().Set("WWW-Authenticate", "application/json") // A revoir le second argument

	// 1. Ici on devrait générer un invoice à l'aide de notre node.
	// 2. Si le client à son macaroon, on devrait aller le chercher dans la request, sinon on devrait en générer un.
	m := PaymentRequest{"", ""}
	b, _ := json.Marshal(m)

	w.Write(b) // Envoyer la requete JSON? Je suis habitué à devoir flush.
}
