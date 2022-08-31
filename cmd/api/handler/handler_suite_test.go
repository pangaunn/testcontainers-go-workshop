package handler_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pangaunn/testcontainers-go-workshop/cmd/api/handler"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/datastore"
	logger "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

type ContainerAddress struct {
	Host      string
	Port      string
	Terminate func()
}

var (
	Maria  ContainerAddress
	ES     ContainerAddress
	Engine *gin.Engine
	Redis  ContainerAddress
)

var _ = BeforeSuite(func() {
	fmt.Println("🟢 BeforeSuite Integration test")
	Maria = setupMariaDBContainer()
	ES = setupElasticSearchContainer()
	Redis = setupRedisContainer()

	// init handler
	dbCredential := datastore.DatabaseCredential{
		Host:     Maria.Host,
		Port:     Maria.Port,
		Username: "root",
		Password: "root",
		Name:     "books",
	}
	esURL := fmt.Sprintf("http://%s:%s", ES.Host, ES.Port)
	redisURL := fmt.Sprintf("%s:%s", Redis.Host, Redis.Port)
	Engine = handler.InitHandler(dbCredential, esURL, redisURL)
})

var _ = AfterSuite(func() {
	fmt.Println("⛔️ AfterSuite Integration test")
	Maria.Terminate()
	ES.Terminate()
})

func setupMariaDBContainer() ContainerAddress {

	ctx := context.Background()
	wd, _ := os.Getwd()
	wd += "/../../../seed/init.sql"

	mariaDBContainerReq := testcontainers.ContainerRequest{
		Image:        "mariadb:10.5.8",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_USERNAME": "root",
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "books",
		},
		Mounts:     testcontainers.Mounts(testcontainers.BindMount(wd, "/docker-entrypoint-initdb.d/init.sql")),
		WaitingFor: wait.ForLog("3306  mariadb.org binary distribution").WithStartupTimeout(time.Second * 300),
	}

	mariaDBContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mariaDBContainerReq,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("error starting mariadb container: %s", err)
	}

	mariaDBHost, _ := mariaDBContainer.Host(ctx)

	mariaDBPort, err := mariaDBContainer.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("mariaDBContainer.MappedPort: %s", err)
	}

	terminateContainer := func() {
		logger.Info("terminating maria container...")
		if err := mariaDBContainer.Terminate(ctx); err != nil {
			log.Fatalf("error terminating maria container: %v\n", err)
		}
	}

	return ContainerAddress{mariaDBHost, mariaDBPort.Port(), terminateContainer}
}

func setupElasticSearchContainer() ContainerAddress {
	ctx := context.Background()
	wd, _ := os.Getwd()
	wd += "/../../../seed/es"

	esContainerReq := testcontainers.ContainerRequest{
		Image:        "elasticsearch:7.17.6",
		ExposedPorts: []string{"9200/tcp"},
		Env: map[string]string{
			"xpack.security.enabled": "false",
			"discovery.type":         "single-node",
		},
		Mounts:     testcontainers.Mounts(testcontainers.BindMount(wd, "/pre-test-script")),
		WaitingFor: wait.ForLog("started").WithStartupTimeout(time.Second * 300),
	}

	esContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: esContainerReq,
		Started:          true,
	})

	if err != nil {
		logger.Fatalf("error starting es container: %s", err)
	}

	_, err = esContainer.Exec(ctx, []string{"sh", "/pre-test-script/es_container_db.sh"})

	if err != nil {
		logger.Fatalf("esContainer.Exec: %s", err)
	}

	esHost, _ := esContainer.Host(ctx)

	esPort, err := esContainer.MappedPort(ctx, "9200")
	if err != nil {
		logger.Fatalf("esContainer.MappedPort: %s", err)
	}

	terminateContainer := func() {
		logger.Info("terminating es container...")
		if err := esContainer.Terminate(ctx); err != nil {
			logger.Fatalf("error terminating es container: %v\n", err)
		}
	}

	return ContainerAddress{esHost, esPort.Port(), terminateContainer}
}

func setupRedisContainer() ContainerAddress {
	ctx := context.Background()
	redisContainerRequest := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
	}
	redisContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: redisContainerRequest,
		Started:          true,
	})

	if err != nil {
		logger.Fatalf("error starting redis container: %s", err)
	}

	redisHost, _ := redisContainer.Host(ctx)
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		logger.Fatalf("redisContainer.MappedPort: %s", err)
	}

	terminateContainer := func() {
		logger.Info("terminating es container...")
		if err := redisContainer.Terminate(ctx); err != nil {
			logger.Fatalf("error terminating es container: %v\n", err)
		}
	}

	return ContainerAddress{redisHost, redisPort.Port(), terminateContainer}
}