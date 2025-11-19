package httpserver

import (
	"gopher/internal/model"
	"gopher/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
)

type server struct {
	router  *mux.Router
	logger  *logrus.Logger
	storage storage.Storage
}

func newServer(storage storage.Storage) *server {
	s := &server{
		router:  mux.NewRouter(),
		logger:  logrus.New(),
		storage: storage,
	}

	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return nil
	}

	s.logger.Info("server has been launched")

	return s
}

func (s *server) configureRouter() {
	s.router.Path("/users").Methods("POST").HandlerFunc(s.handleCreate)
	s.router.Path("/user").Methods("GET").HandlerFunc(s.handleFindUser)
}

func (s *server) handleFindUser(w http.ResponseWriter, r *http.Request) {
	email := "kelmann@mail.ru"

	user, _ := s.storage.User().FindByEmail(email)

	pp.Println(user)
}

func (s *server) handleCreate(w http.ResponseWriter, r *http.Request) {
	user := &model.User{Email: "kovrikovich@mail.ru", Password: "09fdhdjk2"}

	s.storage.User().Create(user)
}

func (s *server) configureLogger() error {
	level, err := logrus.ParseLevel("debug")

	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}
