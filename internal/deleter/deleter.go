package deleter

import (
	"context"
	"sync"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/rs/zerolog"
)

const (
	maxQueueCapacity = 10 // number of records in the queue
)

type Deleter struct {
	service     domain.URLLinkService
	log         zerolog.Logger
	deleteQueue chan domain.DeleteRecordTask
	mu          sync.Mutex
	closeOnce   sync.Once
}

func NewDeleter(service domain.URLLinkService, logger zerolog.Logger) *Deleter {
	return &Deleter{
		service:     service,
		log:         logger,
		deleteQueue: make(chan domain.DeleteRecordTask, maxQueueCapacity),
	}
}

func (d *Deleter) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case req, ok := <-d.deleteQueue:
				if !ok {
					d.log.Info().Msg("Канал удаления закрыт, завершаем горутину")
					return
				}

				if err := d.service.MarkURLsAsDeleted(ctx, req); err != nil {
					d.log.Error().
						Err(err).
						Str("user_id", req.UserID).
						Strs("short_urls", req.ShortURLs).
						Msg("Ошибка при пометке ссылок на удаление")
				} else {
					d.log.Info().
						Int("количество удаленных ссылок", len(req.ShortURLs)).
						Str("user_id", req.UserID).
						Msg("Ссылки успешно помечены на удаление")
				}

			case <-ctx.Done():
				d.log.Info().
					Msg("Получен сигнал завершения через контекст")
				return
			}
		}
	}()
}

func (d *Deleter) Enqueue(task domain.DeleteRecordTask) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.deleteQueue <- task
	// Задача успешно добавлена в очередь
	d.log.Debug().Msg("Задача успешно добавлена в очередь")

	// select {
	// case d.deleteQueue <- task:
	// 	// Задача успешно добавлена в очередь
	// 	d.log.Debug().Msg("Задача успешно добавлена в очередь")
	// default:
	// 	d.log.Warn().Msg("Очередь удаления переполнена, задача не добавлена")
	// }
}

func (d *Deleter) Close() {
	d.closeOnce.Do(func() {
		if d.deleteQueue != nil {
			close(d.deleteQueue)
			d.log.Info().Msg("Канал удаления закрыт")
		}
	})
}

func (d *Deleter) Size() int {
	if d.deleteQueue == nil {
		return 0
	}
	return len(d.deleteQueue)
}
