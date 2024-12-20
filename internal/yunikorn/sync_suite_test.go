package yunikorn

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/unicorn-history-server/test/database"
)

type SyncSuite struct {
	suite.Suite
	tp   *database.TestPostgresContainer
	pool *pgxpool.Pool
}

func (ts *SyncSuite) SetupSuite() {
	ctx := context.Background()
	cfg := database.InstanceConfig{
		User:     "test",
		Password: "test",
		DBName:   "template",
		Host:     "localhost",
		Port:     15434,
	}

	tp, err := database.NewTestPostgresContainer(ctx, cfg)
	require.NoError(ts.T(), err)
	ts.tp = tp
	err = tp.Migrate("../../migrations")
	require.NoError(ts.T(), err)

	ts.pool = tp.Pool(ctx, ts.T(), &cfg)
}

func (ts *SyncSuite) TearDownSuite() {
	err := ts.tp.Container.Terminate(context.Background())
	require.NoError(ts.T(), err)
}

func (ts *SyncSuite) TestSubSuites() {
	ts.T().Run("SyncNodesIntTest", func(t *testing.T) {
		pool := database.CloneDB(t, ts.tp, ts.pool)
		suite.Run(t, &SyncNodesIntTest{pool: pool})
	})
	ts.T().Run("SyncQueuesIntTest", func(t *testing.T) {
		pool := database.CloneDB(t, ts.tp, ts.pool)
		suite.Run(t, &SyncQueuesIntTest{pool: pool})
	})
	ts.T().Run("SyncPartitionIntTest", func(t *testing.T) {
		pool := database.CloneDB(t, ts.tp, ts.pool)
		suite.Run(t, &SyncPartitionIntTest{pool: pool})
	})
	ts.T().Run("SyncApplicationsIntTest", func(t *testing.T) {
		pool := database.CloneDB(t, ts.tp, ts.pool)
		suite.Run(t, &SyncApplicationsIntTest{pool: pool})
	})
}

func TestSyncIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	topSuite := new(SyncSuite)
	suite.Run(t, topSuite)
}
