package srv

import (
	mgo "github.com/cbwfree/micro-core/store/mongo"
	rds "github.com/cbwfree/micro-core/store/redis"
	"github.com/cbwfree/micro-core/web"
	"github.com/micro/go-micro/v2"
)

type WithAPP func(c *App)

func WithMongoDB(dbname string, opts ...mgo.Option) WithAPP {
	return func(a *App) {
		a.Mongo = mgo.NewStore(dbname, opts...)
	}
}

func WithRedisDB(db int, opts ...rds.Option) WithAPP {
	return func(a *App) {
		a.Redis = rds.NewStore(db, opts...)
	}
}

func WithWebServer(name string, opts ...web.Option) WithAPP {
	return func(a *App) {
		a.Web = web.NewServer(name, opts...)
	}
}

func WithPublisher(topic ...string) WithAPP {
	return func(a *App) {
		for _, tp := range topic {
			a.publisher[tp] = micro.NewEvent(tp, a.SrvClient())
		}
	}
}
