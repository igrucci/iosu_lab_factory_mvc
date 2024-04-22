package archive

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
	"iosu_ceh/model"
)

func GetArchiveOrders(ctx *fasthttp.RequestCtx, connPool *pgxpool.Pool) {
	orders := []model.ArchiveOrdering{}

	rows, err := connPool.Query(ctx, "SELECT orderingarchive.id, ordering_id, detail_name, count_detail, customer_name, date_registration, date_adding, status from orderingarchive join detail on detail_id = detail.id join customer on customer_id = customer.id ")

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	orders, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.ArchiveOrdering])
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	render.RenderHtml(ctx, "templates/manage/archive/orders.html", orders)
	ctx.SetStatusCode(fasthttp.StatusOK)

}
