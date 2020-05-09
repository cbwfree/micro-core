package conf

import (
	"context"
	"github.com/cbwfree/micro-core/conv"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

// 配置项
type Conf struct {
	sync.RWMutex
	config config.Config
	data   interface{}
}

func (c *Conf) C() config.Config {
	return c.config
}

func (c *Conf) Data() interface{} {
	c.RLock()
	defer c.RUnlock()

	return c.data
}

// 载入web配置
func (c *Conf) LoadDB(ctx context.Context, col *mongo.Collection, opts ...*options.FindOptions) error {
	c.Lock()
	defer c.Unlock()

	var rows []*Model
	if cur, err := col.Find(ctx, bson.M{}, opts...); err == nil {
		if err := cur.All(context.Background(), &rows); err != nil {
			return err
		}
	} else if err != mongo.ErrNilDocument {
		return err
	}

	source := memory.NewSource(
		memory.WithJSON(conv.ToJson(convert(rows))),
	)

	if err := c.config.Load(source); err != nil {
		return err
	}
	if err := c.config.Scan(c.data); err != nil {
		return err
	}

	return nil
}

// NewConf ...
func NewConf(data interface{}) *Conf {
	c, _ := config.NewConfig()
	return &Conf{
		config: c,
		data:   data,
	}
}
