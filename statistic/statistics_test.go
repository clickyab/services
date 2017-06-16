package statistic_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"time"

	"github.com/fzerorubigd/services/statistic"
	"github.com/fzerorubigd/services/statistic/mock"
)

func TestStatisticStore(t *testing.T) {
	statistic.Register(mock.NewMockStatistic)
	Convey("test statistic store", t, func() {
		store := statistic.GetStatisticStore("test_key", 1*time.Hour)
		So(store.Key(), ShouldEqual, "test_key")
		Convey("check inc and dec", func() {

		})
	})
}
