package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gbrlsnchs/jwt/v2"
	gauthCfg "github.com/hiroaki-yamamoto/gauth/config"
	gauthMw "github.com/hiroaki-yamamoto/gauth/middleware"
	"github.com/hiroaki-yamamoto/todo-sample-backend/auth"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/todo"
	userRepo "github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/user"
	"github.com/hiroaki-yamamoto/todo-sample-backend/graph"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	config, err := LoadConfig("config")
	if err != nil {
		log.Fatal("failed to load config", err)
	}
	port := config.Port
	if port == "" {
		port = defaultPort
	}
	db, err := gorm.Open(postgres.Open(config.DB), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	todoRepo := todo.NewRepo(db)
	prv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(todoRepo)}),
	)

	prv.AddTransport(transport.Options{})
	prv.AddTransport(transport.GET{})
	prv.AddTransport(transport.POST{})

	prv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	prv.Use(extension.Introspection{})
	prv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	authCfg, err := gauthCfg.New(
		"jwt_token",
		gauthCfg.Cookie,
		jwt.NewHS512(config.Secret),
		"Todo-Sample-Backend-Audience",
		"Todo-Sample-Backend-Issuer",
		"Todo-Sample-Backend-Subject",
		24*time.Hour,
		gauthCfg.CookieConfig{
			SameSite: http.SameSiteLaxMode,
		},
	)
	if err != nil {
		log.Fatal("failed to create gauth config", err)
	}
	ur := userRepo.NewRepo(db)
	pub := handler.New(auth.NewExecutableSchema(auth.Config{
		Resolvers: &auth.Resolver{
			UserRepo:    ur,
			GAuthConfig: authCfg,
		},
	}))
	pub.AddTransport(transport.Options{})
	pub.AddTransport(transport.GET{})
	pub.AddTransport(transport.POST{})
	pub.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	pub.Use(extension.Introspection{})
	pub.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	findUserFunc := func(con interface{}, username string) (interface{}, error) {
		userRepoInstance, ok := con.(userRepo.IUserRepo)
		if !ok {
			return nil, errors.New("invalid context format")
		}
		return userRepoInstance.GetByID(context.Background(), username)
	}

	authMiddleware := gauthMw.ContextMiddleware(ur, findUserFunc, authCfg)
	loginRequired := gauthMw.LoginRequired(ur, findUserFunc, authCfg)

	mux := http.NewServeMux()
	mux.Handle("/pub_sandbox", playground.Handler("GraphQL playground", "/pub"))
	mux.Handle("/prv_sandbox", playground.Handler("GraphQL playground", "/prv"))
	mux.Handle("/prv", authMiddleware(loginRequired(prv)))
	mux.Handle("/pub", auth.InjectResponseWriter(authMiddleware(pub)))

	httpSrv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
