package service

// The microservice is just over HTTP, so we just have a single transport.go.

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a profilesvc server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	r.Methods("GET").Path("/IntegrarDiario").Handler(httptransport.NewServer(
		e.IntegrarDiario,
		decodeGetConsetimentosAnteriores,
		encodeResponse,
		options...,
	))

	sh := http.StripPrefix("/apiDocs/", http.FileServer(http.Dir("./swagger/")))
	r.PathPrefix("/apiDocs/").Handler(sh)

	return r
}

func decodeGetConsetimentosAnteriores(_ context.Context, r *http.Request) (request interface{}, err error) {
	queryString := r.URL.Query()

	id := -1

	ids, containsID := queryString["id"]
	origems, containsOrigem := queryString["origem"]

	if containsID {
		offsetString := ids[0]
		id, err = strconv.Atoi(offsetString)
		if err != nil {
			return nil, errors.New("erro a ler id")
		}
	}

	if containsOrigem {
		return nil, errors.New("erro a ler origem")
	}

	return IntegrarDiarioRequest{
		Id:     id,
		Origem: origems[0],
	}, nil

}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//Comentando este pedaço de código, pois a resposta vai ser 200, com sucesso definido com falso
	//w.WriteHeader(codeFrom(err))
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		//como foi um erro, podemos afirmar que o sucesso é falso
		"success": false,
		//error é o resultado do erro gerado
		"error": err.Error(),
	})
}
