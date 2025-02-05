package admin_panel

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/debug"
	"vpn-tg-bot/web/admin_panel/controller"

	"github.com/gorilla/mux"
)

//go:embed public/*
var publicAssets embed.FS

const AssetsPrefix = "public"

type Server struct {
	Addr       string
	storage    storage.Storage
	sessionKey string

	panel *controller.PanelController

	httpServer http.Server
	// ctx context.Context
}

type Settings struct {
	Addr       string
	Storage    storage.Storage
	SessionKey string
}

func New(s Settings) *Server {
	return &Server{
		Addr:       s.Addr,
		storage:    s.Storage,
		sessionKey: s.SessionKey,
	}
}

func (s *Server) Run() (err error) {
	r, err := s.initRouter()
	if err != nil {
		return err
	}
	log.Printf("Admin panel is running on %s", s.Addr)
	return http.ListenAndServe(s.Addr, r)
}

func (s *Server) Assets() fs.FS {
	subFS, err := fs.Sub(publicAssets, AssetsPrefix)
	if err != nil {
		panic("assets prefix must be the same as the embedded directory!")
	}
	return subFS
}

func (s *Server) initRouter() (http.Handler, error) {

	r := mux.NewRouter()

	s.panel = controller.NewPanelController(r, s.storage)

	// Debug
	d, err := os.Getwd()
	fmt.Println("current dir: ", d, err)
	err = debug.DisplayFileSystem(s.Assets())
	if err != nil {
		log.Fatal(err)
	}
	// Static assets router
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.FS(s.Assets()))))
	// assetsRouter.Handle("GET /public/", http.StripPrefix("/public", http.FileServer(http.FS(p.Assets()))))

	return r, nil
}
