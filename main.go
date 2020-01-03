//Package main is where all the magic happens ;)
package main

import (
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
	"time"
)

const stopwordsFileName string = "files/data/stopwords.json"
const configPath string = "files/config"

func main() {
	log.SetOutput(os.Stdout)
	environmentFlag := flag.String("env", "development", "specify your environment for running elasthink (development / staging / production)")

	flag.Parse()

	environment := util.GetEnv(*environmentFlag)
	log.Println("Starting elasthink in environment:", environment)

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
	redisObject, err := redis.InitRedis(*config.GetRedisConfig())
	if err != nil {
		log.Fatalln(err)
		return
	}

	//init module
	module.InitModule(stopwordData, redisObject)

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

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start service! Reason :", err.Error())
	}

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
