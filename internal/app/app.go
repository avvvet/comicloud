package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/common-nighthawk/go-figure"
	"github.com/savioxavier/termlink"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
)

type HttpServer struct {
	port    uint
	context context.Context
}

func NewApp(ctx context.Context, port uint) *HttpServer {
	return &HttpServer{port, ctx}
}
func (h *HttpServer) Status(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		//lets create struct on th fly
		data, _ := json.Marshal(struct {
			Status string
		}{
			Status: "server started...",
		})
		io.WriteString(w, string(data[:]))
	}
}

func Art() {
	fmt.Println()
	fmt.Println("comic cloud ? yup you heard it right")
	myFigure := figure.NewColorFigure("comicloud", "", "gray", true)
	myFigure.Print()
	fmt.Println()
}

func (h *HttpServer) Run() {
	http.HandleFunc("/", h.Status)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(h.port)))
	if err != nil {
		panic(err)
	}
	logger.Sugar().Infof("listening web requests at port 😎️ %v", listener.Addr().(*net.TCPAddr).Port)
	fmt.Println("")
	fmt.Println(termlink.ColorLink("access your recent comic cloud at this link 👉️ ", "http://127.0.0.1:"+strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), "italic blue"))
	if err := http.Serve(listener, nil); err != nil {
		logger.Sugar().Fatal(err)
	}
}
