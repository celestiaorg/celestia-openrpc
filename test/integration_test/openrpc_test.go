package openrpc_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	openrpc "github.com/rollkit/celestia-openrpc"
)

type IntegrationTestSuite struct {
	suite.Suite

	dockerCompose *testcontainers.LocalDockerCompose
}

func (i *IntegrationTestSuite) SetupSuite() {
	composeFilePaths := []string{"docker/test-docker-compose.yml"}
	identifier := strings.ToLower(uuid.New().String())

	i.dockerCompose = testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)
	 i.dockerCompose.WaitForService("bridge0",
	 wait.ForHTTP("/header/1").WithPort("26659").
	 	WithStartupTimeout(60*time.Second).
	 	WithPollInterval(3*time.Second))
	execError := i.dockerCompose.WithCommand([]string{"up", "-d"}).Invoke()
	err := execError.Error
	if err != nil {
		i.Fail("failed to execute docker compose up:", "error: %v\nstdout: %v\nstderr: %v", err, execError.Stdout, execError.Stderr)
	}
}

func (i *IntegrationTestSuite) TearDownSuite() {
	execError := i.dockerCompose.Down()
	if err := execError.Error; err != nil {
		i.Fail("failed to execute docker compose down", "error: %v\nstdout: %v\nstderr: %v", err, execError.Stdout, execError.Stderr)
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (i *IntegrationTestSuite) TestNewClient() {
	client, err := openrpc.NewClient(context.TODO(), "http://localhost:26658", "test-jwt-token")
	i.Require().NoError(err)
	i.NotNil(client)
}

func (i *IntegrationTestSuite) TestDataRoundTrip() {
	time.Sleep(30 * time.Second)
	client, err := openrpc.NewClient(context.Background(), "http://localhost:26659", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.z1FNFoxCRayehAnJHlQsU9kBuCIEmwfESW_1drgJTdk")

	i.Require().NoError(err)
	defer client.Close()
	i.Require().NotNil(client)

	ctx, closer := context.WithTimeout(context.Background(), 1*time.Second)
	defer closer()

	_, err = client.Header.GetByHeight(ctx, 1)
	i.Require().NoError(err)
}
