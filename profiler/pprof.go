package profiler

import (
	"context"
	"path/filepath"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/random"

	"github.com/fzerorubigd/expand"
	"github.com/pkg/profile"
	"github.com/sirupsen/logrus"
)

var (
	appDir, _ = expand.Pwd()
	mode      = config.RegisterString("services.profiler.mode", "disable", "services profiler status")
	root      = config.RegisterString("services.profiler.root", appDir, "services profiler root")
)

type initSystem struct {
	profiler interface {
		Stop()
	}
}

func (pi *initSystem) Initialize(ctx context.Context) {
	done := ctx.Done()
	assert.NotNil(done, "[BUG] the context is not supporting cancel")
	switch mode.String() {
	case "cpu":
		pi.profiler = profile.Start(
			profile.CPUProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(root.String(), <-random.ID)),
		)
	case "mem":
		pi.profiler = profile.Start(
			profile.MemProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(root.String(), <-random.ID)),
		)
	case "trace":
		pi.profiler = profile.Start(
			profile.TraceProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(root.String(), <-random.ID)),
		)
	case "block":
		pi.profiler = profile.Start(
			profile.BlockProfile,
			profile.NoShutdownHook,
			profile.ProfilePath(filepath.Join(root.String(), <-random.ID)),
		)
	default:
		logrus.Debug("Profiler disabled")
	}

	go func() {
		<-done
		if pi.profiler != nil {
			pi.profiler.Stop()
			logrus.Debug("Profiler done")
		}
	}()
}

func init() {
	initializer.Register(&initSystem{}, -100)
}
