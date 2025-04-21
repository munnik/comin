package store

import (
	"testing"

	"github.com/nlewo/comin/internal/deployer"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentCommitAndLoad(t *testing.T) {
	tmp := t.TempDir()
	filename := tmp + "/state.json"
	s := New(filename, 2, 2)
	err := s.Commit()
	assert.Nil(t, err)

	s1 := New(filename, 2, 2)
	err = s1.Load()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(s.Deployments))

	s.DeploymentInsert(deployer.Deployment{UUID: "1", Operation: "switch"})
	_ = s.Commit()
	assert.Nil(t, err)

	s1 = New(filename, 2, 2)
	err = s1.Load()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(s.Deployments))
}

func TestLastDeployment(t *testing.T) {
	s := New("", 2, 2)
	ok, _ := s.LastDeployment()
	assert.False(t, ok)
	s.DeploymentInsert(deployer.Deployment{UUID: "1", Operation: "switch"})
	s.DeploymentInsert(deployer.Deployment{UUID: "2", Operation: "switch"})
	ok, last := s.LastDeployment()
	assert.True(t, ok)
	assert.Equal(t, "2", last.UUID)
}

func TestDeploymentInsert(t *testing.T) {
	s := New("", 2, 2)
	var hasEvicted bool
	var evicted deployer.Deployment
	hasEvicted, _ = s.DeploymentInsert(deployer.Deployment{UUID: "1", Operation: "switch"})
	assert.False(t, hasEvicted)
	hasEvicted, _ = s.DeploymentInsert(deployer.Deployment{UUID: "2", Operation: "switch"})
	assert.False(t, hasEvicted)
	hasEvicted, evicted = s.DeploymentInsert(deployer.Deployment{UUID: "3", Operation: "switch"})
	assert.True(t, hasEvicted)
	assert.Equal(t, "1", evicted.UUID)
	expected := []deployer.Deployment{
		{UUID: "3", Operation: "switch"},
		{UUID: "2", Operation: "switch"},
	}
	assert.Equal(t, expected, s.DeploymentList())

	hasEvicted, _ = s.DeploymentInsert(deployer.Deployment{UUID: "4", Operation: "test"})
	assert.False(t, hasEvicted)
	hasEvicted, _ = s.DeploymentInsert(deployer.Deployment{UUID: "5", Operation: "test"})
	assert.False(t, hasEvicted)
	hasEvicted, evicted = s.DeploymentInsert(deployer.Deployment{UUID: "6", Operation: "test"})
	assert.True(t, hasEvicted)
	assert.Equal(t, "4", evicted.UUID)
	expected = []deployer.Deployment{
		{UUID: "6", Operation: "test"},
		{UUID: "5", Operation: "test"},
		{UUID: "3", Operation: "switch"},
		{UUID: "2", Operation: "switch"},
	}
	assert.Equal(t, expected, s.DeploymentList())

	hasEvicted, evicted = s.DeploymentInsert(deployer.Deployment{UUID: "7", Operation: "switch"})
	assert.True(t, hasEvicted)
	assert.Equal(t, "2", evicted.UUID)
	hasEvicted, evicted = s.DeploymentInsert(deployer.Deployment{UUID: "8", Operation: "switch"})
	assert.True(t, hasEvicted)
	assert.Equal(t, "3", evicted.UUID)
}
