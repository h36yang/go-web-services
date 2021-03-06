package product

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"demo/inventoryservice/database"
)

func getProduct(productID int) (*Product, error) {
	query := `SELECT
		productId,
		manufacturer,
		sku,
		upc,
		pricePerUnit,
		quantityOnHand,
		productName
		FROM products
		WHERE productId = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(ctx, query, productID)
	var product Product
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &product, nil
}

func getProductList() ([]Product, error) {
	query := `SELECT
		productId,
		manufacturer,
		sku,
		upc,
		pricePerUnit,
		quantityOnHand,
		productName
		FROM products`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName,
		)
		products = append(products, product)
	}
	return products, nil
}

func removeProduct(productID int) error {
	query := "DELETE FROM products WHERE productId = ?"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}
	return nil
}

func updateProduct(product Product) error {
	query := `UPDATE products SET
		manufacturer = ?,
		sku = ?,
		upc = ?,
		pricePerUnit = ?,
		quantityOnHand = ?,
		productName = ?
		WHERE productId = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, query,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID,
	)
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	query := `INSERT INTO products (
		manufacturer,
		sku,
		upc,
		pricePerUnit,
		quantityOnHand,
		productName)
		VALUES (?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(ctx, query,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
	)
	if err != nil {
		return 0, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(insertID), nil
}

func searchForProductData(productFilter ReportFilter) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT 
		productId, 
		LOWER(manufacturer), 
		LOWER(sku), 
		upc, 
		pricePerUnit, 
		quantityOnHand, 
		LOWER(productName) 
		FROM products WHERE `)
	if productFilter.NameFilter != "" {
		queryBuilder.WriteString(`productName LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.NameFilter)+"%")
	}
	if productFilter.ManufacturerFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`manufacturer LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.ManufacturerFilter)+"%")
	}
	if productFilter.SKUFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`sku LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.SKUFilter)+"%")
	}

	results, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName,
		)
		products = append(products, product)
	}
	return products, nil
}
