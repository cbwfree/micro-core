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
func SelectOne(ctx context.Context, col *mongo.Collection, filter interface{}, model reflect.Type, options ...*options.FindOneOptions) (interface{}, error) {
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
func SelectAll(ctx context.Context, col *mongo.Collection, filter interface{}, model reflect.Type, options ...*options.FindOptions) ([]interface{}, error) {
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

// 查找单条数据
func FindOne(ctx context.Context, col *mongo.Collection, filter interface{}, result interface{}, options ...*options.FindOneOptions) error {
	if filter == nil {
		filter = bson.M{}
	}

	if err := col.FindOne(ctx, filter, options...).Decode(result); err != nil {
		return err
	}

	return nil
}

// 查找多条数据
func FindAll(ctx context.Context, col *mongo.Collection, filter interface{}, result interface{}, options ...*options.FindOptions) error {
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

// 分段获取数据
func FindScan(ctx context.Context, col *mongo.Collection, cur, size int64, filter interface{}, result interface{}, fn ...func(opts *options.FindOptions) *options.FindOptions) *Scan {
	if filter == nil {
		filter = bson.M{}
	}

	var scan = new(Scan)

	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return scan
	}

	scan = NewScan(count, cur, size)

	if count > 0 {
		opts := scan.FindOptions()
		if len(fn) > 0 && fn[0] != nil {
			opts = fn[0](opts)
		}

		cur, err := col.Find(ctx, filter, opts)
		if err != nil {
			return scan
		}

		if err := cur.All(ctx, result); err != nil {
			return scan
		}
	}

	return scan
}
