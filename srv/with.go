package srv

import (
	mgo "github.com/cbwfree/micro-core/store/mongo"
	rds "github.com/cbwfree/micro-core/store/redis"
	"github.com/cbwfree/micro-core/web"
	"github.com/micro/go-micro/v2"
)

type WithAPP func(c *App)

func WithMongoDB(dbname string) WithAPP {
	return func(a *App) {
		a.Mongo = mgo.NewStore(dbname)
	}
}

func WithRedisDB(db int) WithAPP {
	return func(a *App) {
		a.Redis = rds.NewStore(db)
	}
}

func WithWebServer(name string) WithAPP {
	return func(a *App) {
		a.Web = web.NewServer(name)
	}
}

func WithPublisher(topic ...string) WithAPP {
	return func(a *App) {
		for _, tp := range topic {
			a.publisher[tp] = micro.NewEvent(tp, a.SrvClient())
		}
	}
}
