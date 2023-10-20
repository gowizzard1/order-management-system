package product

import (
	"context"
	"errors"
	"github.com/leta/order-management-system/orders/pkg/models"
	"time"

	"cloud.google.com/go/firestore"
	db "github.com/leta/order-management-system/orders/db/firebase"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/product"
	"google.golang.org/api/iterator"
)

//var _ repository.ProductRepository = (*ProductRepository)(nil)

type productRepository struct {
	db *db.FirestoreService
}

func NewProductRepository(db *db.FirestoreService) product.ProductRepositoryInterface {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CheckPreconditions() {
	if r.db == nil {
		panic("no DB service provided")
	}
}

func (r *productRepository) productCollection() *firestore.CollectionRef {
	r.CheckPreconditions()

	return r.db.Client.Collection("products")
}

func (r *productRepository) CreateProduct(ctx context.Context, product *product.Product) (*product.Product, error) {
	r.CheckPreconditions()

	// Set CreatedAt and UpdatedAt to the current time
	currentTime := time.Now()

	product.CreatedAt = currentTime.Format(time.RFC3339)
	product.UpdatedAt = currentTime.Format(time.RFC3339)

	err := product.Validate()
	if err != nil {
		return nil, err
	}

	productModel := r.marshallProduct(product)

	docRef, _, writeErr := r.productCollection().Add(ctx, productModel)
	if writeErr != nil {
		return nil, writeErr
	}

	product.Id = docRef.ID

	return product, nil
}

func (r *productRepository) GetProduct(ctx context.Context, id string) (*product.Product, error) {
	r.CheckPreconditions()

	if id == "" {
		return nil, errors.New("some error")
	}

	docRef, err := r.productCollection().Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	productModel := &models.ProductModel{}
	if err := docRef.DataTo(productModel); err != nil {
		return nil, err
	}

	product := r.unmarshallProduct(productModel)

	product.Id = docRef.Ref.ID

	return product, nil
}

func (r *productRepository) ListProducts(ctx context.Context) ([]*product.Product, error) {
	r.CheckPreconditions()

	iter := r.productCollection().Documents(ctx)

	var products []*product.Product

	for {
		docRef, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		productModel := &models.ProductModel{}
		if err := docRef.DataTo(productModel); err != nil {
			return nil, err
		}

		product := r.unmarshallProduct(productModel)
		product.Id = docRef.Ref.ID

		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, id string, update *product.ProductUpdate,
) (*product.Product, error) {
	r.CheckPreconditions()

	product, getErr := r.GetProduct(ctx, id)
	if getErr != nil {
		return nil, getErr
	}

	if p := update.Name; p != nil {
		product.Name = *p
	}

	if p := update.Description; p != nil {
		product.Description = *p
	}

	if p := update.Price; p != nil {
		product.Price = *p
	}

	err := product.Validate()
	if err != nil {
		return nil, err
	}

	timeNow := time.Now()
	product.UpdatedAt = timeNow.Format(time.RFC3339)

	productModel := r.marshallProduct(product)

	_, err = r.productCollection().Doc(id).Set(ctx, productModel)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, id string) error {
	r.CheckPreconditions()

	_, err := r.productCollection().Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	return err
}

func (r *productRepository) marshallProduct(product *product.Product) *models.ProductModel {
	return &models.ProductModel{
		Name:        product.Name,
		Description: product.Description,
		Price:       int(product.Price),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func (r *productRepository) unmarshallProduct(productModel *models.ProductModel) *product.Product {
	return &product.Product{
		Name:        productModel.Name,
		Description: productModel.Description,
		Price:       uint(productModel.Price),
		CreatedAt:   productModel.CreatedAt,
		UpdatedAt:   productModel.UpdatedAt,
	}
}
