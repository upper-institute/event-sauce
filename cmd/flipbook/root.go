package flipbook

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/upper-institute/event-sauce/internal/eventstore"
	"github.com/upper-institute/event-sauce/internal/snapshotstore"
	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const rootCmdUse = "flipbook"

var (
	cfgFile string

	listenAddr string
	enableTls  bool
	tlsKey     string
	tlsCert    string

	grpcServerListener net.Listener
	grpcServer         *grpc.Server

	eventStoreService    = &eventstore.EventStoreServer{}
	snapshotStoreService = &snapshotstore.SnapshotStoreServer{}

	rootCmd = &cobra.Command{
		Use:   rootCmdUse,
		Short: "EventSauce - Snapshot store",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			opts := []grpc.ServerOption{}

			if enableTls {

				cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
				if err != nil {
					log.Fatalln("failed load TLS certificate (", tlsCert, ") or key (", tlsKey, ") because", err)
				}

				config := &tls.Config{
					Certificates: []tls.Certificate{cert},
					ClientAuth:   tls.VerifyClientCertIfGiven,
				}

				opts = append(opts, grpc.Creds(credentials.NewTLS(config)))

			}

			lis, err := net.Listen("tcp", listenAddr)
			if err != nil {
				log.Fatalln("failed to listen to store address", listenAddr, "because", err)
			}

			grpcServerListener = lis

			grpcServer = grpc.NewServer(opts...)

		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+rootCmdUse+".yaml)")

	rootCmd.PersistentFlags().StringVar(&listenAddr, "listenAddr", "0.0.0.0:6336", "Bind address to store gRPC server")
	rootCmd.PersistentFlags().BoolVar(&enableTls, "tls", false, "Enable TLS protocol only on gRPC server")
	rootCmd.PersistentFlags().StringVar(&tlsKey, "tlsKey", "", "PEM encoded private key file path")
	rootCmd.PersistentFlags().StringVar(&tlsCert, "tlsCert", "", "PEM encoded certificate file path")

}

func serveGrpcServer() {

	if eventStoreService.Backend != nil {
		apiv1.RegisterEventStoreServer(grpcServer, eventStoreService)
	}

	if snapshotStoreService.Backend != nil {
		apiv1.RegisterSnapshotStoreServer(grpcServer, snapshotStoreService)
	}

	reflection.Register(grpcServer)

	log.Println("Server listening at:", grpcServerListener.Addr())

	if err := grpcServer.Serve(grpcServerListener); err != nil {
		log.Fatalln("Failed to serve because", err)
	}

}

func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name "." + rootCmdUse (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("." + rootCmdUse)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
