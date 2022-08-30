package handler_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

type ContainerAddress struct {
	Host string
	Port string
}

var _ = BeforeSuite(func() {
	fmt.Println("🟢 BeforeSuite Integration test")
	maria := setupMariaDBContainer()
	es := setupElasticSearchContainer()

	fmt.Println(maria, es)
})

func setupMariaDBContainer() ContainerAddress {

	ctx := context.Background()
	wd, _ := os.Getwd()
	wd += "/../../../seed/init.sql"

	mariadbContainerReq := testcontainers.ContainerRequest{
		Image:        "mariadb:10.5.8",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_USERNAME": "root",
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "books",
		},
		Mounts:     testcontainers.Mounts(testcontainers.BindMount(wd, "/docker-entrypoint-initdb.d/init.sql")),
		WaitingFor: wait.ForLog("3306  mariadb.org binary distribution").WithStartupTimeout(time.Second * 10),
	}

	mariaDBContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mariadbContainerReq,
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

	return ContainerAddress{mariaDBHost, mariaDBPort.Port()}
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
		WaitingFor: wait.ForLog("started").WithStartupTimeout(time.Second * 10),
	}

	esContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: esContainerReq,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("error starting es container: %s", err)
	}

	_, err = esContainer.Exec(ctx, []string{"sh", "/pre-test-script/es_container_db.sh"})

	if err != nil {
		log.Fatalf("esContainer.Exec: %s", err)
	}

	esHost, _ := esContainer.Host(ctx)

	esPort, err := esContainer.MappedPort(ctx, "9200")
	if err != nil {
		log.Fatalf("esContainer.MappedPort: %s", err)
	}

	return ContainerAddress{esHost, esPort.Port()}
}
