//Package module is the core package of elasthink
package module

import (
	"github.com/SurgicalSteel/elasthink/entity"
	"github.com/SurgicalSteel/elasthink/redis"
	"github.com/SurgicalSteel/elasthink/util"
)

//Module is the main struct to represent a core module
type Module struct {
	StopwordSet map[string]int
	Redis       *redis.Redis
}

var moduleObj *Module

//InitModule is a function that initializes a module object and its requirements (dependencies)
func InitModule(stopwordData entity.StopwordData, redisObject *redis.Redis) {
	moduleObj = new(Module)
	moduleObj.StopwordSet = util.CreateWordSet(stopwordData.Words)
	moduleObj.Redis = redisObject
}
