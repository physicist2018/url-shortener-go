package deleter

import (
	"context"
	"sync"
	"time"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/rs/zerolog"
)

const (
	maxQueueCapacity = 20              // number of records in the queue
	maxBatchSize     = 10              // number of records in the batch
	fushInterval     = 5 * time.Second // time interval for flushing the batch
)

type Deleter struct {
	service     domain.URLLinkService
	log         zerolog.Logger
	deleteQueue chan domain.URLLink
	mu          sync.Mutex
	closeOnce   sync.Once
}

func NewDeleter(service domain.URLLinkService, logger zerolog.Logger) *Deleter {
	return &Deleter{
		service:     service,
		log:         logger,
		deleteQueue: make(chan domain.URLLink, maxQueueCapacity),
	}
}

func (d *Deleter) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		var (
			batch       []domain.URLLink
			flushTicker = time.NewTicker(fushInterval)
		)

		// функция удаления пачки коротких url-ссылок
		flushBatch := func() {
			if len(batch) > 0 {
				if err := d.service.MarkURLsAsDeleted(ctx, batch); err != nil {
					d.log.Error().
						Err(err).
						Str("user_id", batch[0].UserID). // Предполагается, что в пачке один пользователь
						Msg("Ошибка при пометке ссылок на удаление")
				} else {
					d.log.Info().
						Int("количество удаленных ссылок", len(batch)).
						Str("user_id", batch[0].UserID).
						Msg("Ссылки успешно помечены на удаление")
				}
				batch = batch[0:] // Очищаем пачку после обработки
			}
		}

		for {
			select {
			case req, ok := <-d.deleteQueue:
				if !ok {
					d.log.Info().Msg("Канал удаления закрыт, завершаем горутину")
					flushBatch()
					return
				}

				batch = append(batch, req) // Добавляем в batch
				if len(batch) >= maxBatchSize {
					flushBatch()
				}

			case <-flushTicker.C:
				// Если таймер сработал, то отправляем текущую пачку на обработку
				flushBatch()

			case <-ctx.Done():
				d.log.Info().
					Msg("Получен сигнал завершения через контекст")
				flushBatch()
				return
			}
		}
	}()
}

func (d *Deleter) Enqueue(tasks ...domain.URLLink) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, t := range tasks {
		d.deleteQueue <- t
		d.log.Debug().Int("Задача", i).Str("Котортая ссылка", t.ShortURL).Msg("успешно добавлена в очередь для удаления")
	}
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
	return len(d.deleteQueue)
}
