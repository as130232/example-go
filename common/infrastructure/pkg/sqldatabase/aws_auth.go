package sqldatabase

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/go-sql-driver/mysql"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

type AWSAuthProvider struct {
	Region         string
	DbEndpoint     string
	DbUser         string
	DbName         string
	AuthTokenMux   *sync.Mutex
	AuthTokenCache *cache.Cache
}

var once sync.Once

const AuthTokenCacheKey = "AuthToken"

func NewAWSAuthProvider(region string, dbEndpoint string, dbUser string, dbName string) *AWSAuthProvider {
	provider := &AWSAuthProvider{Region: region, DbEndpoint: dbEndpoint, DbUser: dbUser, DbName: dbName}
	provider.AuthTokenMux = &sync.Mutex{}
	provider.AuthTokenCache = cache.New(10*time.Minute, 10*time.Minute)
	return provider
}

func (a *AWSAuthProvider) Register() {
	once.Do(func() {
		RegisterRDSMysqlCerts()
	})
}

func (a *AWSAuthProvider) DataSourceName() string {
	return a.dataSourceName("emptyToken")
}

func (a *AWSAuthProvider) dataSourceName(token string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true&tls=rds&charset=utf8mb4&parseTime=True&loc=%s",
		a.DbUser, token, a.DbEndpoint, a.DbName, "Asia%2fTaipei")
}

func (a *AWSAuthProvider) lockAndAuthToken() (string, string) {
	authToken := ""
	defer a.AuthTokenMux.Unlock()
	a.AuthTokenMux.Lock()
	result, exist := a.AuthTokenCache.Get(AuthTokenCacheKey)
	useCache := "true"
	if exist {
		authToken = result.(string)
	} else {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(a.Region))
		if err != nil {
			panic(err)
		}
		authToken, err = auth.BuildAuthToken(context.TODO(), a.DbEndpoint, a.Region, a.DbUser, cfg.Credentials)
		if err != nil {
			panic(err)
		}
		a.AuthTokenCache.Set(AuthTokenCacheKey, authToken, 10*time.Minute)
		useCache = "false"
	}

	return authToken, useCache
}

func (a *AWSAuthProvider) AuthToken() string {
	start := time.Now()
	authToken := ""
	result, exist := a.AuthTokenCache.Get(AuthTokenCacheKey)
	useCache := "true"
	if exist {
		authToken = result.(string)
	} else {
		authToken, useCache = a.lockAndAuthToken()
	}

	log.Println(fmt.Sprintf("[AWSAuthProvider] AuthToken authToken:%s,useCache:%s,elapse:%s", authToken, useCache, time.Since(start)))

	return authToken
}

func RegisterRDSMysqlCerts() {
	resp, err := http.DefaultClient.Get("https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem")
	if err != nil {
		panic("failed to download RDS certificate: " + err.Error())
	}

	pem, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("failed to read RDS certificate: " + err.Error())
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		panic("failed to append RDS certificate to root cert pool")
	}

	err = mysql.RegisterTLSConfig("rds", &tls.Config{RootCAs: rootCertPool, InsecureSkipVerify: true})
	if err != nil {
		panic("failed to register RDS certificate: " + err.Error())
	}
}
