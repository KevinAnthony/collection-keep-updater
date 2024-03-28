package utils_test

//func TestGetFlagOrDefault(t *testing.T) {
//	t.Parallel()
//
//	Convey("GetFlagOrDefault", t, func() {
//		cmdMock := types.NewICommandMock(t)
//
//		flagCall := cmdMock.On("Flag", "test_key").Maybe()
//
//		flag := &pflag.Flag{Changed: true}
//		Convey("should return flag value when", func() {
//			Convey("flag is not nil and has been changed", func() {
//				flagCall.Once().Return(flag)
//
//				actual := utils.GetFlagOrDefault(cmdMock, "test_key", "yes", "no")
//
//				So(actual, ShouldEqual, "yes")
//			})
//		})
//		Convey("should return default when", func() {
//			Convey("flag returns nil", func() {
//				flagCall.Once().Return(nil)
//
//				actual := utils.GetFlagOrDefault(cmdMock, "test_key", "yes", "no")
//
//				So(actual, ShouldEqual, "no")
//			})
//			Convey("flag is not changed", func() {
//				flag.Changed = false
//				flagCall.Once().Return(flag)
//
//				actual := utils.GetFlagOrDefault(cmdMock, "test_key", "yes", "no")
//
//				So(actual, ShouldEqual, "no")
//			})
//		})
//	})
//}
