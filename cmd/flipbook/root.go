package flipbook

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/upper-institute/flipbook/internal"
	"github.com/upper-institute/flipbook/internal/helpers"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	rootCmdUse = "flipbook"

	grpcAddress_flag = "grpc.address"
)

var (
	cfgFile string

	server   *grpc.Server
	flagCont helpers.FlagController

	rootCmd = &cobra.Command{
		Use:   rootCmdUse,
		Short: "flipbook - Snapshot store",
		Run: func(cmd *cobra.Command, args []string) {

			addr := viper.GetString(grpcAddress_flag)

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalln("failed to listen to store address", addr, "because", err)
			}

			server = grpc.NewServer()

			storeAdapter := internal.GetStoreAdapter(flagCont)

			storeDriver, err := storeAdapter.New(flagCont)
			if err != nil {

			}

			store := internal.NewEventStore(storeDriver)

			flipbookv1.RegisterEventStoreServer(server, store)

			reflection.Register(server)

			log.Println("Server listening at:", lis.Addr())

			if err := server.Serve(lis); err != nil {
				log.Fatalln("Failed to serve because", err)
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

	flagCont.BindString(grpcAddress_flag, "0.0.0.0:6333", "Bind address for gRPC server listener")

	internal.BindWellKnownDrivers(flagCont)

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
