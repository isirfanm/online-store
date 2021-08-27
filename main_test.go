package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isirfanm/online-store/api"
	"github.com/isirfanm/online-store/config"
	"github.com/isirfanm/online-store/inventory"
	log "github.com/sirupsen/logrus"
)

var testRouter *gin.Engine

func init() {
	// Init DB
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:15432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to database. " + err.Error())
	}
	defer db.Close()

	db.Exec(`drop table "order"`)
	db.Exec(`drop table product`)

	_, err = db.Exec(`
	create table product (
		sku text not null,
		stock integer not null,
		constraint product_pk primary key (sku)
	)`)
	if err != nil {
		log.Fatal("cannot initialize table product. " + err.Error())
	}

	_, err = db.Exec(`
	create table "order" (
		id uuid not null,
		sku text not null references product(sku),
		quantity integer not null,
		"status" text not null,
		constraint order_pk primary key (id)
	)`)
	if err != nil {
		log.Fatal("cannot initialize table order. " + err.Error())
	}

	_, err = db.Exec(`insert into product (sku, stock) values ('aaaa', 1000)`)
	if err != nil {
		log.Fatal("cannot initialize product sku aaaa. " + err.Error())
	}

	_, err = db.Exec(`insert into product (sku, stock) values ('bbbb', 2000)`)
	if err != nil {
		log.Fatal("cannot initialize product sku bbbb. " + err.Error())
	}

	// setup
	config.SetupAll()

	// init router
	testRouter = api.SetupRouter()
}

func TestGetProduct(t *testing.T) {
	// make req aaaa
	req, err := http.NewRequest("GET", "/products/aaaa", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatal("failed to get product sku aaaa")
	}

	// make req aaaa
	req, err = http.NewRequest("GET", "/products/bbbb", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatal("failed to get product sku bbbb")
	}
}

func TestCreateOrder(t *testing.T) {
	// get product aaaa
	ap, err := inventory.Repo.FindProduct("aaaa")
	if err != nil {
		t.Fatal(err)
	}

	// make req
	aco := inventory.OrderCreate{
		SKU:      "aaaa",
		Quantity: 1,
	}
	aps, err := json.Marshal(aco)
	if err != nil {
		t.Fatal(err)
	}

	aReq, err := http.NewRequest("POST", "/orders", bytes.NewReader(aps))
	if err != nil {
		t.Error(err)
		return
	}

	aResp := httptest.NewRecorder()
	testRouter.ServeHTTP(aResp, aReq)

	if aResp.Code != 200 {
		t.Fatal("failed to create order sku aaaa")
	}

	// get product result
	apx, err := inventory.Repo.FindProduct("aaaa")
	if err != nil {
		t.Fatal(err)
	}

	if (ap.Stock - apx.Stock) != 1 {
		t.Fatal("inconsistent stock when create order sku aaaa")
	}

	// get product bbbb
	bp, err := inventory.Repo.FindProduct("bbbb")
	if err != nil {
		t.Fatal(err)
	}

	// make req
	bco := inventory.OrderCreate{
		SKU:      "bbbb",
		Quantity: 1,
	}
	bps, err := json.Marshal(bco)
	if err != nil {
		t.Fatal(err)
	}

	bReq, err := http.NewRequest("POST", "/orders", bytes.NewReader(bps))
	if err != nil {
		t.Error(err)
		return
	}

	bResp := httptest.NewRecorder()
	testRouter.ServeHTTP(aResp, bReq)

	if bResp.Code != 200 {
		t.Fatal("failed to create order sku bbbb")
	}

	// get product result
	bpx, err := inventory.Repo.FindProduct("bbbb")
	if err != nil {
		t.Fatal(err)
	}

	if (bp.Stock - bpx.Stock) != 1 {
		t.Fatal("inconsistent stock when create order sku bbbb")
	}
}

func TestConcurrentCreateOrder(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	// go for aaaa
	go func() {
		var awg sync.WaitGroup
		awg.Add(2000)

		for i := 0; i < 2000; i++ {
			go func() {
				co := inventory.OrderCreate{
					SKU:      "aaaa",
					Quantity: 1,
				}
				ps, _ := json.Marshal(co)
				req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(ps))
				resp := httptest.NewRecorder()
				testRouter.ServeHTTP(resp, req)

				awg.Done()
			}()
		}

		awg.Wait()
		wg.Done()
	}()

	// go for bbbb
	go func() {
		var awg sync.WaitGroup
		awg.Add(3000)

		for i := 0; i < 3000; i++ {
			go func() {
				co := inventory.OrderCreate{
					SKU:      "bbbb",
					Quantity: 1,
				}
				ps, _ := json.Marshal(co)
				req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(ps))
				resp := httptest.NewRecorder()
				testRouter.ServeHTTP(resp, req)

				awg.Done()
			}()
		}

		awg.Wait()
		wg.Done()
	}()

	wg.Wait()

	// get product aaaa
	ap, err := inventory.Repo.FindProduct("aaaa")
	if err != nil {
		t.Fatal(err)
	}

	if ap.Stock != 0 {
		t.Fatal("inconsistent stock when concurrent create order sku aaaa")
	}

	// get product bbbb
	bp, err := inventory.Repo.FindProduct("bbbb")
	if err != nil {
		t.Fatal(err)
	}

	if bp.Stock != 0 {
		t.Fatal("inconsistent stock when concurrent create order sku bbbb")
	}
}
