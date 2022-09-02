package main

import (
	"context"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/grace"
	"github.com/guoyk93/grace/graceconf"
	"github.com/guoyk93/grace/gracemain"
	"github.com/guoyk93/grace/gracenotify"
	"github.com/guoyk93/grace/gracetrack"
)

func main() {
	var (
		err error

		ctx, _ = gracemain.WithSignalCancel(
			gracetrack.Init(context.Background()),
		)
	)

	defer gracemain.Exit(&err)
	defer gracenotify.Notify("[GITDUMP]", &ctx, &err)
	defer grace.Guard(&err)

	opts := grace.Must(graceconf.LoadYAMLFlagConf[gitdump.Options]())

	_ = gracemain.WriteLastRun(opts.Dir)
}
