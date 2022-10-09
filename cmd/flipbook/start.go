package flipbook

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	snapshotStoreBackend string
	eventStoreBackend    string

	awsAccessKeyID     string
	awsSecretAccessKey string
	awsSessionToken    string
	awsRegion          string = "us-east-1"

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start Flipbook server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			switch eventStoreBackend {
			case "dynamodb":
				err = dynamodbEventStoreBackend()
			case "":
				break

			default:
				return fmt.Errorf("invalid event store backend: %s", eventStoreBackend)

			}

			if err != nil {
				return err
			}

			switch snapshotStoreBackend {
			case "redis":
				err = redisSnapshotStoreBackend()
			case "":
				break

			default:
				return fmt.Errorf("invalid snapshot store backend: %s", snapshotStoreBackend)

			}

			if err != nil {
				return err
			}

			serveGrpcServer()

			return nil

		},
	}
)

func init() {

	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVar(&snapshotStoreBackend, "snapshotStoreBackend", "", "SnapshotStore backend, supported values are: redis")
	startCmd.PersistentFlags().StringVar(&eventStoreBackend, "eventStoreBackend", "", "EventStore backend, supported values are: dynamodb")

	startCmd.PersistentFlags().StringVar(&awsAccessKeyID, "awsAccessKeyID", "", "AWS Access Key")
	startCmd.PersistentFlags().StringVar(&awsSecretAccessKey, "awsSecretAccessKey", "", "AWS Secret Access Key")
	startCmd.PersistentFlags().StringVar(&awsSessionToken, "awsSessionToken", "", "AWS Session Token")
	startCmd.PersistentFlags().StringVar(&awsRegion, "awsRegion", awsRegion, "AWS Default Region")

}
