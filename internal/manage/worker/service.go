package worker

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
	"iosu_ceh/model"
)

func GetWorkers(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	workers := []model.Worker{}
	rows, err := pool.Query(ctx, "select worker.id, worker_name, position_name, place_name, is_available from Worker join Place on Worker.place_id = Place.id join position on worker.position_id = position.id")
	if err != nil {

		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	workers, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Worker])

	render.RenderHtml(ctx, "templates/manage/workers/workers.html", workers)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetWorker(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	workerId := ctx.UserValue("id").(string)
	worker := model.Worker{}
	rows, err := pool.Query(ctx, "select worker.id, worker_name, position_name, place_name, is_available from Worker join Place on Worker.place_id = Place.id join position on worker.position_id = position.id where worker.id = $1", workerId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	worker, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Worker])
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	render.RenderHtml(ctx, "templates/manage/workers/worker.html", worker)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func UpdateWorker(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {

	workerId := ctx.UserValue("id")
	workerName := string(ctx.QueryArgs().Peek("workerName"))
	position := string(ctx.QueryArgs().Peek("position"))
	place := string(ctx.QueryArgs().Peek("place"))
	isAvailable := string(ctx.QueryArgs().Peek("isAvailable"))

	worker := model.Worker{

		WorkerName:  workerName,
		Position:    position,
		Place:       place,
		IsAvailable: isAvailable,
	}

	tx, err := pool.Begin(context.Background())
	if err != nil {
		logrus.Error("Failed to begin transaction: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())
	_, _ = pool.Exec(ctx, "insert into Position(position_name) values ($1)", worker.Position)
	_, _ = pool.Exec(ctx, "insert into Place(place_name) values ($1)", worker.Place)
	_, err = pool.Exec(ctx, "UPDATE Worker SET worker_name = $1, position_id = (SELECT id FROM Position WHERE position_name = $2), place_id = (SELECT id FROM Place WHERE place_name = $3), is_available = $4 WHERE id = $5", worker.WorkerName, worker.Position, worker.Place, worker.IsAvailable, workerId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	err = tx.Commit(context.Background())
	if err != nil {
		logrus.Error("Failed to commit transaction ", err)
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func DeleteWorker(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	workerId := ctx.UserValue("id").(string)
	_, err := pool.Exec(ctx, "DELETE FROM Worker WHERE id = $1", workerId)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func CreateWorker(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	workerName := string(ctx.QueryArgs().Peek("workerName"))
	position := string(ctx.QueryArgs().Peek("position"))
	place := string(ctx.QueryArgs().Peek("place"))
	isAvailable := string(ctx.QueryArgs().Peek("isAvailable"))

	worker := model.Worker{

		WorkerName:  workerName,
		Position:    position,
		Place:       place,
		IsAvailable: isAvailable,
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to begin transaction ", err)
		return
	}
	defer tx.Rollback(context.Background())

	_, _ = pool.Exec(ctx, "insert into Position(position_name) values ($1)", worker.Position)
	_, _ = pool.Exec(ctx, "insert into Place(place_name) values ($1)", worker.Place)
	_, err = pool.Exec(ctx, "insert into Worker(worker_name, position_id, place_id, is_available) values ($1, (SELECT id FROM Position WHERE position_name = $2), (SELECT id FROM Place WHERE place_name = $3), $4)", worker.WorkerName, worker.Position, worker.Place, worker.IsAvailable)
	if err != nil {
		logrus.Error(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	err = tx.Commit(context.Background())
	if err != nil {
		logrus.Error("Failed to commit transaction ", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
