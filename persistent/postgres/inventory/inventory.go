package inventory

import (
	"database/sql"

	"github.com/isirfanm/online-store/inventory"
	_ "github.com/lib/pq"
)

type RepoImpl struct {
	db *sql.DB
}

// Setup inventory package. Don't forget to call this function before using this package.
func NewRepo(d *sql.DB) *RepoImpl {
	return &RepoImpl{db: d}
}

func (r *RepoImpl) FindProduct(ID string) (*inventory.Product, error) {
	p := &inventory.Product{}
	row := r.db.QueryRow("select sku, stock from product where sku = $1", ID)
	err := row.Scan(&p.SKU, &p.Stock)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *RepoImpl) FindProductTx(tx *sql.Tx, ID string) (*inventory.Product, error) {
	p := &inventory.Product{}
	row := tx.QueryRow("select sku, stock from product where sku = $1", ID)
	err := row.Scan(&p.SKU, &p.Stock)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *RepoImpl) UpdateProductTx(tx *sql.Tx, p *inventory.Product) (*inventory.Product, error) {
	// save Product
	_, err := tx.Exec(
		"update product set stock=$2 where sku=$1",
		p.SKU,
		p.Stock,
	)
	if err != nil {
		return nil, err
	}

	// save Order
	o, err := r.SaveOrderTx(tx, p.Order)
	if err != nil {
		return nil, err
	}
	p.Order = o
	p.Saved = true

	return p, nil
}

func (r *RepoImpl) FindOrder(ID string) (*inventory.Order, error) {
	o := &inventory.Order{}
	row := r.db.QueryRow("select id, sku, quantity, \"status\" from \"order\" where id = $1", ID)
	err := row.Scan(
		&o.ID,
		&o.SKU,
		&o.Quantity,
		&o.Status,
	)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (r *RepoImpl) SaveOrderTx(tx *sql.Tx, o *inventory.Order) (*inventory.Order, error) {
	_, err := tx.Exec(
		`insert into "order" (id, sku, quantity, "status") VALUES ($1, $2, $3, $4)`,
		o.ID,
		o.SKU,
		o.Quantity,
		o.Status,
	)
	if err != nil {
		return nil, err
	}

	return o, nil
}
