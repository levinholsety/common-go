package web

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

// Module represents a module.
type Module interface {
	Action(name string) Action
}

// Action represents an action of service.
type Action interface {
	setLogger(*Logger)
	setResponseWriter(http.ResponseWriter)
	setRequest(*http.Request)
	Method(name string) func()
}

var mdlMap = map[string]Module{}

// Register registers a service.
func Register(name string, mdl Module) {
	name = path.Join("/api", name) + "/"
	mdlMap[name] = mdl
}

func handleDir(relativePath string, mux *http.ServeMux, logger *Logger) {
	dirPath := path.Join("contents", relativePath)
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	pattern := path.Clean("/" + relativePath)
	if len(pattern) > 1 {
		pattern += "/"
	}
	absPath, _ := filepath.Abs(dirPath)
	logger.Logi(1, "- {pattern: %s, path: '%s'}", pattern, absPath)
	mux.Handle(pattern, http.FileServer(http.Dir(dirPath)))
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			handleDir(path.Join(relativePath, fileInfo.Name()), mux, logger)
		}
	}
}

func handleFunc(mux *http.ServeMux, mdlName string, mdl Module) {
	mux.HandleFunc(mdlName, func(w http.ResponseWriter, r *http.Request) {
		logger := NewLogger()
		logger.Log("---")
		logger.Log("request:")
		logger.Logi(1, "remoteAddr: %s", r.RemoteAddr)
		logger.Logi(1, "url: %s", r.URL)
		array := strings.SplitN(r.URL.Path[len(mdlName):], "/", 2)
		if len(array) != 2 {
			writeNotFound(w, logger)
			return
		}
		actName, methodName := array[0], array[1]
		act := mdl.Action(actName)
		if act == nil {
			writeNotFound(w, logger)
			return
		}
		act.setLogger(logger)
		act.setResponseWriter(w)
		act.setRequest(r)
		method := act.Method(methodName)
		if method == nil {
			writeNotFound(w, logger)
			return
		}
		method()
	})
}

func writeError(err error, code int, w http.ResponseWriter, logger *Logger) {
	logger.Logi(1, `error: "%d - %s"`, code, err)
	http.Error(w, err.Error(), code)
}

func writeNotFound(w http.ResponseWriter, logger *Logger) {
	writeError(errors.New("page not found"), http.StatusNotFound, w, logger)
}

// Listen listens on the address for handling requests.
func Listen(addr string) error {
	logger := NewLogger()
	mux := http.NewServeMux()
	logger.Log("---")
	logger.Log("contents:")
	handleDir(".", mux, logger)
	logger.Log("services:")
	for mdlName, mdl := range mdlMap {
		logger.Logi(1, "- {pattern: %s, service: %s}", mdlName, reflect.TypeOf(mdl).Elem().String())
		handleFunc(mux, mdlName, mdl)
	}
	return http.ListenAndServe(addr, mux)
}
