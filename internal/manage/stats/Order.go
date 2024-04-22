package stats

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
	"iosu_ceh/model"
)

func GetDynamic(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {

	rows, err := pool.Query(context.Background(), "SELECT\n    detail_name,\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 1 THEN count_detail ELSE 0 END) AS \"January\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 2 THEN count_detail ELSE 0 END) AS \"February\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 3 THEN count_detail ELSE 0 END) AS \"March\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 4 THEN count_detail ELSE 0 END) AS \"April\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 5 THEN count_detail ELSE 0 END) AS \"May\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 6 THEN count_detail ELSE 0 END) AS \"June\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 7 THEN count_detail ELSE 0 END) AS \"July\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 8 THEN count_detail ELSE 0 END) AS \"August\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 9 THEN count_detail ELSE 0 END) AS \"September\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 10 THEN count_detail ELSE 0 END) AS \"October\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 11 THEN count_detail ELSE 0 END) AS \"November\",\n    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 12 THEN count_detail ELSE 0 END) AS \"December\"\nFROM\n    detail left join ordering on detail.id = ordering.detail_id\nGROUP BY\n    detail_name;;")
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()

	var detailsStats []model.DetailStats
	detailsStats, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.DetailStats])

	render.RenderHtml(ctx, "templates/manage/dynamic.html", detailsStats)
}
