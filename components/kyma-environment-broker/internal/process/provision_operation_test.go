package process

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage"
)

func Test_Provision_RetryOperationOnce(t *testing.T) {
	// given
	memory := storage.NewMemoryStorage()
	operations := memory.Operations()
	opManager := NewProvisionOperationManager(operations)
	op := internal.ProvisioningOperation{}
	op.UpdatedAt = time.Now()
	retryInterval := time.Hour
	errMsg := errors.New("ups ... ")

	// this is required to avoid storage retries (without this statement there will be an error => retry)
	err := operations.InsertProvisioningOperation(op)
	require.NoError(t, err)

	// then - first call
	op, when, err := opManager.RetryOperationOnce(op, errMsg.Error(), errMsg, retryInterval, fixLogger())

	// when - first retry
	assert.True(t, when > 0)
	assert.Nil(t, err)

	// then - second call
	t.Log(op.UpdatedAt.String())
	op.UpdatedAt = op.UpdatedAt.Add(-retryInterval - time.Second) // simulate wait of first retry
	t.Log(op.UpdatedAt.String())
	op, when, err = opManager.RetryOperationOnce(op, errMsg.Error(), errMsg, retryInterval, fixLogger())

	// when - second call => no retry
	assert.True(t, when == 0)
	assert.NotNil(t, err)
}

func Test_Provision_RetryOperation(t *testing.T) {
	// given
	memory := storage.NewMemoryStorage()
	operations := memory.Operations()
	opManager := NewProvisionOperationManager(operations)
	op := internal.ProvisioningOperation{}
	op.UpdatedAt = time.Now()
	retryInterval := time.Hour
	errorMessage := "ups ... "
	errOut := errors.New("error occurred")
	maxtime := time.Hour * 3 // allow 2 retries

	// this is required to avoid storage retries (without this statement there will be an error => retry)
	err := operations.InsertProvisioningOperation(op)
	require.NoError(t, err)

	// then - first call
	op, when, err := opManager.RetryOperation(op, errorMessage, errOut, retryInterval, maxtime, fixLogger())

	// when - first retry
	assert.True(t, when > 0)
	assert.Nil(t, err)

	// then - second call
	t.Log(op.UpdatedAt.String())
	op.UpdatedAt = op.UpdatedAt.Add(-retryInterval - time.Second) // simulate wait of first retry
	t.Log(op.UpdatedAt.String())
	op, when, err = opManager.RetryOperation(op, errorMessage, errOut, retryInterval, maxtime, fixLogger())

	// when - second call => retry
	assert.True(t, when > 0)
	assert.Nil(t, err)
}
