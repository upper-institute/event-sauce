package flipbook

import (
	"fmt"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/upper-institute/flipbook/internal"
	"github.com/upper-institute/flipbook/internal/helpers"
	"github.com/upper-institute/flipbook/internal/logging"
	"github.com/upper-institute/flipbook/internal/server"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	rootCmdUse = "flipbook"
)

var (
	cfgFile string

	flagCont helpers.FlagController

	rootCmd = &cobra.Command{
		Use:   rootCmdUse,
		Short: "flipbook - Event store controller",
		Run: func(cmd *cobra.Command, args []string) {

			logging.Load(flagCont)

			defer logging.Flush()

			log := logging.Logger.Named("RootCmd").Sugar()

			listener := server.CreateListener(flagCont)

			opts := []grpc.ServerOption{
				grpc.StreamInterceptor(
					grpc_middleware.ChainStreamServer(
						otelgrpc.StreamServerInterceptor(),
						grpc_zap.StreamServerInterceptor(logging.Logger.Named("GrpcServer")),
					),
				),
				grpc.UnaryInterceptor(
					grpc_middleware.ChainUnaryServer(
						otelgrpc.UnaryServerInterceptor(),
						grpc_zap.UnaryServerInterceptor(logging.Logger.Named("GrpcServer")),
					),
				),
			}

			grpcServer := grpc.NewServer(opts...)

			storeAdapter := internal.GetStoreAdapter(flagCont)

			storeDriver, err := storeAdapter.New(flagCont)
			if err != nil {
				log.Fatalw("Failed to instantiate store driver", "error", err)
			}

			defer storeAdapter.Destroy(storeDriver)

			store := internal.NewEventStore(storeDriver)

			flipbookv1.RegisterEventStoreServer(grpcServer, store)

			reflection.Register(grpcServer)

			log.Infow("Server listening", "address", listener.Addr())

			if err := grpcServer.Serve(listener); err != nil {
				log.Fatalw("Failed to serve gRPC server", "error", err)
			}

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

	flagCont = helpers.NewFlagController(viper.GetViper(), rootCmd.PersistentFlags())

	internal.BindWellKnownDrivers(flagCont)
	logging.BindOptions(flagCont)
	server.BindOptions(flagCont)

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
