package model

import "github.com/jackc/pgx/v5/pgtype"

type Order struct {
	Id               int8               `json:"ordering_id" db:"id"`
	CustomerName     string             `json:"customer" db:"customer_name"`
	DetailName       string             `json:"detail_name" db:"detail_name"`
	CountDetail      int                `json:"count_detail" db:"count_detail"`
	DateRegistration pgtype.Timestamptz `json:"date_registration" db:"date_registration"`
	Status           string             `json:"status" db:"status"`
}

type Casting struct {
	OrderId        int    `json:"orderId"`
	CountCasting   int    `json:"countCasting"`
	TypeCasting    string `json:"typeCasting"`
	MaterialName   string `json:"materialName"`
	CountResources int    `json:"countResources"`
}

type Forging struct {
	ForgingId    int    `json:"forgingId"`
	OrderId      int    `json:"orderId"`
	TypeForging  string `json:"typeForging"`
	CountForging int    `json:"countForging"`
	CountCasting int    `json:"countCasting"`
}
type Detail struct {
	Id            int    `json:"id" db:"id"`
	DetailName    string `json:"detail_name" db:"detail_name"`
	MaterialName  string `json:"material_name" db:"type_name"`
	CountMaterial int    `json:"count_material" db:"count_material"`
}

type Customer struct {
	Id           int    `json:"id" db:"id"`
	CustomerName string `json:"customer_name" db:"customer_name"`
}

type Worker struct {
	Id          int    `json:"id" db:"id"`
	WorkerName  string `json:"worker_name" db:"worker_name"`
	Position    string `json:"position" db:"position_name"`
	Place       string `json:"place" db:"place_name"`
	IsAvailable string `json:"isAvailable" db:"is_available"`
}
type RequestData struct {
	CustomerName string `json:"customer_name"`
	DetailName   string `json:"detail_name"`
	Count        int    `json:"count"`
}
type ArchiveOrdering struct {
	Id               int                `db:"id"`
	OrderingId       int                `db:"ordering_id"`
	DetailName       string             `db:"detail_name"`
	CountDetail      int                `db:"count_detail"`
	CustomerName     string             `db:"customer_name"`
	DateRegistration pgtype.Timestamptz `db:"date_registration"`
	DateEnd          pgtype.Timestamptz `db:"date_adding"`
	Status           string             `db:"status"`
}

type DetailStats struct {
	DetailName string
	January    int
	February   int
	March      int
	April      int
	May        int
	June       int
	July       int
	August     int
	September  int
	October    int
	November   int
	December   int
}

type FinishData struct {
	OrderingId       int                `db:"ordering_id"`
	DetailName       string             `db:"detail_name"`
	CountDetail      int                `db:"count_detail"`
	CustomerName     string             `db:"customer_name"`
	CountForging     int                `db:"count_forging"`
	TypeForging      string             `db:"forging_type"`
	CountCasting     int                `db:"count_casting"`
	TypeCasting      string             `db:"casting_type"`
	CountResources   int                `db:"material_count"`
	DateRegistration pgtype.Timestamptz `db:"date_registration"`
	DateEnd          pgtype.Timestamptz `db:"date_end"`
	Status           string             `db:"status"`
}
