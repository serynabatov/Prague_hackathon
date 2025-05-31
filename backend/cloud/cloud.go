package cloud

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	smpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func CreateUserSecret(ctx context.Context, client *secretmanager.Client, projectID string, userID uint) error {
	parent := fmt.Sprintf("projects/%s", projectID)
	secretID := fmt.Sprintf("%d-privatekey", userID)

	req := &smpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &smpb.Secret{
			Replication: &smpb.Replication{
				Replication: &smpb.Replication_Automatic_{
					Automatic: &smpb.Replication_Automatic{},
				},
			},
		},
	}

	_, err := client.CreateSecret(ctx, req)

	if err != nil {
		return fmt.Errorf("create secret error: %w", err)
	}

	return nil
}

func GetUserPrivateKey(ctx context.Context, client *secretmanager.Client, projectID string, userID uint) ([]byte, error) {
	name := fmt.Sprintf("projects/%s/secrets/%d-privatekey/versions/latest", projectID, userID)
	req := &smpb.AccessSecretVersionRequest{Name: name}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	return result.Payload.Data, nil
}

func NewSecretManagerClient(ctx context.Context) (*secretmanager.Client, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	return client, nil
}

func AddUserSecretVersion(ctx context.Context, client *secretmanager.Client, projectID string, userID uint, privateKey []byte) error {
	secretName := fmt.Sprintf("projects/%s/secrets/%d-privatekey", projectID, userID)
	req := &smpb.AddSecretVersionRequest{
		Parent: secretName,
		Payload: &smpb.SecretPayload{
			Data: privateKey,
		},
	}

	_, err := client.AddSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("add secret version error: %w", err)
	}
	return nil
}
