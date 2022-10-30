package handler

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"io"
	"uacs_store_bot/pkg/log"
)

var (
	adminID int64 = 1380785175
)

func NewBotHandler(bot *tele.Bot) *BotHandler {
	return &BotHandler{
		bot,
	}
}

type BotHandler struct {
	*tele.Bot
}

var (
	textSubmitter          = &tele.ReplyMarkup{}
	mediaSubmitter         = &tele.ReplyMarkup{}
	btnUserSubmitWithText  = textSubmitter.Data("Подтвердить", "submit_text")
	btnUserSubmitWithMedia = mediaSubmitter.Data("Подтвердить", "submit_media")
)

func (o *BotHandler) Serve(ctx context.Context) {
	o.Handle(tele.OnMedia, o.handleMediaMessage(ctx))
	o.Handle(tele.OnText, o.handleTextMessage(ctx))
	o.Handle(&btnUserSubmitWithMedia, o.handleSubmitMedia(ctx))
	o.Handle(&btnUserSubmitWithText, o.handleSubmitText(ctx))
	o.Handle("/start", o.start())
}

func (o *BotHandler) start() tele.HandlerFunc {
	return func(c tele.Context) error {
		msg :=
			`Привет, чтобы опубликовать объявление напиши его сюда. 
Так же ты можешь добавлять фото и видео.🙂`
		return c.Send(msg)
	}
}

func (o *BotHandler) handleTextMessage(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		textSubmitter.Inline(
			textSubmitter.Row(btnUserSubmitWithText),
		)

		return c.Reply("Если все верно, нажмите 'Подтвердить', чтобы отправить пост на модерацию", textSubmitter)
	}
}

func (o *BotHandler) handleMediaMessage(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		mediaSubmitter.Inline(
			mediaSubmitter.Row(btnUserSubmitWithMedia),
		)

		return c.Reply("Если все верно, нажмите 'Подтвердить', чтобы отправить пост на модерацию", mediaSubmitter)
	}
}

func (o *BotHandler) handleSubmitMedia(ctx context.Context) tele.HandlerFunc {
	logger := log.FromContext(ctx).Sugar()
	return func(c tele.Context) error {
		var (
			err       error
			recipient = &tele.User{ID: adminID}
		)

		if _, err = o.Send(recipient, "Новый запрос на модерацию!"); err != nil {
			return err
		}

		if c.Callback().Message.ReplyTo.Photo != nil {
			if err = o.SendPhoto(c); err != nil {
				logger.Error(err)
				return err
			}
		}

		if c.Callback().Message.ReplyTo.Video != nil {
			if err = o.SendVideo(c); err != nil {
				logger.Error(err)
				return err
			}
		}

		err = c.Send("Ваше объявление отправлено на модерацию!")
		return err
	}
}

func (o *BotHandler) SendVideo(c tele.Context) error {
	var (
		reader    io.ReadCloser
		recipient = &tele.User{ID: adminID}
		err       error
	)

	if reader, err = o.File(c.Callback().Message.ReplyTo.Video.MediaFile()); err != nil {
		return err
	}
	vid := tele.Video{
		File:    tele.FromReader(reader),
		Caption: authorCaption(c),
	}

	if _, err = vid.Send(o.Bot, recipient, nil); err != nil {
		return err
	}

	return nil
}

func (o *BotHandler) SendPhoto(c tele.Context) error {
	var (
		reader    io.ReadCloser
		recipient = &tele.User{ID: adminID}
		err       error
	)

	if reader, err = o.File(c.Callback().Message.ReplyTo.Photo.MediaFile()); err != nil {
		return err
	}

	ph := tele.Photo{
		File:    tele.FromReader(reader),
		Caption: authorCaption(c),
	}

	if _, err = ph.Send(o.Bot, recipient, nil); err != nil {
		return err
	}

	return nil
}

func authorCaption(c tele.Context) string {
	return fmt.Sprintf("%s \n\nАвтор: @%s", c.Callback().Message.ReplyTo.Caption, c.Callback().Message.ReplyTo.Sender.Username)
}

func (o *BotHandler) handleSubmitText(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		_, err := o.Send(&tele.User{ID: adminID}, "Новый запрос на модерацию!")

		completeText := fmt.Sprintf("%s \n\nАвтор: @%s", c.Callback().Message.ReplyTo.Text, c.Callback().Message.ReplyTo.Sender.Username)
		_, err = o.Send(&tele.User{ID: adminID}, completeText)

		err = c.Send("Ваше объявление отправлено на модерацию!")

		return err
	}
}
