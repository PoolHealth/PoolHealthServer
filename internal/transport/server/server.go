package server

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vektah/gqlparser/v2/ast"

	generated "github.com/PoolHealth/PoolHealthServer/pkg/api/v1/graphql"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type Server struct {
	resolvers   generated.ResolverRoot
	router      *mux.Router
	middlewares []graphql.HandlerExtension
	wsInitFunc  transport.WebsocketInitFunc

	log log.Logger
}

type Middleware func(handler http.Handler) http.Handler

func (s *Server) Run(ctx context.Context) error {
	http.Handle("/", s.router)

	originsOk := handlers.AllowedOrigins([]string{"*"})

	s.log.Info("server start on port 8080")

	server := &http.Server{Addr: ":8080", Handler: handlers.CORS(originsOk)(s.router)}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			s.log.Error(err)
		}
	}()

	s.router.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	s.router.Handle("/metrics", promhttp.Handler())

	return server.ListenAndServe()
}

func (s *Server) InitV1Api() {
	srv := handler.New(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: s.resolvers}))

	srv.AddTransport(&transport.Websocket{
		InitFunc:              s.wsInitFunc,
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(apollotracing.Tracer{})
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	for _, m := range s.middlewares {
		srv.Use(m)
	}
	// srv.Use(extension.FixedComplexityLimit(100))

	s.router.Use()
	s.router.Handle("/v1/", playground.Handler("GraphQL playground", "/v1/query"))
	s.router.Handle("/v1/query", srv)
}

func NewServer(
	resolvers generated.ResolverRoot,
	middlewares []graphql.HandlerExtension,
	wsInitFunc transport.WebsocketInitFunc,
	logger log.Logger,
) *Server {
	return &Server{
		resolvers:   resolvers,
		router:      mux.NewRouter(),
		middlewares: middlewares,
		wsInitFunc:  wsInitFunc,
		log:         logger,
	}
}
