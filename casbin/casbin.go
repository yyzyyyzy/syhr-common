package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	redis2 "github.com/yyzyyyzy/syhr-common/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type CasbinConf struct {
	ModelText string `json:"ModelText,optional,env=CASBIN_MODEL_TEXT"`
}

func (c *CasbinConf) getDefaultModel() string {
	return `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
`
}

func (c *CasbinConf) getModel() (model.Model, error) {
	modelText := c.ModelText
	if modelText == "" {
		modelText = c.getDefaultModel()
	}
	return model.NewModelFromString(modelText)
}

func (c *CasbinConf) NewEnforcer(adapter persist.Adapter) (*casbin.Enforcer, error) {
	m, err := c.getModel()
	logx.Must(err)

	enforcer, err := casbin.NewEnforcer(m, adapter)
	logx.Must(err)

	err = enforcer.LoadPolicy()
	logx.Must(err)
	return enforcer, nil
}

func (c *CasbinConf) NewWatcher(redisConf redis2.RedisConf, updateCallback func(string)) (persist.Watcher, error) {
	opts := &redis.Options{
		Network:  "tcp",
		Addr:     redisConf.Host,
		Username: redisConf.Username,
		Password: redisConf.Pass,
		DB:       redisConf.Db,
	}

	watcher, err := rediswatcher.NewWatcher(redisConf.Host, rediswatcher.WatcherOptions{
		Options:    *opts,
		Channel:    fmt.Sprintf("%s-%d", redis2.RedisCasbinChannel, redisConf.Db),
		IgnoreSelf: false,
	})
	logx.Must(err)

	err = watcher.SetUpdateCallback(updateCallback)
	logx.Must(err)

	return watcher, nil
}

func (c *CasbinConf) NewEnforcerWithWatcher(dbType, dsn string, redisConf redis2.RedisConf) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapter(dbType, dsn)
	logx.Must(err)

	enforcer, err := c.NewEnforcer(adapter)
	logx.Must(err)

	watcher, err := c.NewWatcher(redisConf, rediswatcher.DefaultUpdateCallback(enforcer))
	logx.Must(err)

	err = enforcer.SetWatcher(watcher)
	logx.Must(err)

	err = enforcer.SavePolicy()
	logx.Must(err)

	return enforcer, nil
}

func (c *CasbinConf) MustNewEnforcer(dbType, dsn string, conf redis2.RedisConf) *casbin.Enforcer {
	enforcer, err := c.NewEnforcerWithWatcher(dbType, dsn, conf)
	if err != nil {
		logx.Errorf("Failed to initialize Casbin: %v", err)
		logx.Errorf("Casbin initialization failed: %v", err)
	}
	return enforcer
}
