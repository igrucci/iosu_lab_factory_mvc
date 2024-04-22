package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"iosu_ceh/config"
	"iosu_ceh/internal/factory"
	"iosu_ceh/internal/manage/archive"
	"iosu_ceh/internal/manage/castingForging"
	"iosu_ceh/internal/manage/customer"
	"iosu_ceh/internal/manage/detail"
	"iosu_ceh/internal/manage/order"
	"iosu_ceh/internal/manage/stats"
	"iosu_ceh/internal/manage/worker"
	"iosu_ceh/internal/render"
	"log"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	connPool, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := fasthttprouter.New()

	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		render.RenderHtml(ctx, "templates/factory/home.html", nil)
	})

	router.POST("/factory/order", func(ctx *fasthttp.RequestCtx) {
		factory.TakeOrder(ctx, connPool)
	})
	router.GET("/factory/manage", func(ctx *fasthttp.RequestCtx) {
		factory.GetOrders(ctx, connPool)
	})
	router.GET("/factory/manage/:id", func(ctx *fasthttp.RequestCtx) {
		factory.CheckResources(ctx, connPool)
	})

	router.GET("/factory/casting/:id", func(ctx *fasthttp.RequestCtx) {
		factory.GetCasting(ctx, connPool)
	})
	router.GET("/factory/forging/:id", func(ctx *fasthttp.RequestCtx) {
		factory.GetForging(ctx, connPool)
	})

	router.GET("/factory/finish/:id", func(ctx *fasthttp.RequestCtx) {
		factory.FinishOrder(ctx, connPool)
	})

	// manage menu
	router.GET("/manage", func(ctx *fasthttp.RequestCtx) {
		render.RenderHtml(ctx, "templates/manage/menu.html", nil)
	})

	router.GET("/manage/order", func(ctx *fasthttp.RequestCtx) {
		order.GetOrders(ctx, connPool)
	})
	router.GET("/manage/order/:id", func(ctx *fasthttp.RequestCtx) {
		order.GetOrder(ctx, connPool)
	})

	router.PUT("/manage/order/:id", func(ctx *fasthttp.RequestCtx) {
		order.UpdateOrder(ctx, connPool)
	})

	router.DELETE("/manage/order/:id", func(ctx *fasthttp.RequestCtx) {
		order.DeleteOrder(ctx, connPool)
	})

	router.POST("/manage/order", func(ctx *fasthttp.RequestCtx) {
		order.CreateOrder(ctx, connPool)
	})

	router.GET("/manage/detail", func(ctx *fasthttp.RequestCtx) {
		detail.GetDetails(ctx, connPool)
	})

	router.GET("/manage/detail/:id", func(ctx *fasthttp.RequestCtx) {
		detail.GetDetail(ctx, connPool)
	})

	router.PUT("/manage/detail/:id", func(ctx *fasthttp.RequestCtx) {
		detail.UpdateDetail(ctx, connPool)
	})

	router.DELETE("/manage/detail/:id", func(ctx *fasthttp.RequestCtx) {
		detail.DeleteDetail(ctx, connPool)
	})

	router.POST("/manage/detail", func(ctx *fasthttp.RequestCtx) {
		detail.CreateDetail(ctx, connPool)
	})

	// customers

	router.GET("/manage/customer", func(ctx *fasthttp.RequestCtx) {
		customer.GetCustomers(ctx, connPool)
	})

	router.GET("/manage/customer/:id", func(ctx *fasthttp.RequestCtx) {
		customer.GetCustomer(ctx, connPool)
	})

	router.PUT("/manage/customer/:id", func(ctx *fasthttp.RequestCtx) {
		customer.UpdateCustomer(ctx, connPool)
	})

	router.DELETE("/manage/customer/:id", func(ctx *fasthttp.RequestCtx) {
		customer.DeleteCustomer(ctx, connPool)
	})

	router.POST("/manage/customer", func(ctx *fasthttp.RequestCtx) {
		customer.CreateCustomer(ctx, connPool)
	})

	//for workers
	router.GET("/manage/worker", func(ctx *fasthttp.RequestCtx) {

		worker.GetWorkers(ctx, connPool)
	})

	router.GET("/manage/worker/:id", func(ctx *fasthttp.RequestCtx) {
		worker.GetWorker(ctx, connPool)
	})

	router.PUT("/manage/worker/:id", func(ctx *fasthttp.RequestCtx) {
		worker.UpdateWorker(ctx, connPool)
	})

	router.DELETE("/manage/worker/:id", func(ctx *fasthttp.RequestCtx) {
		worker.DeleteWorker(ctx, connPool)
	})

	router.POST("/manage/worker", func(ctx *fasthttp.RequestCtx) {
		worker.CreateWorker(ctx, connPool)
	})
	router.GET("/manage/archive/orders", func(ctx *fasthttp.RequestCtx) {
		archive.GetArchiveOrders(ctx, connPool)
	})

	router.GET("/manage/dynamic", func(ctx *fasthttp.RequestCtx) {
		stats.GetDynamic(ctx, connPool)
	})

	router.GET("/manage/cftypes", func(ctx *fasthttp.RequestCtx) {
		castingForging.GetAllTypes(ctx, connPool)
	})
	if err := fasthttp.ListenAndServe(":8080", router.Handler); err != nil {
		log.Fatalf("%v", err)
	}
}
