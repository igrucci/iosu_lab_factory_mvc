package factory

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"iosu_ceh/internal/render"
	"iosu_ceh/model"
	"math"
	"math/rand"
	"os"
	"strconv"
)

func TakeOrder(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {

	customerName := string(ctx.QueryArgs().Peek("customerName"))
	detailName := string(ctx.QueryArgs().Peek("detailName"))
	countStr := string(ctx.QueryArgs().Peek("count"))
	count, err := strconv.Atoi(countStr)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)

		logrus.Error("Failed to convert count to int: %v", err)
		return
	}

	reqData := model.RequestData{
		CustomerName: customerName,
		DetailName:   detailName,
		Count:        count,
	}

	tx, err := pool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		logrus.Error("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(context.Background())
	//count, err := strconv.Atoi(reqData.Count)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)

		logrus.Error("Failed to convert count to int: %v", err)
		return
	}

	_, _ = pool.Exec(context.Background(), "insert into Customer(customer_name) values ($1)", reqData.CustomerName)

	if _, err := pool.Exec(context.Background(), "insert into Ordering (detail_id, count_detail, customer_id) SELECT  (SELECT id FROM Detail WHERE detail_name = $1), $2, (SELECT id FROM Customer WHERE customer_name = $3 ORDER BY customer.id ASC  LIMIT 1);", reqData.DetailName, reqData.Count, reqData.CustomerName); err != nil {

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		logrus.Error("Failed to insert into Ordering: %v", err)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		logrus.Error("Failed to commit transaction: %v", err)

		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func GetOrders(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {

	var rows pgx.Rows
	var err error
	if Param := ctx.QueryArgs().Peek("param"); len(Param) > 0 {
		orderId, err := strconv.Atoi(string(Param))
		if err != nil {
			logrus.Error("Failed to convert param to int: %v", err)
		}
		rows, err = pool.Query(ctx, "select ordering.id, detail_name, count_detail, customer_name, date_registration, status from ordering join detail on ordering.detail_id = detail.id\njoin customer on ordering.customer_id = customer.id where ordering.id = $1 and status = 'available';", orderId)
		if err != nil {
			logrus.Error("Query failed: %v\n", err)
		}

	} else {
		rows, err = pool.Query(ctx, "select ordering.id, detail_name, count_detail,customer_name, date_registration, status \nfrom ordering join detail on ordering.detail_id = detail.id\njoin customer on ordering.customer_id = customer.id where status = 'available';")
		if err != nil {
			logrus.Error("Query failed: %v\n", err)
		}
	}
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])

	if orders == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)

		logrus.Error("No orders found")
		return
	}

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("QueryRow failed: %v\n", err)
		return
	}

	render.RenderHtml(ctx, "templates/factory/orders.html", orders)
	ctx.SetStatusCode(fasthttp.StatusOK)
	return

}

