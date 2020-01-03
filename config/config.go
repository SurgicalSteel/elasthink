//package config is where we parse configurations from configuration files (.ini) into configuration objects
package config

// Elasthink, An alternative to elasticsearch engine written in Go for small set of documents that uses inverted index to build the index and utilizes redis to store the indexes.
// Copyright (C) 2020 Yuwono Bangun Nagoro (a.k.a SurgicalSteel)
//
// This file is part of Elasthink
//
// Elasthink is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Elasthink is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
