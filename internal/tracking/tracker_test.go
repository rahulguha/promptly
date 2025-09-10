package tracking

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDynamoDBClient is a mock implementation of the dynamodbiface.DynamoDBAPI.
type MockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *MockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *MockDynamoDBClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func TestNewDynamoDBTracker(t *testing.T) {
	tracker, err := NewDynamoDBTracker("us-east-1", "test-table")
	assert.NoError(t, err)
	assert.NotNil(t, tracker)
	assert.NotNil(t, tracker.db)
	assert.Equal(t, "test-table", tracker.tableName)
}

func TestUserExists(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	tracker := &DynamoDBTracker{
		db:        mockClient,
		tableName: "test-table",
	}

	// Case 1: User exists
	mockClient.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(1),
	}, nil).Once()
	exists, err := tracker.UserExists("test@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockClient.AssertExpectations(t)

	// Case 2: User does not exist
	mockClient.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(0),
	}, nil).Once()
	exists, err = tracker.UserExists("nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)

	// Case 3: DynamoDB error
	mockClient.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(nil, errors.New("dynamodb error")).Once()
	exists, err = tracker.UserExists("error@example.com")
	assert.Error(t, err)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}

func TestCreateUserRecord(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	tracker := &DynamoDBTracker{
		db:        mockClient,
		tableName: "test-table",
	}

	// Case 1: Successful creation
	mockClient.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil).Once()
	err := tracker.CreateUserRecord("user123", "newuser@example.com", "New User", 1678886400)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)

	// Case 2: DynamoDB error
	mockClient.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return(nil, errors.New("dynamodb error")).Once()
	err = tracker.CreateUserRecord("user456", "erroruser@example.com", "Error User", 1678886400)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestUpdateUserRecord(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	tracker := &DynamoDBTracker{
		db:        mockClient,
		tableName: "test-table",
	}

	// Case 1: Successful update
	mockClient.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(&dynamodb.UpdateItemOutput{}, nil).Once()
	err := tracker.UpdateUserRecord("user123", "updated@example.com", "Updated Name", 1678886500)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)

	// Case 2: DynamoDB error
	mockClient.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, errors.New("dynamodb error")).Once()
	err = tracker.UpdateUserRecord("user456", "errorupdate@example.com", "Error Update", 1678886500)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}
