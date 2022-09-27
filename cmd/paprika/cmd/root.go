package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const rootCmdUse = "paprika"

var (
	cfgFile string

	listenAddr string
	enableTls  bool
	tlsKey     string
	tlsCert    string

	grpcServerListener net.Listener
	grpcServer         *grpc.Server

	rootCmd = &cobra.Command{
		Use:   rootCmdUse,
		Short: "EventSauce - Paprika snapshot store",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			opts := []grpc.ServerOption{}

			if enableTls {

				cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
				if err != nil {
					log.Fatalln("Failed load TLS certificate (", tlsCert, ") or key (", tlsKey, ") because", err)
				}

				config := &tls.Config{
					Certificates: []tls.Certificate{cert},
					ClientAuth:   tls.VerifyClientCertIfGiven,
				}

				opts = append(opts, grpc.Creds(credentials.NewTLS(config)))

			}

			lis, err := net.Listen("tcp", listenAddr)
			if err != nil {
				log.Fatalln("Failed to listen to store address", listenAddr, "because", err)
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

	rootCmd.PersistentFlags().StringVar(&listenAddr, "listenAddr", "localhost:57755", "Bind address to store gRPC server")
	rootCmd.PersistentFlags().BoolVar(&enableTls, "tls", false, "Enable TLS protocol only on gRPC server")
	rootCmd.PersistentFlags().StringVar(&tlsKey, "tlsKey", "", "PEM encoded private key file path")
	rootCmd.PersistentFlags().StringVar(&tlsCert, "tlsCert", "", "PEM encoded certificate file path")

}

func serveGrpcServer() {

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
