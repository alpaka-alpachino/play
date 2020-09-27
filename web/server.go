package web

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/play/config"
	"github.com/play/notifications/bot"
)

var NotificationBot *bot.NotificationBot

type Perceptron struct {
	input        [][]float64
	actualOutput []float64
	weights      []float64
	bias         float64
	epochs       int
}

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	addr := *flag.String("address", ":1433", "address of server")
	flag.Parse()

	router := newRouter()

	return &Server{
		&http.Server{
			Addr:    addr,
			Handler: ddosFilter(router),
		},
	}
}

func (ws *Server) Run() error {
	return ws.server.ListenAndServe()
}

func init() {
	configPath := flag.String("cp", "../config", "Path to config file")
	configName := flag.String("cn", "viper.config", "Name of config file")

	viper, err := config.NewConfig(*configName, *configPath)
	if err != nil {
		log.Fatal(err)
	}
	botConfig := viper.NewNotificationBotConfig()
	NotificationBot, err = bot.NewNotificationBot(botConfig)
	if err != nil {
		log.Println(err)
		return
	}
	go NotificationBot.Bot()
}

func ddosFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		header := r.Method
		requestURL := r.URL.String()
		contLenth := r.ContentLength

		h, u, l := 0.0, 0.0, 0.0
		if header == http.MethodGet {
			h = 1.0
		}
		if requestURL == "/getname?id=5" {
			u = 1.0
		}
		if contLenth == -1 {
			l = 1.0
		}

		goPerceptron := Perceptron{
			input: [][]float64{{0, 0, 1},
				{1, 1, 1},
				{1, 0, 1},
				{0, 1, 0}},
			actualOutput: []float64{0, 1, 1, 0},
			epochs:       1000,
		}
		goPerceptron.initialize()
		goPerceptron.Train()

		test := goPerceptron.ForwardPass([]float64{h, u, l})

		notif := "suspicious activity"
		if test > 0.9 {
			NotificationBot.SendNotification(410287439, notif)
		}
		notif = "all is ok"
		if test < 0.9 {
			NotificationBot.SendNotification(410287439, notif)
		}
	})
}

func dotProduct(v1, v2 []float64) float64 {
	dot := 0.0
	for i := 0; i < len(v1); i++ {
		dot += v1[i] * v2[i]
	}
	return dot
}

func vecAdd(v1, v2 []float64) []float64 {
	add := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		add[i] = v1[i] + v2[i]
	}
	return add
}

func scalarMatMul(s float64, mat []float64) []float64 {
	result := make([]float64, len(mat))
	for i := 0; i < len(mat); i++ {
		result[i] += s * mat[i]
	}
	return result
}

func (a *Perceptron) initialize() {
	rand.Seed(time.Now().UnixNano())
	a.bias = 0.0
	a.weights = make([]float64, len(a.input[0]))
	for i := 0; i < len(a.input[0]); i++ {
		a.weights[i] = rand.Float64()
	}
}

func (a *Perceptron) sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func (a *Perceptron) ForwardPass(x []float64) (sum float64) {
	return a.sigmoid(dotProduct(a.weights, x) + a.bias)
}

func (a *Perceptron) gradW(x []float64, y float64) []float64 {
	pred := a.ForwardPass(x)
	return scalarMatMul(-(pred-y)*pred*(1-pred), x)
}

func (a *Perceptron) gradB(x []float64, y float64) float64 {
	pred := a.ForwardPass(x)
	return -(pred - y) * pred * (1 - pred)
}

func (a *Perceptron) Train() {
	for i := 0; i < a.epochs; i++ {
		dw := make([]float64, len(a.input[0]))
		db := 0.0
		for length, val := range a.input {
			dw = vecAdd(dw, a.gradW(val, a.actualOutput[length]))
			db += a.gradB(val, a.actualOutput[length])
		}
		dw = scalarMatMul(2/float64(len(a.actualOutput)), dw)
		a.weights = vecAdd(a.weights, dw)
		a.bias += db * 2 / float64(len(a.actualOutput))
	}
}
