package inventory

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var db *sql.DB
var Repo Repository

// Setup inventory package. Don't forget to call this function before using this package.
func Setup(d *sql.DB, r Repository) {
	db = d
	Repo = r
}

// Product entity
type Product struct {
	SKU string `json:"sku"`
	// MAYBE We may split stock to multiple bucket to spread hot spot in more advance implementation.
	Stock int `json:"stock"`
	// Order to handle.
	// For this simple implementation we will handle only one Order.
	// MAYBE We may handle multiple Order in more advance implementation, ex: CQRS.
	Order *Order
	// [transient]
	// Saved flag whether this state have been persisted.
	Saved bool
}

// Apply command.
// MAYBE We may handle multiple Order Create in more advance implementation, ex: CQRS.
func (p *Product) Apply(oc *OrderCreate) error {
	// check stock
	if p.Stock < oc.Quantity {
		return errors.New("product is not available")
	}

	// decrease stock by order quantity
	p.Stock -= oc.Quantity

	// make new order
	p.Order = &Order{
		ID:       uuid.New().String(),
		SKU:      oc.SKU,
		Quantity: oc.Quantity,
		Status:   OrderStatusApproved,
	}

	return nil
}

// Order entity
type Order struct {
	ID       string `json:"id"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

// Order status consts for Order.Status
const (
	OrderStatusApproved = "approved"
	OrderStatusCanceled = "canceled"
)

// OrderCreate command
type OrderCreate struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

// Repository of inventory
type Repository interface {
	// FindProduct find a product by ID
	FindProduct(ID string) (*Product, error)
	// FindProductTx find a product by ID with transaction
	FindProductTx(tx *sql.Tx, ID string) (*Product, error)
	// UpdateProductTx update a product with transaction
	UpdateProductTx(tx *sql.Tx, p *Product) (*Product, error)
	// FindOrder find an order
	FindOrder(ID string) (*Order, error)
}

// ProcessOrderCreate process create order
func ProcessOrderCreate(oc *OrderCreate) (*Order, error) {
	var o *Order

	// retry until commit success
	var has_committed bool = false
	for ok := true; ok; ok = !has_committed {
		// begin transaction with repeatable read isolation level
		tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		// find product
		p, err := Repo.FindProductTx(tx, oc.SKU)
		if err != nil {
			return nil, err
		}

		// apply command
		err = p.Apply(oc)
		if err != nil {
			return nil, err
		}

		// save product
		p, err = Repo.UpdateProductTx(tx, p)
		if err != nil {
			if strings.Contains(err.Error(), "could not serialize access due to concurrent update") {
				// rejected because of concurrent access
				// let's retry
				continue
			}

			return nil, err
		}

		// commit
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
		has_committed = true

		// set order for output
		o = p.Order
	}

	return o, nil
}
