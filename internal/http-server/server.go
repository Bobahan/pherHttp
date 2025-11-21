package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"gopher/internal/dto"
	"gopher/internal/model"
	"gopher/internal/storage"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "bobahan"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestId
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuth                  = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router  *mux.Router
	logger  *logrus.Logger
	storage storage.Storage
	session sessions.Store
}

func newServer(storage storage.Storage, session sessions.Store) *server {
	s := &server{
		router:  mux.NewRouter(),
		logger:  logrus.New(),
		storage: storage,
		session: session,
	}

	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return nil
	}

	s.logger.Info("server has been launched")

	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestIdMiddleware)
	s.router.Use(s.logRequestMiddleware)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.Path("/users").Methods("POST").HandlerFunc(s.handleCreateUser)
	s.router.Path("/users/{id}").Methods("GET").HandlerFunc(s.handleFindUser)
	s.router.Path("/sessions").Methods("POST").HandlerFunc(s.handleCreateSession)

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authMiddleware)
	private.HandleFunc("/whoami", s.handleWhoami).Methods("GET")
}

func (s *server) handleFindUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user, err := s.storage.User().FindById(id)

	if err != nil {
		return
	}

	user.Sanitize()

	s.response(w, r, http.StatusOK, user)
}

func (s *server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if err != nil {
		return
	}

	var userDto dto.User

	if err := json.Unmarshal(b, &userDto); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	user := &model.User{Email: userDto.Email, Password: userDto.Passowrd}

	if err := s.storage.User().Create(user); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	user.Sanitize()
	s.response(w, r, http.StatusOK, user)
}

func (s *server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if err != nil {
		return
	}

	var userDto dto.User

	if err := json.Unmarshal(b, &userDto); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := s.storage.User().FindByEmail(userDto.Email)

	if err != nil || !user.CheckPassword(userDto.Passowrd) {
		s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		return
	}

	session, err := s.session.Get(r, sessionName)

	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = user.Id

	if err := s.session.Save(r, w, session); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.response(w, r, http.StatusOK, nil)

}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.session.Get(r, sessionName)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]

		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuth)
			return
		}

		user, err := s.storage.User().FindById(strconv.Itoa(id.(int)))

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuth)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	})
}

func (s *server) setRequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()

		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestId, id)))
	})
}

func (s *server) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestId),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		startAt := time.Now()
		next.ServeHTTP(w, r)

		logger.Infof("completed in %v", time.Since(startAt))
	})
}

func (s *server) handleWhoami(w http.ResponseWriter, r *http.Request) {
	s.response(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
}

func (s *server) error(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error) {
	s.response(w, r, httpStatusCode, map[string]string{"error": err.Error()})
}

func (s *server) response(w http.ResponseWriter, r *http.Request, httpStatusCode int, data interface{}) {
	w.WriteHeader(httpStatusCode)

	if data != nil {
		b, err := json.Marshal(data)

		if err != nil {
			return
		}

		w.Write(b)
	}
}

func (s *server) configureLogger() error {
	level, err := logrus.ParseLevel("debug")

	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}