func CheckResources(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	type OrderData struct {
		DetailName  string `db:"detail_name"`
		DetailCount int    `db:"count_detail"`
	}
	//передаваемое айди
	orderId := ctx.UserValue("id")

	//начало транзакции
	tx, err := pool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(context.Background())

	// достаем имя и количество в заказе
	rows, err := pool.Query(ctx, "select detail_name, count_detail \nfrom ordering join detail on ordering.detail_id = detail.id where ordering.id = $1;", orderId)
	if err != nil {
		logrus.Error("Query failed: %v\n", err)
	}
	orderData := OrderData{}
	orderData, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[OrderData])

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("QueryRow failed: %v\n", err)
		return
	}

	// выделяем материалы для всех деталей
	var countResources int
	if err := pool.QueryRow(ctx, "select count_material from Detail where detail_name = $1", orderData.DetailName).Scan(&countResources); err != nil {
		logrus.Error(err)
	}

	СountResourсesForOrder := orderData.DetailCount * countResources
	_, err = pool.Exec(ctx, "update material set material_count = material_count - $1 where material.id = (select material_id from detail where detail_name = $2)", СountResourсesForOrder, orderData.DetailName)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23514" {
			logrus.Println("Not enough resources")

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to update material: %v", err)
	}

	// проверка сотрудников для плавки
	res, err := pool.Exec(ctx, "update worker SET is_available = $1 WHERE worker.id = ( SELECT worker.id FROM worker WHERE place_id = 1 AND position_id = 1 AND is_available = 'available' ORDER BY worker.id ASC  LIMIT 1);", orderId)
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		logrus.Println("no available workers for casting")
		return
	}
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	// проверка сотрудников для ковки
	res, err = pool.Exec(ctx, "update worker SET is_available = $1 WHERE worker.id = ( SELECT worker.id FROM worker WHERE place_id = 2 AND position_id = 1 AND is_available = 'available' ORDER BY worker.id ASC  LIMIT 1);", orderId)
	rowsAffected = res.RowsAffected()
	if rowsAffected == 0 {
		logrus.Println("no available workers for forging")
		return
	}
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	// проверка оборудования для плавки
	res, err = pool.Exec(ctx, "update equipment set is_available = $1 where equipment.id = ( SELECT equipment.id FROM equipment WHERE place_id = 1 AND is_available = 'available' ORDER BY equipment.id ASC  LIMIT 1); ", orderId)
	rowsAffected = res.RowsAffected()
	if rowsAffected == 0 {
		logrus.Println("no available equipment for casting")
		return
	}
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	// проверка оборудования для ковки
	res, err = pool.Exec(ctx, "update equipment set is_available = $1 where equipment.id = ( SELECT equipment.id FROM equipment WHERE place_id = 2  AND is_available = 'available' ORDER BY equipment.id ASC  LIMIT 1); ", orderId)
	rowsAffected = res.RowsAffected()
	if rowsAffected == 0 {
		logrus.Println("no available equipment for forging")
		return
	}
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	// изменем статус заказа
	if _, err := pool.Exec(ctx, "update ordering set status = 'in-progress' where id = $1", orderId); err != nil {
		logrus.Error(err.Error())
	}

	// конец транзакции
	err = tx.Commit(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to commit transaction: %v", err)
		return
	}

	redirectURl := fmt.Sprintf("/factory/casting/%s", orderId)
	ctx.Redirect(redirectURl, fasthttp.StatusOK)

}

func GetCasting(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	type OrderDataForCasting struct {
		DetailName  string `db:"detail_name"`
		DetailCount int    `db:"count_detail"`
	}
	//передаваемое айди
	OrderId := ctx.UserValue("id")

	orderIdStr, ok := OrderId.(string)
	if !ok {
		logrus.Error("wrong id")
	}
	orderIdInt, err := strconv.Atoi(orderIdStr)
	if err != nil {
		fmt.Fprint(ctx, err)
	}

	tx, err := pool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to begin transaction: %v", err)
		return
	}
	// достаем имя и количество в заказе
	rows, err := pool.Query(ctx, "select detail_name, count_detail from ordering join detail on ordering.detail_id = detail.id where ordering.id = $1;", OrderId)
	if err != nil {
		logrus.Error("Query failed: %v\n", err)
	}

	orderData := OrderDataForCasting{}
	orderData, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[OrderDataForCasting])

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("QueryRow failed: %v\n", err)
		return
	}

	// выделяем материалы для всех деталей
	var countResources int
	if err := pool.QueryRow(ctx, "select count_material from Detail where detail_name = $1", orderData.DetailName).Scan(&countResources); err != nil {
		logrus.Error(err)
	}

	// количество затрачиваемых материалов
	CountResourcesForOrder := orderData.DetailCount * countResources

	//количество отливок
	countCasting := int(math.Round((rand.Float64()*0.1 + 0.8) * float64(CountResourcesForOrder)))
	typeCasting := rand.Intn(3) + 1

	// вставить в таблицу casting
	if _, err := pool.Exec(ctx, "insert into casting (casting_type_id, count_casting, material_id, count_material, place_id)\nvalues ( $1,\n        $2,\n        (select material_id from detail where detail_name = $3),\n        $4,\n        1);", typeCasting, countCasting, orderData.DetailName, CountResourcesForOrder); err != nil {
		logrus.Error(err)
	}
	var typeCastingName string
	err = pool.QueryRow(ctx, "select type_name from casting join castingtype on casting_type_id = castingtype.id where casting.id = (select id from casting order by id desc limit 1) ").Scan(&typeCastingName)
	if err != nil {
		logrus.Error(err)
	}
	var materialTypeName string
	err = pool.QueryRow(ctx, "select type_name from casting join material on material_id = material.id  join public.materialtype on material_type_id = materialtype.id where casting.id = (select id from casting order by id desc limit 1) ").Scan(&materialTypeName)
	if err != nil {
		logrus.Error(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to commit transaction: %v", err)
		return
	}

	casting := model.Casting{
		OrderId:        orderIdInt,
		CountCasting:   countCasting,
		TypeCasting:    typeCastingName,
		MaterialName:   materialTypeName,
		CountResources: CountResourcesForOrder,
	}

	render.RenderHtml(ctx, "templates/factory/casting.html", casting)
	ctx.SetStatusCode(fasthttp.StatusOK)

}

