package repo

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repo struct {
	db *sqlx.DB
}

func New() *Repo {
	dbConfig := mysql.Config{
		User:      "root",
		Passwd:    "1234",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "resto-app",
		ParseTime: true,
	}

	db, err := sqlx.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic(err)
	}

	return &Repo{
		db: db,
	}
}

type RestaurantFilters struct {
	UUID *uuid.UUID
	Name *string
}
type Restaurant struct {
	UUID          *uuid.UUID `db:"uuid"`
	Name          string     `db:"name" json:"name" form:"name"`
	Description   string     `db:"description" json:"description" form:"description"`
	ContactPhone  string     `db:"contact_phone" json:"contactPhone" form:"contactPhone"`
	CoverImageURL string     `db:"cover_image_url" json:"coverImageURL" form:"coverImageUrl"`
	Address       string     `db:"address" json:"address" form:"address"`
}
type Table struct {
	UUID           *uuid.UUID `db:"uuid"`
	RestaurantUUID uuid.UUID  `db:"restaurant_uuid" json:"restaurantUuid" form:"restaurantUuid"`
	Number         int        `db:"number" json:"number" form:"restaurantUuid"`
}
type Reservation struct {
	UUID           *uuid.UUID `db:"uuid" json:"uuid" form:"uuid"`
	ClientPhone    string     `db:"client_phone" json:"clientPhone" form:"clientPhone"`
	StartDate      time.Time  `db:"start_date" json:"startDate" form:"startDate"`
	EndDate        time.Time  `db:"end_date" json:"endDate" form:"endDate"`
	RestaurantUUID uuid.UUID  `db:"restaurant_uuid" json:"restaurantUuid" form:"restaurantUuid"`
	TableUUID      uuid.UUID  `db:"table_uuid" json:"tableUuid" form:"tableUuid"`
}
type TableWithReservations struct {
	UUID           uuid.UUID     `json:"uuid"`
	RestaurantUUID uuid.UUID     `json:"restaurantUuid"`
	Number         int           `json:"number"`
	Reservations   []Reservation `json:"reservations"`
}
type RestaurantWithTables struct {
	UUID          uuid.UUID               `json:"uuid"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	ContactPhone  string                  `json:"contactPhone"`
	CoverImageURL string                  `json:"coverImageUrl"`
	Address       string                  `json:"address"`
	Tables        []TableWithReservations `json:"tables"`
}

func (r *Repo) GetOneRestaurant(uuid *uuid.UUID) (*RestaurantWithTables, error) {
	if uuid == nil {
		return nil, fmt.Errorf("repo error: uuid should be passed")
	}

	restaurant := Restaurant{}
	restaurantStatement := squirrel.Select("*").From("restaurants")
	restaurantStatement = restaurantStatement.Where(squirrel.Expr("uuid = uuid_to_bin(?)", uuid))
	restaurantSQL, args, err := restaurantStatement.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&restaurant, restaurantSQL, args...)
	if err != nil {
		return nil, err
	}

	tables := make([]Table, 0, 30)
	tablesStatement := squirrel.Select("*").From("restaurant_tables").Where(squirrel.Expr("restaurant_uuid = uuid_to_bin(?)", uuid))
	tablesSQL, args, err := tablesStatement.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&tables, tablesSQL, args...)
	if err != nil {
		return nil, err
	}

	reservations := make([]Reservation, 0, 30)
	reservationsStatement := squirrel.Select("*").From("table_reservations").Where(squirrel.Expr("restaurant_uuid = uuid_to_bin(?)", uuid))
	reservationsSQL, args, err := reservationsStatement.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&reservations, reservationsSQL, args...)
	if err != nil {
		return nil, err
	}

	fullTables := make([]TableWithReservations, 0, len(tables))
	for i, table := range tables {
		fullTables = append(fullTables, TableWithReservations{
			UUID:           *table.UUID,
			RestaurantUUID: table.RestaurantUUID,
			Number:         table.Number,
			Reservations:   make([]Reservation, 0, 30),
		})
		for _, reservation := range reservations {
			if *table.UUID == reservation.TableUUID {
				fullTables[i].Reservations = append(fullTables[i].Reservations, reservation)
			}
		}
	}
	fullRestaurant := RestaurantWithTables{
		UUID:          *restaurant.UUID,
		Name:          restaurant.Name,
		Description:   restaurant.Description,
		ContactPhone:  restaurant.ContactPhone,
		CoverImageURL: restaurant.CoverImageURL,
		Address:       restaurant.Address,
		Tables:        fullTables,
	}

	return &fullRestaurant, nil
}

func (r *Repo) CreateNewRestaurant(restaurant *Restaurant) (*uuid.UUID, error) {
	newUUID, _ := uuid.NewV7()

	sql, args, err := squirrel.Insert("restaurants").
		Columns("uuid", "name", "description", "contact_phone", "cover_image_url", "address").
		Values(squirrel.Expr("uuid_to_bin(?)", newUUID), restaurant.Name, restaurant.Description, restaurant.ContactPhone, restaurant.CoverImageURL, restaurant.Address).ToSql()
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	return &newUUID, nil
}

func (r *Repo) CreateNewTable(table *Table) (*uuid.UUID, error) {
	newUUID, _ := uuid.NewV7()

	sql, args, err := squirrel.Insert("restaurant_tables").
		Columns("uuid", "restaurant_uuid", "number").
		Values(squirrel.Expr("uuid_to_bin(?)", newUUID), squirrel.Expr("uuid_to_bin(?)", table.RestaurantUUID), table.Number).ToSql()
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	return &newUUID, nil
}

func (r *Repo) CreateNewReservation(reservation *Reservation) (*uuid.UUID, error) {
	newUUID, _ := uuid.NewV7()

	sql, args, err := squirrel.Insert("table_reservations").
		Columns("uuid", "restaurant_uuid", "table_uuid", "client_phone", "start_date", "end_date").
		Values(squirrel.Expr("uuid_to_bin(?)", newUUID), squirrel.Expr("uuid_to_bin(?)", reservation.RestaurantUUID), squirrel.Expr("uuid_to_bin(?)", reservation.TableUUID), reservation.ClientPhone, reservation.StartDate, reservation.EndDate).ToSql()
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	return &newUUID, nil
}
