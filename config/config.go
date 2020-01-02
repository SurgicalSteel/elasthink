//package config is where we parse configurations from configuration files (.ini) into configuration objects
package config

import (
	"log"
)

const configTag string = "[CONFIG]"

//InitConfig initializes all configuration (database, mq, api calls, redis, etc.)
func InitConfig(path, env string) error {

	// err := readDatabaseConfig(path, env)
	// if err != nil {
	// 	log.Println(configTag, "Error on reading DB config. Detail :", err.Error())
	// 	return err
	// }
	//
	// err = readMQConfig(path, env)
	// if err != nil {
	// 	log.Println(configTag, "Error on reading MQ config. Detail :", err.Error())
	// 	return err
	// }

	err := readRedisConfig(path, env)
	if err != nil {
		log.Println(configTag, "Error on reading Redis cofig. Detail :", err.Error())
		return err
	}

	return nil
}
