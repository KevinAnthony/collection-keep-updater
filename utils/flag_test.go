package utils_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/pflag"
)

func TestGetFlagInt(t *testing.T) {
	t.Parallel()

	Convey("GetFlagInt", t, func() {
		cmdMock := types.NewICommandMock(t)
		flags := &pflag.FlagSet{}
		key := "fake flag"
		expected := 42

		flagCall := cmdMock.On("PersistentFlags").Maybe()

		Convey("should return flag value when key is found", func() {
			flags.Int(key, expected, "")
			flagCall.Return(flags)

			actual := utils.GetFlagInt(cmdMock, key)

			So(actual, ShouldResemble, expected)
		})
		Convey("should return default value when key is missing", func() {
			flagCall.Return(flags)

			actual := utils.GetFlagInt(cmdMock, key)

			So(actual, ShouldResemble, 0)
		})
	})
}

func TestGetFlagIntPtr(t *testing.T) {
	t.Parallel()
	Convey("GetFlagIntPtr", t, func() {
		cmdMock := types.NewICommandMock(t)
		flags := &pflag.FlagSet{}
		key := "fake flag"
		expected := 42

		flagCall := cmdMock.On("PersistentFlags").Maybe()

		Convey("should return flag value when key is found", func() {
			flags.Int(key, expected, "")
			flagCall.Return(flags)

			actual := utils.GetFlagIntPtr(cmdMock, key)

			So(*actual, ShouldResemble, expected)
		})
		Convey("should return default value when key is missing", func() {
			flagCall.Return(flags)

			actual := utils.GetFlagInt(cmdMock, key)

			So(actual, ShouldResemble, 0)
		})
	})
}

func TestGetFlagBool(t *testing.T) {
	t.Parallel()
	Convey("GetFlagBool", t, func() {
		cmdMock := types.NewICommandMock(t)
		flags := &pflag.FlagSet{}
		key := "fake flag"
		expected := true

		flagCall := cmdMock.On("PersistentFlags").Maybe()

		Convey("should return flag value when key is found", func() {
			flags.Bool(key, expected, "")
			flagCall.Return(flags)

			actual := utils.GetFlagBool(cmdMock, key)

			So(actual, ShouldResemble, expected)
		})
		Convey("should return default value when key is missing", func() {
			flagCall.Return(flags)

			actual := utils.GetFlagBool(cmdMock, key)

			So(actual, ShouldResemble, false)
		})
	})
}

func TestGetFlagIntSlice(t *testing.T) {
	t.Parallel()
	Convey("GetFlagIntSlice", t, func() {
		cmdMock := types.NewICommandMock(t)
		flags := &pflag.FlagSet{}
		key := "fake flag"
		expected := []int{42, 81}

		flagCall := cmdMock.On("PersistentFlags").Maybe()

		Convey("should return flag value when key is found", func() {
			flags.IntSlice(key, expected, "")
			flagCall.Return(flags)

			actual := utils.GetFlagIntSlice(cmdMock, key)

			So(actual, ShouldResemble, expected)
		})
		Convey("should return default value when key is missing", func() {
			flagCall.Return(flags)

			actual := utils.GetFlagIntSlice(cmdMock, key)

			So(actual, ShouldBeNil)
		})
	})
}

func TestGetFlagString(t *testing.T) {
	t.Parallel()
	Convey("GetFlagString", t, func() {
		cmdMock := types.NewICommandMock(t)
		flags := &pflag.FlagSet{}
		key := "fake flag"
		expected := "forty-two"

		flagCall := cmdMock.On("PersistentFlags").Maybe()

		Convey("should return flag value when key is found", func() {
			flags.String(key, expected, "")
			flagCall.Return(flags)

			actual := utils.GetFlagString(cmdMock, key)

			So(actual, ShouldResemble, expected)
		})
		Convey("should return default value when key is missing", func() {
			flagCall.Return(flags)

			actual := utils.GetFlagString(cmdMock, key)

			So(actual, ShouldBeEmpty)
		})
	})
}

func TestGetFlagStringPtr(t *testing.T) {
	t.Parallel()
	Convey("GetFlagStringPtr", t, func() {
		cmdMock := types.NewICommandMock(t)
		flags := &pflag.FlagSet{}
		key := "fake flag"
		expected := "forty-two"

		flagCall := cmdMock.On("PersistentFlags").Maybe()

		Convey("should return flag value when key is found", func() {
			flags.String(key, expected, "")
			flagCall.Return(flags)

			actual := utils.GetFlagStringPtr(cmdMock, key)

			So(*actual, ShouldResemble, expected)
		})
		Convey("should return default value when key is missing", func() {
			flagCall.Return(flags)

			actual := utils.GetFlagStringPtr(cmdMock, key)

			So(actual, ShouldBeNil)
		})
	})
}

func TestGetFlagStringSlice(t *testing.T) {
	t.Parallel()

	Convey("GetFlagStringSlice", t, func() {
		cmdMock := types.NewICommandMock(t)
		flags := &pflag.FlagSet{}
		key := "fake flag"
		expected := []string{"forty", "two"}

		flagCall := cmdMock.On("PersistentFlags").Maybe()

		Convey("should return flag value when key is found", func() {
			flags.StringSlice(key, expected, "")
			flagCall.Return(flags)

			actual := utils.GetFlagStringSlice(cmdMock, key)

			So(actual, ShouldResemble, expected)
		})
		Convey("should return default value when key is missing", func() {
			flagCall.Return(flags)

			actual := utils.GetFlagStringSlice(cmdMock, key)

			So(actual, ShouldBeNil)
		})
	})
}
