package conf

import (
	"context"
	"errors"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	ErrInvalid  = errors.New("invalid config filed")
	ErrNotFound = errors.New("not found filed")
)

// 系统配置数据模型
type Item struct {
	Field string      `json:"field"` // 设置字段
	Type  string      `json:"type"`  // 数据类型
	Value interface{} `json:"value"` // 设置值
}

// 配置项
type Conf struct {
	sync.RWMutex
	config config.Config
	source map[string]*Model
	data   map[string]interface{}
}

func (c *Conf) C() config.Config {
	return c.config
}

func (c *Conf) Source() map[string]*Model {
	c.RLock()
	defer c.RUnlock()

	return c.source
}

func (c *Conf) SetSource(rows []interface{}) {
	c.Lock()
	defer c.Unlock()

	c.source = make(map[string]*Model)
	for _, row := range rows {
		if d, ok := row.(*Model); ok {
			c.source[d.Field] = d
		}
	}
}

func (c *Conf) Data() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()

	return c.data
}

func (c *Conf) Model(field string) *Model {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.source[field]; ok {
		return v
	}
	return nil
}

func (c *Conf) Get(field string) interface{} {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v
	}
	return nil
}

func (c *Conf) Set(field string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	c.data[field] = value
	c.config.Set(value, field)
}

func (c *Conf) Bool(field string) bool {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.(bool)
	}
	return false
}

func (c *Conf) Str(field string) string {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.(string)
	}
	return ""
}

func (c *Conf) Int(field string) int {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.(int)
	}
	return 0
}

func (c *Conf) Int32(field string) int32 {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.(int32)
	}
	return 0
}

func (c *Conf) Int64(field string) int64 {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.(int64)
	}
	return 0
}

func (c *Conf) Float32(field string) float32 {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.(float32)
	}
	return 0
}

func (c *Conf) Float64(field string) float64 {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.(float64)
	}
	return 0
}

func (c *Conf) SliceStr(field string) []string {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.([]string)
	}
	return nil
}

func (c *Conf) SliceInt(field string) []int {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.([]int)
	}
	return nil
}

func (c *Conf) SliceInt32(field string) []int32 {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.([]int32)
	}
	return nil
}

func (c *Conf) SliceInt64(field string) []int64 {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.data[field]; ok {
		return v.([]int64)
	}
	return nil
}

// 载入web配置
func (c *Conf) LoadDB(ctx context.Context, col *mongo.Collection) error {
	c.Lock()
	defer c.Unlock()

	var rows []*Model
	if cur, err := col.Find(ctx, bson.M{}); err == nil {
		if err := cur.All(context.Background(), &rows); err != nil {
			return err
		}
	} else if err != mongo.ErrNilDocument {
		return err
	}

	source := memory.NewSource(
		memory.WithJSON(toDataJson(rows)),
	)

	if err := c.config.Load(source); err != nil {
		return err
	}

	if err := c.config.Scan(&c.data); err != nil {
		return err
	}

	return nil
}

// 更新配置
func (c *Conf) Update(ctx context.Context, col *mongo.Collection, field string, value string, reset ...bool) error {
	model := c.Model(field)
	if model == nil {
		return ErrNotFound
	}

	var opts = options.Update()
	if len(reset) > 0 && reset[0] {
		opts = opts.SetUpsert(true)
	}

	var update = bson.M{
		"field": model.Field,
		"type":  model.Type,
		"value": value,
	}
	if _, err := col.UpdateOne(ctx, bson.M{"field": field}, bson.M{"$set": update}, opts); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrInvalid
		}
		return err
	}

	c.Set(field, Convert(model.Type, value))

	return nil
}

// 重置配置
func (c *Conf) Reset(ctx context.Context, col *mongo.Collection, field string) (*Model, error) {
	model := c.Model(field)
	if model == nil {
		return nil, ErrNotFound
	}

	opts := options.Update().SetUpsert(true)
	if _, err := col.UpdateOne(ctx, bson.M{"field": field}, bson.M{"$set": model}, opts); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInvalid
		}
		return nil, err
	}

	c.Set(field, Convert(model.Type, model.Value))

	return model, nil
}

// NewConf ...
func NewConf() *Conf {
	c, _ := config.NewConfig()
	return &Conf{
		config: c,
		source: make(map[string]*Model),
		data:   make(map[string]interface{}),
	}
}
