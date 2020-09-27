package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/play/store"
)

type DBmetric struct {
	db *sql.DB
}

const (
	sourceFormat = "host=localhost port=5432 user=%v password=%v dbname=%v sslmode=disable"
	serviceName  = "ServiceName"
)

func NewDBmetric() *DBmetric {

	s := fmt.Sprintf(sourceFormat, os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB"))

	database, err := sql.Open("postgres", s)

	if err != nil {
		database.Close()
		log.Fatal(err)
	}

	err = database.Ping()

	if err != nil {
		log.Println(err)
	}

	return &DBmetric{
		db: database,
	}
}

func (database *DBmetric) MetricsCreate(w http.ResponseWriter, r *http.Request) {

	var metrics store.Metrics
	var Status int

	err := json.NewDecoder(r.Body).Decode(&metrics)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		database.BadRequestHandler(w, r)
		return
	}

	vars := mux.Vars(r)
	metrics.ServiceName = vars[serviceName]

	time := time.Now()

	Status = http.StatusOK
	sqlStatement := `
	INSERT INTO metrictab(servicename,metricvalue,metricname,time,status)
	VALUES($1,$2,$3,$4,$5) `

	_, err = database.db.Exec(sqlStatement, metrics.ServiceName, metrics.MetricValue, metrics.MetricName, time, Status)

	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (database *DBmetric) BadRequestHandler(w http.ResponseWriter, r *http.Request) {

	var Status int

	vars := mux.Vars(r)
	ServiceName := vars[serviceName]

	time := time.Now()

	Status = http.StatusBadRequest

	sqlStatement := `
	INSERT INTO metrictab(servicename,metricvalue,metricname,time,status)
	VALUES($1,$2,$3,$4,$5) `

	_, err := database.db.Exec(sqlStatement, ServiceName, 0, 0, time, Status)

	if err != nil {
		fmt.Print("b")
		log.Println(err)
		return
	}
}

func (database *DBmetric) GetMetricsForService(w http.ResponseWriter, r *http.Request) {

	var metrics store.Metrics

	vars := mux.Vars(r)
	metrics.ServiceName = vars[serviceName]

	postgresqlQuerry := "SELECT servicename,metricvalue, metricname,time,status FROM metrictab WHERE servicename=$1"

	rows, err := database.db.Query(postgresqlQuerry, metrics.ServiceName)
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	ms := make([]*store.Metrics, 0)

	for rows.Next() {
		m := new(store.Metrics)
		err := rows.Scan(&m.ServiceName, &m.MetricValue, &m.MetricName, &m.Time, &m.Status)
		if err != nil {
			log.Fatal(err)
		}
		ms = append(ms, m)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, m := range ms {
		fmt.Fprintf(w, "%s %g %s %d %s \n", m.ServiceName, m.MetricValue, m.MetricName, m.Status, m.Time)
	}

	if err != nil {
		log.Println(err)
	}
}

func (database *DBmetric) GetSuccessNumberFromAll(w http.ResponseWriter, r *http.Request) {

	var metrics store.Metrics

	metrics.Status = 200

	postgresqlQuerry := "SELECT COUNT(*) FROM metrictab WHERE status=$1"

	rows, err := database.db.Query(postgresqlQuerry, metrics.Status)

	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Fprintf(w, "Number of successfully handled requests from all nodes is %d\n", count)
}

func (database *DBmetric) GetSuccessAndFailedForOne(w http.ResponseWriter, r *http.Request) {

	var metrics store.Metrics

	vars := mux.Vars(r)
	metrics.ServiceName = vars[serviceName]
	metrics.Status = 200

	postgresqlQuerry := "SELECT COUNT(status) FROM metrictab WHERE servicename=$1 AND status = $2 "
	srows, err := database.db.Query(postgresqlQuerry, metrics.ServiceName, metrics.Status)
	if err != nil {
		log.Println(err)
		return
	}

	defer srows.Close()

	var counts int

	for srows.Next() {
		if err := srows.Scan(&counts); err != nil {
			log.Fatal(err)
		}
	}

	metrics.Status = 400

	postgresqlQuerry = "SELECT COUNT(status) FROM metrictab WHERE servicename=$1 AND status = $2 "
	frows, err := database.db.Query(postgresqlQuerry, metrics.ServiceName, metrics.Status)
	if err != nil {
		log.Println(err)
		return
	}

	defer frows.Close()

	var countf int

	for frows.Next() {
		if err := frows.Scan(&countf); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Fprintf(w, "Number of successfully handled requests from all nodes is %d\n", counts)
	fmt.Fprintf(w, "Number of failed requests from all nodes is %d\n", countf)
}

func (database *DBmetric) HandledRequestsForDate(w http.ResponseWriter, r *http.Request) {

	var metrics store.Metrics

	vars := mux.Vars(r)

	str1 := vars["from"]
	mtime1, err := time.Parse(time.RFC3339, str1)

	metrics.Time = &mtime1
	f := metrics.Time

	str2 := vars["to"]
	mtime2, err := time.Parse(time.RFC3339, str2)

	metrics.Time = &mtime2
	t := metrics.Time

	dataRange := "SELECT *FROM metrictab WHERE time >= $1 AND time <= $2"
	data, err := database.db.Query(dataRange, f, t)

	if err != nil {
		log.Println(err)
		return
	}

	metrics.Status = 200

	postgresqlQuerry := "SELECT COUNT(status) FROM metrictab WHERE time=$1 AND status = $2 "
	srows, err := database.db.Query(postgresqlQuerry, data, metrics.Status)
	if err != nil {
		log.Println(err)
		return
	}

	defer srows.Close()

	var counts int

	for srows.Next() {
		if err := srows.Scan(&counts); err != nil {
			log.Fatal(err)
		}
	}

	metrics.Status = 400

	postgresqlQuerry = "SELECT COUNT(status) FROM metrictab WHERE time=$1 AND status = $2 "
	frows, err := database.db.Query(postgresqlQuerry, data, metrics.Status)
	if err != nil {
		log.Println(err)
		return
	}

	defer frows.Close()

	var countf int

	for frows.Next() {
		if err := frows.Scan(&countf); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Fprintf(w, "Number of successfully handled requests from all nodes is %d\n", counts)
	fmt.Fprintf(w, "Number of failed requests from all nodes is %d\n", countf)
}
