package services

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseService handles Firebase operations
type FirebaseService struct {
	app    *firebase.App
	client *auth.Client
}

// NewFirebaseService creates a new Firebase service
func NewFirebaseService(credentialsPath, projectID string) (*FirebaseService, error) {
	ctx := context.Background()

	var opts []option.ClientOption
	if credentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsPath))
	}

	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %w", err)
	}

	log.Printf("Firebase service initialized for project: %s", projectID)

	return &FirebaseService{
		app:    app,
		client: client,
	}, nil
}

// VerifyIDToken verifies a Firebase ID token and returns the UID
func (s *FirebaseService) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := s.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %w", err)
	}

	return token.UID, nil
}

// GetUser retrieves user information by UID
func (s *FirebaseService) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := s.client.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return user, nil
}
