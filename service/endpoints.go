package service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	stdjwt "github.com/golang-jwt/jwt/v4"
	"os"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	IntegrarDiario endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		IntegrarDiario: MakeIntegrarDiario(s),
	}
}

// Esta função faz verificações preliminares, devolvendo a chave privada para a verificação do JWT.
// A chave é obtida da variável de ambiente
func keyFunc(token *stdjwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*stdjwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	claim, claimsOk := token.Claims.(*stdjwt.StandardClaims)
	if !claimsOk {
		return nil, fmt.Errorf("erro a ler jwt")
	}

	if claim.Subject == "" {
		return nil, fmt.Errorf("atributo sub obrigatótio para identificar utilizador")
	}

	//Obter segredo da chave JWT
	jwtSecret := os.Getenv("jwtSecret")

	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	return []byte(jwtSecret), nil
}

// Existem duas formas de dar erro nos endpoints
// Por exemplo, no Login,
// return postJDTResponse{}, errors.New("parâmetros em falta"), usar em erros de transporte
// Ou
// return postJDTResponse{err:errors.New("parâmetros em falta")},nil, //usar em erros na camada de negocios

// O jwt.NewParser devolve um endpoint.Middleware, que é uma função que dado um Endpoint recebe outro Endpoint.
// Assim sendo, o parser recebe o Endpoint e devolve outro, que vai primeiro executar a verificação do JWT, e depois executa o nosso Endpoint
var jwtParser = jwt.NewParser(keyFunc, stdjwt.SigningMethodHS256, jwt.StandardClaimsFactory)

// MakeGetUsersEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeIntegrarDiario(s Service) endpoint.Endpoint {
	return jwtParser(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(IntegrarDiarioRequest)

		//Nota: Primeiro passa pelo middleware, depois pelo service!
		e := s.IntegrarDiario(ctx, req.Id, req.Origem)
		if e != nil {
			return IntegrarDiarioResponse{Success: false, Err: e}, nil
		} else {
			return IntegrarDiarioResponse{Success: true, Err: nil}, nil
		}
	})
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

//  *Body contains the expected body as JSON
//	*Request contains body+query parameters ready for service
//	*Response contains JSON response

type IntegrarDiarioRequest struct {
	Id     int
	Origem string
}

type IntegrarDiarioResponse struct {
	Success bool  `json:"success"`
	Err     error `json:"err,omitempty"`
}

func (r IntegrarDiarioResponse) error() error { return r.Err }
