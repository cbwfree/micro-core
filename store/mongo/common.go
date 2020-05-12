package mgo

import (
	"context"
	"github.com/cbwfree/micro-core/conv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

// SelectOne 通过反射查询单条记录
func SelectOne(col *mongo.Collection, filter interface{}, model reflect.Type, options ...*options.FindOneOptions) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadWriteTimeout)
	defer cancel()

	if filter == nil {
		filter = bson.M{}
	}

	var one = conv.Elem(model).Addr().Interface()
	if err := col.FindOne(ctx, filter, options...).Decode(one); err != nil {
		return nil, err
	}

	return one, nil
}

// SelectAll 通过反射查询多条记录
func SelectAll(col *mongo.Collection, filter interface{}, model reflect.Type, options ...*options.FindOptions) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadWriteTimeout)
	defer cancel()

	if filter == nil {
		filter = bson.M{}
	}

	var cur, err = col.Find(ctx, filter, options...)
	if err != nil {
		return nil, err
	}

	var rows = conv.ElemSlice(model)
	if err := cur.All(context.Background(), rows.Addr().Interface()); err != nil {
		return nil, err
	}

	if rows.IsNil() {
		return nil, nil
	}

	var result []interface{}
	for i := 0; i < rows.Len(); i++ {
		result = append(result, rows.Index(i).Interface())
	}

	return result, nil
}

func FindOne(col *mongo.Collection, filter interface{}, result interface{}, options ...*options.FindOneOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadWriteTimeout)
	defer cancel()

	if filter == nil {
		filter = bson.M{}
	}

	if err := col.FindOne(ctx, filter, options...).Decode(result); err != nil {
		return err
	}

	return nil
}

func FindAll(col *mongo.Collection, filter interface{}, result interface{}, options ...*options.FindOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadWriteTimeout)
	defer cancel()

	if filter == nil {
		filter = bson.M{}
	}

	if cur, err := col.Find(ctx, filter, options...); err == nil {
		if err := cur.All(context.Background(), result); err != nil {
			return err
		}
	} else if err != mongo.ErrNilDocument {
		return err
	}

	return nil
}