func GetForging(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {

	type DataCasting struct {
		CastingId    int `db:"id"`
		CastingCount int `db:"count_casting"`
	}

	//передаваемое айди
	OrderId := ctx.UserValue("id")

	OrderIdStr, ok := OrderId.(string)
	if !ok {
		logrus.Error("wrong id")
	}
	orderIdInt, err := strconv.Atoi(OrderIdStr)
	if err != nil {
		logrus.Error(err)
	}

	//начало транзакции
	tx, err := pool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(context.Background())

	// достаем имя и количество из плавки
	rows, err := pool.Query(ctx, "select id, count_casting from casting where id = (select id from casting order by id desc limit 1);")
	if err != nil {
		//	logrus.Error("No casting found")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		logrus.Error("Query failed: %v\n", err)

	}

	dataCasting := DataCasting{}
	dataCasting, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[DataCasting])

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("QueryRow failed: %v\n", err)
		return
	}

	firstTypeForging := rand.Intn(3) + 1

	secondTypeForging := rand.Intn(3) + 1

	firstCount := int(math.Round((0.5 * float64(dataCasting.CastingCount))))
	secondCount := dataCasting.CastingCount - firstCount

	if _, err := pool.Exec(ctx, "insert into forging(forging_type_id, casting_id, count_forging, place_id) values ($1, $2, $3, 2 )", firstTypeForging, dataCasting.CastingId, firstCount); err != nil {
		logrus.Error(err)
	}
	var firstId int

	if err := pool.QueryRow(ctx, "select id from forging where id = (select id from forging order by id desc limit 1 )").Scan(&firstId); err != nil {
		logrus.Error(err)
	}

	if _, err := pool.Exec(ctx, "insert into forging(forging_type_id, casting_id, count_forging, place_id) values ($1, $2, $3, 2 )", secondTypeForging, dataCasting.CastingId, secondCount); err != nil {

		logrus.Error(err)
	}
	var secondId int

	if err := pool.QueryRow(ctx, "select id from forging where id = (select id from forging order by id desc limit 1 )").Scan(&secondId); err != nil {
		logrus.Error(err)
	}
	var firstTypeForgingName string
	if err := pool.QueryRow(ctx, "select type_name from forging join forgingtype on forging_type_id = forgingType.id where forging.forging_type_id = $1 ", firstTypeForging).Scan(&firstTypeForgingName); err != nil {
		logrus.Error(err)
	}
	var secondTypeForgingName string
	if err := pool.QueryRow(ctx, "select type_name from forging join forgingtype on forging_type_id = forgingType.id where forging_type_id = $1 ", secondTypeForging).Scan(&secondTypeForgingName); err != nil {
		logrus.Error(err)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to commit transaction: %v", err)
		return
	}

	// Формируем данные о ковках для передачи в шаблон
	forgingData := map[string]interface{}{
		"OrderId":            orderIdInt,
		"FirstTypeForging":   firstTypeForgingName,
		"FirstCountForging":  firstCount,
		"SecondTypeForging":  secondTypeForgingName,
		"SecondCountForging": secondCount,
	}

	_, err = pool.Exec(ctx, "insert into TechnicalProcess(ordering_id, forging_id) values ($1, $2)", orderIdInt, firstId)
	if err != nil {
		logrus.Error(err)
	}
	_, err = pool.Exec(ctx, "insert into TechnicalProcess(ordering_id, forging_id) values ($1, $2)", orderIdInt, secondId)
	if err != nil {
		logrus.Error(err)
	}

	render.RenderHtml(ctx, "templates/factory/forging.html", forgingData)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func FinishOrder(ctx *fasthttp.RequestCtx, pool *pgxpool.Pool) {
	orderId := ctx.UserValue("id")

	tx, err := pool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(context.Background())

	rows, err := pool.Query(ctx, "select \n    o.id AS ordering_id,\n    d.detail_name,\n    o.count_detail,\n    c.customer_name,\n    f.count_forging,\n    ft.type_name AS forging_type,\n    ca.count_casting,\n    ct.type_name AS casting_type,\n    m.material_count,\n    o.date_registration,\n    tp.date_end AS date_end,\n    o.status\nFROM\n    Ordering o\n        JOIN\n    Detail d ON o.detail_id = d.id\n        JOIN\n    Customer c ON o.customer_id = c.id\n        LEFT JOIN\n    TechnicalProcess tp ON o.id = tp.ordering_id\n        LEFT JOIN\n    Forging f ON tp.forging_id = f.id\n        LEFT JOIN\n    ForgingType ft ON f.forging_type_id = ft.id\n        LEFT JOIN\n    Casting ca ON f.casting_id = ca.id\n        LEFT JOIN\n    CastingType ct ON ca.casting_type_id = ct.id\n        LEFT JOIN\n    Material m ON d.material_id = m.id where ordering_id = $1;", orderId)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	finData := model.FinishData{}
	finData, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[model.FinishData])
	if err != nil {
		fmt.Fprint(ctx, "not exist order")
		logrus.Error(err.Error())
		return
	}

	if _, err := pool.Exec(ctx, "update ordering set status = 'completed' where id = $1", orderId); err != nil {

		logrus.Error(err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err := pool.Exec(ctx, "update worker set is_available = 'available' where is_available = $1", orderId); err != nil {

		logrus.Error(err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err := pool.Exec(ctx, "update equipment set is_available = 'available' where is_available = $1", orderId); err != nil {
		logrus.Error(err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	err = tx.Commit(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("Failed to commit transaction: %v", err)
		return
	}
	PrintData(finData)
	render.RenderHtml(ctx, "templates/factory/finish.html", finData)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
func PrintData(finData model.FinishData) {

	if err := os.MkdirAll("storage", 0777); err != nil {
		logrus.Println("Error creating directory:", err)
		return
	}

	if _, err := os.Create("storage/order" + strconv.Itoa(finData.OrderingId)); err != nil {
		logrus.Println("Error creating file:", err)
		return
	}
	file, err := os.OpenFile("storage/order"+strconv.Itoa(finData.OrderingId), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	var data string
	data = "Договор на выполнение заказа изготовления деталей\nМежду:\nООО ИОСУ ЛАБ, именуемое в дальнейшем \"Исполнитель\", с одной стороны, и " + finData.CustomerName + ", именуемое в дальнейшем \"Заказчик\", с другой стороны, совместно именуемые \"Стороны\".\n\n1. Предмет договора:\nИсполнитель обязуется выполнить заказ на изготовление деталей, указанных в приложении к настоящему договору, а Заказчик обязуется принять и оплатить указанные детали согласно условиям настоящего договора.\n\nДетали заказа:\n\nOrdering ID: " + strconv.Itoa(finData.OrderingId) + "\n\nНаименование детали: " + finData.DetailName + "\n\nКоличество деталей: " + strconv.Itoa(finData.CountDetail) + "\n\nИмя заказчика: " + finData.CustomerName + "\n\nКоличество ковок: " + strconv.Itoa(finData.CountForging) + "\n\nТип ковки: " + finData.TypeForging + "\n\nКоличество литья: " + strconv.Itoa(finData.CountCasting) + "\n\nТип литья: " + finData.TypeCasting + "\n\nКоличество ресурсов: " + strconv.Itoa(finData.CountResources) + "\n\nДата регистрации: " + finData.DateRegistration.Time.Format("2006-01-02 15:04") + "\n\nДата завершения:" + finData.DateEnd.Time.Format("2006-01-02 15:04") + " \n\n2. Условия выполнения заказа:\n2.1 Исполнитель обязуется выполнить заказ в соответствии с требованиями и спецификациями, указанными в приложении к настоящему договору.\n\n2.2 В случае необходимости внесения изменений в заказ, стороны обязуются согласовать их письменно.\n\n3. Приемка и гарантия:\n3.1 Гарантийный срок на изготовленные детали составляет 24 месяца с момента их приемки."
	_, err = file.WriteString(data)
	if err != nil {
		logrus.Println("Error writing to file:", err)
		return
	}
}
