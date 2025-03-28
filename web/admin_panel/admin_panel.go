package admin_panel

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"vpn-tg-bot/internal/service/subscription"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/debug"
	"vpn-tg-bot/web/admin_panel/controller"
	"vpn-tg-bot/web/admin_panel/views"

	"github.com/gorilla/mux"
)

//go:embed public/*
var publicAssets embed.FS

// var publicAssetsMeta map[string]entity.FileMeta = make(map[string]entity.FileMeta)

const AssetsPrefix = "public"

type Server struct {
	Addr       string
	Scheme     string
	storage    storage.Storage
	sessionKey string
	certFile   string
	keyFile    string

	panel        *controller.PanelController
	user         *controller.UserController
	server       *controller.ServerController
	subscription *controller.SubscriptionController

	httpServer          http.Server
	subscriptionService subscription.VPN_API
	// ctx context.Context
}

type Settings struct {
	Addr       string
	Scheme     string
	Storage    storage.Storage
	SessionKey string

	CertFile string
	KeyFile  string

	SubscriptionService subscription.VPN_API
}

func New(s Settings) *Server {
	return &Server{
		Addr:                s.Addr,
		Scheme:              s.Scheme,
		storage:             s.Storage,
		sessionKey:          s.SessionKey,
		certFile:            s.CertFile,
		keyFile:             s.KeyFile,
		subscriptionService: s.SubscriptionService,
	}
}

func (s *Server) Run() (err error) {
	err = s.CheckSettings()
	if err != nil {
		return err
	}
	r, err := s.initRouter()
	if err != nil {
		return err
	}

	if s.Scheme != "https" {
		s.Scheme = "http"
	}

	path := fmt.Sprintf("%s://%s", s.Scheme, s.Addr)
	log.Printf("Admin panel is running on %s", path)

	s.specifyViewPaths(path)
	if s.Scheme == "https" {
		return http.ListenAndServeTLS(s.Addr, s.certFile, s.keyFile, r)
	}
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
	s.user = controller.NewUserController(r, s.storage, s.subscriptionService)
	s.server = controller.NewServerController(r, s.storage)
	s.subscription = controller.NewSubscriptionController(r, s.storage)

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

func (s *Server) specifyViewPaths(basePath string) {
	views.BasePath = basePath
	views.PublicPath = basePath + "/" + AssetsPrefix
	views.PublicDestPath = views.PublicPath + "/dest"

	// err := s.collectFileMetadata()
	// if err != nil {
	// 	log.Printf("[ERR] Can't collect file metadata: %v", err)
	// }
	// views.PublicAssetsMeta = publicAssetsMeta
}

func (s *Server) CheckSettings() error {
	if s.subscriptionService == nil {
		return fmt.Errorf("subscription service is not set")
	}
	if s.storage == nil {
		return fmt.Errorf("storage is not set")
	}
	if s.sessionKey == "" {
		return fmt.Errorf("session key is not set")
	}
	if s.Addr == "" {
		s.Addr = ":8080"
	}
	if s.Scheme == "" {
		s.Scheme = "http"
	}
	return nil
}

// func (s *Server) collectFileMetadata() (err error) {
// 	defer func() { e.WrapIfErr("can't collect file metadata", err) }()

// 	var paths []string

// 	err = fs.WalkDir(publicAssets, AssetsPrefix, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !d.IsDir() {
// 			paths = append(paths, path)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return e.Wrap("can't walk assets directories", err)
// 	}

// 	fmt.Printf("paths: %+v", paths)

// 	for _, path := range paths {
// 		info, err := os.Stat("/" + path)
// 		if err != nil {
// 			return e.Wrap("can't stat file", err)
// 		}
// 		path = strings.TrimPrefix(path, "/")
// 		path = strings.TrimPrefix(path, "/"+AssetsPrefix+"/")

// 		publicAssetsMeta[path] = entity.FileMeta{
// 			Path:    path,
// 			ModTime: info.ModTime(),
// 		}
// 	}

// 	return err
// }
