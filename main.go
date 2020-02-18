//Package main is where all the magic happens ;)
package main

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
	"context"
	"encoding/json"
	"flag"
	"github.com/SurgicalSteel/elasthink/config"
	"github.com/SurgicalSteel/elasthink/entity"
	"github.com/SurgicalSteel/elasthink/module"
	"github.com/SurgicalSteel/elasthink/redis"
	"github.com/SurgicalSteel/elasthink/router"
	"github.com/SurgicalSteel/elasthink/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const stopwordsFileName string = "files/data/stopwords_id.json"
const configPath string = "files/config"

func main() {
	log.SetOutput(os.Stdout)
	environmentFlag := flag.String("env", "development", "specify your environment for running elasthink (development / staging / production)")
	stopwordsRemovalUsageFlag := flag.Bool("swr", false, "option to use stopwords removal during create index & update index & searching (default false)")

	flag.Parse()

	environment := util.GetEnv(*environmentFlag)
	log.Println("Environment for elasthink:", environment)

	isUsingStopwordsRemoval := *stopwordsRemovalUsageFlag

	//read stop words file
	stopwordData, err := readStopwordsFile(stopwordsFileName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	//init config
	err = config.InitConfig(configPath, environment)
	if err != nil {
		log.Fatalln(err)
		return
	}

	//init redis
	redisObject := redis.InitRedis(*config.GetRedisConfig())

	//init entity data
	entity.Entity.Initialize(stopwordData)

	//init module
	module.InitModule(entity.Entity.GetStopwordData(), redisObject, isUsingStopwordsRemoval)

	routing := router.InitializeRoute()
	routing.RegisterHandler()
	routing.RegisterAppHandler()
	routing.RegisterInternalHandler()
	server := &http.Server{
		Addr: "0.0.0.0:9000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      routing.Router,
	}
	log.Println("Starting elasthink in port 9000...")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Failed to start service! Reason :", err.Error())
		}
	}()
	log.Println("Server Started")

	<-done
	log.Println("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		// extra handling here
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v\n", err)
	}

	log.Println("Server Exited Properly")
	log.Println("👋")
}

func readStopwordsFile(fileName string) (entity.StopwordData, error) {
	var stopwordData entity.StopwordData

	rawStopwordsFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Failed to open stopwords file. Reason :", err.Error())
		return stopwordData, err
	}
	defer rawStopwordsFile.Close()

	rawStopwordsBody, err := ioutil.ReadAll(rawStopwordsFile)
	if err != nil {
		log.Println("Failed to read stopwords file. Reason :", err.Error())
		return stopwordData, err
	}

	err = json.Unmarshal(rawStopwordsBody, &stopwordData)
	if err != nil {
		log.Println("Failed to unmarshal raw stopwords file. Reason :", err.Error())
		return stopwordData, err
	}

	return stopwordData, nil
}
