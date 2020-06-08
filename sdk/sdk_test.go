package sdk

import (
	"reflect"
	"testing"
	"time"

	er "github.com/SurgicalSteel/elasthink/redis"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestInitSDK(t *testing.T) {
	initializeSpec := InitializeSpec{
		RedisConfig: RedisConfig{
			Address:   "a.redis.address",
			MaxActive: 30,
			MaxIdle:   10,
			Timeout:   10,
		},
		SdkConfig: SdkConfig{
			IsUsingStopWordsRemoval: true,
			StopWordRemovalData:     getDummyStopwords(),
			AvailableDocumentType:   getDummyDocumentType(),
		},
	}

	expectedRedis := &er.Redis{
		Pool: &redigo.Pool{
			MaxIdle:     initializeSpec.RedisConfig.MaxIdle,
			MaxActive:   initializeSpec.RedisConfig.MaxActive,
			IdleTimeout: time.Duration(initializeSpec.RedisConfig.Timeout) * time.Second,
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", initializeSpec.RedisConfig.Address)
			},
		},
	}

	actualElasthinkSDK := Initialize(initializeSpec)
	assert.Equal(t, expectedRedis.Pool.MaxActive, actualElasthinkSDK.Redis.Pool.MaxActive)
	assert.Equal(t, expectedRedis.Pool.MaxIdle, actualElasthinkSDK.Redis.Pool.MaxIdle)
	assert.Equal(t, expectedRedis.Pool.IdleTimeout, actualElasthinkSDK.Redis.Pool.IdleTimeout)
	assert.NotNil(t, actualElasthinkSDK.Redis.Pool.Dial)
	assert.Equal(t, true, actualElasthinkSDK.isUsingStopWordsRemoval)

	if len(actualElasthinkSDK.availableDocumentType) != len(getDummyDocumentType()) {
		t.Fatal("Result 'DocumentTypes' is not same with what we expected.")
	}

	equalStopwords := reflect.DeepEqual(getDummyStopwords(), actualElasthinkSDK.stopWordRemovalData)
	if !equalStopwords {
		t.Fatal("Result 'Stopwords' is not same with what we expected.")
	}
}

//private functions
func getDummyInitializedSDK() ElasthinkSDK {
	return ElasthinkSDK{}
}

func getDummyStopwords() []string {
	return []string{"ada", "adalah", "adanya", "adapun", "waktu", "waktunya", "walau", "walaupun", "wong", "yaitu", "yakin", "yakni", "yang"}
}

func getDummyDocumentType() []string {
	return []string{"campaign", "advertisement"}
}
