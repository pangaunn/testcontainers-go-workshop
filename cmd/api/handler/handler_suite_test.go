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

var _ = BeforeSuite(func() {
	fmt.Println("🟢 BeforeSuite Integration test")
	SetupContainer()
})

func SetupContainer() {

	ctx := context.Background()
	wd, _ := os.Getwd()
	wd = wd + "/../../../seed/init.sql"
	fmt.Println(wd)
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

	fmt.Println(mariaDBHost, mariaDBPort)

	// time.Sleep(100000 * time.Hour)

}
