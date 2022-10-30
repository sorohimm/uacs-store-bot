package bot

import (
	"context"
	"errors"
	tele "gopkg.in/telebot.v3"
	stdl "log"
	"os"
	"time"
	"uacs_store_bot/internal/service/bot/config"
	"uacs_store_bot/internal/service/bot/handler"
	"uacs_store_bot/pkg/conf"
	"uacs_store_bot/pkg/log"
)

func NewService() *Service {
	return &Service{}
}

type Service struct {
}

func (o *Service) initConfigs(ctx context.Context) context.Context {
	appConf := &config.Config{}
	if err := conf.New(appConf); err != nil {
		if errors.Is(err, conf.ErrHelp) {
			os.Exit(0)
		}
		stdl.Fatalf("failed to read app config: %v", err)
	}
	return config.WithContext(ctx, appConf)
}

func (o *Service) initLogger(ctx context.Context, version, built, appName string) context.Context {
	appConf := config.FromContext(ctx)
	// init logger
	l, err := log.NewZap(
		appConf.Log.Level,
		appConf.Log.EncType)
	if err != nil {
		stdl.Fatalf("failed to init logger: %v", err)
	}
	logger := l.Sugar().With("v", version, "built", built, "app", appName)
	return log.CtxWithLogger(ctx, logger.Desugar())
}

func (o *Service) Init(ctx context.Context, version, built, appName string) {
	ctx = o.initConfigs(ctx)
	ctx = o.initLogger(ctx, version, built, appName)

	bot := o.initBot(ctx)

	o.handle(ctx, bot)

	defer func() {
		if r := recover(); r != nil {
			bot.Start()
		}
	}()

	bot.Start()
}

func (o *Service) initBot(ctx context.Context) *tele.Bot {

	logger := log.FromContext(ctx).Sugar()
	cfg := config.FromContext(ctx)

	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		logger.Fatal(err)
	}

	return b
}

func (o *Service) handle(ctx context.Context, b *tele.Bot) {
	h := handler.NewBotHandler(b)
	h.Serve(ctx)
}
