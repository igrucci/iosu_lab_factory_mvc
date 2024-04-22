package detail

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
	"iosu_ceh/model"
	"strconv"
)

func GetDetails(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	rows, err := pool.Query(ctx, "SELECT detail.id, detail_name, materialtype.type_name, count_material FROM detail\nJOIN material ON detail.material_id = material.id\nJOIN materialtype ON material.material_type_id = materialtype.id;")
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	details := []model.Detail{}
	details, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Detail])

	render.RenderHtml(ctx, "templates/manage/details/details.html", details)
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func GetDetail(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	detailId := ctx.UserValue("id").(string)
	rows, err := pool.Query(ctx, "SELECT detail.id, detail_name, materialtype.type_name, count_material FROM detail\nJOIN material ON detail.material_id = material.id\nJOIN materialtype ON material.material_type_id = materialtype.id where detail.id = $1", detailId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	detail := model.Detail{}
	detail, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Detail])

	render.RenderHtml(ctx, "templates/manage/details/detail.html", detail)
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func UpdateDetail(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	detailId := ctx.UserValue("id").(string)
	detailName := string(ctx.QueryArgs().Peek("detailName"))
	materialName := string(ctx.QueryArgs().Peek("materialName"))

	countMaterialStr := string(ctx.QueryArgs().Peek("countMaterial"))
	countMaterial, err := strconv.Atoi(countMaterialStr)

	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	detail := model.Detail{
		DetailName:    detailName,
		MaterialName:  materialName,
		CountMaterial: countMaterial,
	}

	_, _ = pool.Exec(ctx, "insert into MaterialType(type_name) values ($1)", detail.MaterialName)
	_, _ = pool.Exec(ctx, "insert into Material(material_type_id) values ((SELECT id FROM MaterialType WHERE type_name = $1))", detail.MaterialName)
	_, err = pool.Exec(ctx, "UPDATE Detail SET detail_name = $1, material_id = (SELECT material.id FROM Material join MaterialType on material_type_id = MaterialType.id WHERE type_name = $2), count_material = $3 WHERE id = $4",
		detail.DetailName, detail.MaterialName, detail.CountMaterial, detailId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func DeleteDetail(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	detailId := ctx.UserValue("id").(string)

	_, err := pool.Exec(ctx, "DELETE FROM Detail WHERE id = $1", detailId)

	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func CreateDetail(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {

	detailName := string(ctx.QueryArgs().Peek("detailName"))
	materialName := string(ctx.QueryArgs().Peek("materialName"))

	countMaterialStr := string(ctx.QueryArgs().Peek("countMaterial"))
	countMaterial, err := strconv.Atoi(countMaterialStr)

	detail := model.Detail{

		DetailName:    detailName,
		MaterialName:  materialName,
		CountMaterial: countMaterial,
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}
	defer tx.Rollback(context.Background())

	_, _ = pool.Exec(ctx, "insert into MaterialType(type_name) values ($1)", detail.MaterialName)
	_, _ = pool.Exec(ctx, "insert into Material(material_type_id) values ((SELECT id FROM MaterialType WHERE type_name = $1))", detail.MaterialName)
	_, err = pool.Exec(ctx, "insert into Detail(detail_name, material_id, count_material) values ($1, (SELECT material.id FROM Material join MaterialType on material_type_id = MaterialType.id WHERE type_name = $2), $3)",
		detail.DetailName, detail.MaterialName, detail.CountMaterial)
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
