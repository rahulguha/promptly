package tracking

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Tracker defines the interface for user tracking operations.
type Tracker interface {
	UserExists(email string) (bool, error)
	CreateUserRecord(userID, email, name string, timestamp int64) error
	UpdateUserRecord(userID, email, name string, timestamp int64) error
	CreateActivityLog(userID, email string, timestamp int64, activityType, activityResult string, activityDetails map[string]interface{}) error
}

// DynamoDBTracker implements the Tracker interface using AWS DynamoDB.
type DynamoDBTracker struct {
	db                dynamodbiface.DynamoDBAPI
	tableName         string // For Users table
	activityTableName string // For ActivityLogs table
}

// NewDynamoDBTracker creates a new tracker instance with a DynamoDB client.
func NewDynamoDBTracker(awsRegion, usersTableName, activityTableName string) (*DynamoDBTracker, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws session: %w", err)
	}

	return &DynamoDBTracker{
		db:                dynamodb.New(sess),
		tableName:         usersTableName,
		activityTableName: activityTableName,
	}, nil
}

// UserExists checks if a user with the given email exists in DynamoDB using the email-index GSI.
func (t *DynamoDBTracker) UserExists(email string) (bool, error) {
	log.Printf("Checking if user exists by email: %s", email)

	result, err := t.db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(t.tableName),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :e"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":e": {S: aws.String(email)},
		},
		Limit: aws.Int64(1),
		Select: aws.String("COUNT"),
	})

	if err != nil {
		return false, fmt.Errorf("failed to query dynamodb by email: %w", err)
	}

	return *result.Count > 0, nil
}

// CreateUserRecord creates a new user record in DynamoDB.
func (t *DynamoDBTracker) CreateUserRecord(userID, email, name string, timestamp int64) error {
	log.Printf("Creating user record for userID: %s, email: %s", userID, email)
	record := UserRecord{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Timestamp: timestamp,
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return fmt.Errorf("failed to marshal user record: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(t.tableName),
		ConditionExpression: aws.String("attribute_not_exists(user_id)"), // Ensure user_id does not exist
	}

	_, err = t.db.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put item to dynamodb: %w", err)
	}

	return nil
}

// UpdateUserRecord updates an existing user record in DynamoDB.
func (t *DynamoDBTracker) UpdateUserRecord(userID, email, name string, timestamp int64) error {
	log.Printf("Updating user record for userID: %s, email: %s", userID, email)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(t.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {S: aws.String(userID)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("name"),
			"#T": aws.String("timestamp"),
			"#E": aws.String("email"), // Include email in update expression if it can change
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {S: aws.String(name)},
			":t": {N: aws.String(fmt.Sprintf("%d", timestamp))},
			":e": {S: aws.String(email)},
		},
		UpdateExpression:    aws.String("SET #N = :n, #T = :t, #E = :e"),
		ReturnValues:        aws.String("UPDATED_NEW"),
	}

	_, err := t.db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update item in dynamodb: %w", err)
	}

	return nil
}

// CreateActivityLog creates a new activity log record in DynamoDB.
func (t *DynamoDBTracker) CreateActivityLog(userID, email string, timestamp int64, activityType, activityResult string, activityDetails map[string]interface{}) error {
	log.Printf("Creating activity log for userID: %s, activityType: %s", userID, activityType)
	record := ActivityLogRecord{
		UserID:        userID,
		Email:         email,
		Timestamp:     timestamp,
		ActivityType:  activityType,
		ActivityResult: activityResult,
		ActivityDetails: activityDetails,
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return fmt.Errorf("failed to marshal activity log record: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(t.activityTableName),
	}

	_, err = t.db.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put activity log item to dynamodb: %w", err)
	}

	return nil
}

// UserRecord represents the data structure for a user in DynamoDB.
type UserRecord struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

// ActivityLogRecord represents the data structure for an activity log in DynamoDB.
type ActivityLogRecord struct {
	UserID        string                 `json:"user_id"`
	Email         string                 `json:"email"`
	Timestamp     int64                  `json:"timestamp"`
	ActivityType  string                 `json:"activity_type"`
	ActivityResult string                 `json:"activity_result"`
	ActivityDetails map[string]interface{} `json:"activity_details,omitempty"`
}