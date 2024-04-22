package customer

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
	"iosu_ceh/model"
)

func GetCustomers(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	customers := []model.Customer{}
	rows, err := pool.Query(ctx, "SELECT customer.id , customer_name FROM Customer")
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	customers, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Customer])

	render.RenderHtml(ctx, "templates/manage/customers/customers.html", customers)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetCustomer(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	customerId := ctx.UserValue("id").(string)
	customer := model.Customer{}
	rows, err := pool.Query(ctx, "SELECT * FROM Customer WHERE id = $1", customerId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	customer, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Customer])
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	render.RenderHtml(ctx, "templates/manage/customers/customer.html", customer)
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func UpdateCustomer(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	customerId := ctx.UserValue("id").(string)
	customerName := string(ctx.QueryArgs().Peek("customerName"))
	customer := model.Customer{
		CustomerName: customerName,
	}

	_, err := pool.Exec(ctx, "UPDATE Customer SET customer_name = $1 WHERE id = $2", customer.CustomerName, customerId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func DeleteCustomer(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	customerId := ctx.UserValue("id").(string)
	_, err := pool.Exec(ctx, "DELETE FROM Customer WHERE id = $1", customerId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func CreateCustomer(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	customerName := string(ctx.QueryArgs().Peek("customerName"))
	customer := model.Customer{
		CustomerName: customerName,
	}

	_, err := pool.Exec(ctx, "INSERT INTO Customer(customer_name) VALUES ($1)", customer.CustomerName)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
