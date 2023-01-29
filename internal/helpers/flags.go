package helpers

import (
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagGetter interface {
	GetInt64(string) int64
	GetDuration(string) time.Duration
	GetBool(string) bool
	GetString(string) string
	GetStringSlice(string) []string
}

type FlagBinder interface {
	BindInt64(string, int64, string)
	BindDuration(string, time.Duration, string)
	BindBool(string, bool, string)
	BindString(string, string, string)
	BindStringSlice(string, []string, string)
}

type FlagController interface {
	FlagBinder
	FlagGetter
}

type flagController struct {
	Viper   *viper.Viper
	FlagSet *pflag.FlagSet
}

func NewFlagController(v *viper.Viper, f *pflag.FlagSet) FlagController {
	return &flagController{
		Viper:   v,
		FlagSet: f,
	}
}

func flagSetName(key string) string {
	return strings.ReplaceAll(key, ".", "-")
}

func (f *flagController) bind(key, name string) {
	err := f.Viper.BindPFlag(key, f.FlagSet.Lookup(name))
	if err != nil {
		panic(err)
	}
}

func (f *flagController) BindBool(key string, value bool, usage string) {
	name := flagSetName(key)
	f.FlagSet.Bool(name, value, usage)
	f.bind(key, name)
}

func (f *flagController) GetBool(key string) bool {
	return f.Viper.GetBool(key)
}

func (f *flagController) BindString(key string, value string, usage string) {
	name := flagSetName(key)
	f.FlagSet.String(name, value, usage)
	f.bind(key, name)
}

func (f *flagController) GetString(key string) string {
	return f.Viper.GetString(key)
}

func (f *flagController) BindStringSlice(key string, value []string, usage string) {
	name := flagSetName(key)
	f.FlagSet.StringSlice(name, value, usage)
	f.bind(key, name)
}

func (f *flagController) GetStringSlice(key string) []string {
	return f.Viper.GetStringSlice(key)
}

func (f *flagController) BindDuration(key string, value time.Duration, usage string) {
	name := flagSetName(key)
	f.FlagSet.Duration(name, value, usage)
	f.bind(key, name)
}

func (f *flagController) GetDuration(key string) time.Duration {
	return f.Viper.GetDuration(key)
}

func (f *flagController) BindInt64(key string, value int64, usage string) {
	name := flagSetName(key)
	f.FlagSet.Int64(name, value, usage)
	f.bind(key, name)
}

func (f *flagController) GetInt64(key string) int64 {
	return f.Viper.GetInt64(key)
}
