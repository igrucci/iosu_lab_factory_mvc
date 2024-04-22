package castingForging

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
)

func GetAllTypes(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	type Type struct {
		Id       int    `db:"id"`
		TypeName string `db:"type_name"`
	}
	types := []Type{}
	rows, err := pool.Query(ctx, "select id, type_name from ForgingType UNION select id, type_name from CastingType")
	if err != nil {
		logrus.Error(err)
		return
	}
	types, err = pgx.CollectRows(rows, pgx.RowToStructByName[Type])

	if err != nil {
		logrus.Error(err)
	}
	render.RenderHtml(ctx, "templates/manage/forgingcasting.html", types)
}
