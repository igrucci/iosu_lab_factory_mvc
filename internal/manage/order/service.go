package order

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
	"iosu_ceh/model"
	"strconv"
	"time"
)

func GetOrders(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	var rows pgx.Rows
	var err error
	if Param := string(ctx.QueryArgs().Peek("param")); len(Param) > 0 {

		rows, err = pool.Query(ctx, "select ordering.id, detail_name, count_detail, customer_name, date_registration, status from ordering join detail on detail_id = detail.id join customer on customer_id = customer.id  where status = $1", Param)
		if err != nil {
			logrus.Error(err)
		}

	} else {
		rows, err = pool.Query(ctx, "select ordering.id, detail_name, count_detail, customer_name, date_registration, status from ordering join detail on detail_id = detail.id join customer on customer_id = customer.id  ")
		if err != nil {
			logrus.Error("QueryRow failed: %v\n", err)
		}
	}
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])

	if orders == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		logrus.Println("No orders found")
		return
	}

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("QueryRow failed: %v\n", err)
		return
	}

	render.RenderHtml(ctx, "templates/manage/orders/orders.html", orders)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}

func GetOrder(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	orderId := ctx.UserValue("id").(string)

	rows, err := pool.Query(ctx, "select ordering.id, detail_name, count_detail, customer_name, date_registration, status from ordering join detail on detail_id = detail.id join customer on customer_id = customer.id  WHERE ordering.id = $1", orderId)

	if err != nil {
		logrus.Error("QueryRow failed: %v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	order := model.Order{}
	order, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		logrus.Error("QueryRow failed: %v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	render.RenderHtml(ctx, "templates/manage/orders/order.html", order)
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func UpdateOrder(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	orderId := ctx.UserValue("id").(string)

	customerName := string(ctx.QueryArgs().Peek("customerName"))
	detailName := string(ctx.QueryArgs().Peek("detailName"))
	countDetailStr := string(ctx.QueryArgs().Peek("countDetail"))
	countDetail, err := strconv.Atoi(countDetailStr)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	dateRegistrationStr := string(ctx.QueryArgs().Peek("dateRegistration"))

	parsedTime, err := time.Parse("2006-01-02 15:04:05-07", dateRegistrationStr)
	if err != nil {
		logrus.Println("Error parsing time: ", err)
		return
	}

	fmt.Println(parsedTime)
	dateRegistration := pgtype.Timestamptz{
		Time:             parsedTime,
		InfinityModifier: pgtype.Infinity,
		Valid:            true,
	}

	status := string(ctx.QueryArgs().Peek("status"))

	order := model.Order{
		CustomerName:     customerName,
		DetailName:       detailName,
		CountDetail:      countDetail,
		DateRegistration: dateRegistration,
		Status:           status,
	}

	_, _ = pool.Exec(ctx, "insert into Customer(customer_name) values ($1)", order.CustomerName)

	_, err = pool.Exec(ctx, "UPDATE ordering SET detail_id = (SELECT id FROM Detail WHERE detail_name = $1), count_detail = $2, customer_id = (SELECT id FROM Customer WHERE customer_name = $3 ORDER BY customer.id ASC  LIMIT 1),date_registration = $4, status = $5 WHERE id = $6",
		order.DetailName, order.CountDetail, order.CustomerName, parsedTime, order.Status, orderId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

}

func DeleteOrder(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	orderId := ctx.UserValue("id").(string)

	type OrderDate struct {
		TimeStart pgtype.Timestamptz `db:"date_registration"`
		TimeEnd   pgtype.Timestamptz `db:"date_end"`
	}
	rows, err := pool.Query(ctx, "select date_registration, date_end from ordering join technicalprocess on TechnicalProcess.ordering_id = ordering.id where ordering.id = $1 limit 1", orderId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	orderDate := OrderDate{}
	orderDate, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[OrderDate])
	if err != nil {
		fmt.Fprint(ctx, err)
	}

	if _, err := pool.Exec(ctx, "insert into orderingarchive(ordering_id, detail_id, count_detail, customer_id, date_registration, status, date_adding) values ($1, (select detail_id from ordering where id = $1), (select count_detail from ordering where id = $1), (select customer_id from ordering where id = $1), $2, 'closed', $3)", orderId, orderDate.TimeStart, orderDate.TimeEnd); err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	_, err = pool.Exec(ctx, "delete from ordering WHERE id = $1 ", orderId)

	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func CreateOrder(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {

	customerName := string(ctx.QueryArgs().Peek("customerName"))
	detailName := string(ctx.QueryArgs().Peek("detailName"))
	countDetailStr := string(ctx.QueryArgs().Peek("countDetail"))
	countDetail, err := strconv.Atoi(countDetailStr)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	dateRegistrationStr := string(ctx.QueryArgs().Peek("dateRegistration"))

	parsedTime, err := time.Parse("2006-01-02 15:04:05-07", dateRegistrationStr)
	if err != nil {
		logrus.Println("Error parsing time: ", err)
		return
	}

	dateRegistration := pgtype.Timestamptz{
		Time:  parsedTime,
		Valid: true,
	}

	status := string(ctx.QueryArgs().Peek("status"))

	order := model.Order{
		DetailName:       detailName,
		CountDetail:      countDetail,
		CustomerName:     customerName,
		DateRegistration: dateRegistration,
		Status:           status,
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to begin transaction ", err)
		return
	}
	defer tx.Rollback(context.Background())
	_, _ = pool.Exec(ctx, "insert into Customer(customer_name) values ($1)", order.CustomerName)

	_, err = pool.Exec(ctx, "INSERT INTO ordering (detail_id, count_detail, customer_id, date_registration, status) VALUES ((SELECT id FROM Detail WHERE detail_name = $1), $2, (SELECT id FROM Customer WHERE customer_name = $3), $4, $5)",
		order.DetailName, order.CountDetail, order.CustomerName, order.DateRegistration, order.Status)

	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
