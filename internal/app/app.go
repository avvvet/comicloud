package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

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

// - picture url
// - title / description
// - web url for browser view
// - publishing date

// Sort the resulting feed by publishing date from recent to older.

type ComicData struct {
	PicUrl      string
	WebUrl      string
	Title       string
	PublishDate int64
}

/*
  structure returned
  from xkcd.com
*/
type XkcdData struct {
	Month      string
	Num        int
	Link       string
	Year       string
	News       string
	Safe_title string
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

var ComicStore []*ComicData

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

func (h *HttpServer) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		/*
			call external api
			example URL
			apiURL := "https://xkcd.com/614/info.0.json"
		*/

		limit := 10
		var lastUpdate int = 2677 //this will be dynamic

		for i := 0; i < limit; i++ {
			resBody, err := fetchRecentComic("https://xkcd.com/"+strconv.Itoa(lastUpdate-i)+"/info.0.json", limit)
			if err != nil {
				logger.Sugar().Warn(err)
				continue
			}

			//json interface for the api
			xkcd := &XkcdData{}
			err = json.Unmarshal(resBody, xkcd)
			if err != nil {
				logger.Sugar().Warn(err)
				continue
			}

			//convert to common structure
			c := &ComicData{
				PicUrl: xkcd.Img,
				WebUrl: xkcd.Link,
				Title:  xkcd.Title,
			}

			// parse date
			y, _ := strconv.Atoi(xkcd.Year)
			m, _ := strconv.Atoi(xkcd.Month)
			d, _ := strconv.Atoi(xkcd.Day)
			c.PublishDate = parseXkcdDate(y, m, d)

			//store
			ComicStore = append(ComicStore, c)
			//fmt.Printf("client: response body: %++v %s\n", xkcd, time.Unix(c.PublishDate, 0).UTC())
		}

		//response
		data, _ := json.Marshal(struct {
			ComicStore []*ComicData
			Length     int
		}{
			ComicStore: ComicStore,
			Length:     len(ComicStore),
		})

		io.WriteString(w, string(data[:]))
	}
}

func parseXkcdDate(y int, m int, d int) int64 {
	publishDate := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return publishDate.Unix()
}

func fetchRecentComic(apiURL string, limit int) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	return resBody, nil
}

func (h *HttpServer) Run() {
	http.HandleFunc("/", h.Index)
	http.HandleFunc("/status", h.Status)

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

func Art() {
	fmt.Println()
	fmt.Println("comic cloud ? yup you heard it right")
	myFigure := figure.NewColorFigure("comicloud", "", "gray", true)
	myFigure.Print()
	fmt.Println()
}
