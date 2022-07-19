package clock

type Option func(cron *Cron)

func WithSecond(isStartup bool) Option {
	return func(cron *Cron) {
		cron.Parse.IsWithSecond = isStartup
	}
}
