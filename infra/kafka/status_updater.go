package kafka

import (
	"context"
	"log"
	"messaggio/internal/repository"
	"sync"
)

// StatusUpdater представляет структуру для асинхронного обновления статуса.
type StatusUpdater struct {
	messageRepo repository.MessageRepository
	statusChan  chan statusUpdate
	wg          sync.WaitGroup
}

// NewStatusUpdater создает новый экземпляр StatusUpdater.
func NewStatusUpdater(messageRepo repository.MessageRepository) *StatusUpdater {
	su := &StatusUpdater{
		messageRepo: messageRepo,
		statusChan:  make(chan statusUpdate, 100),
	}

	// Запускаем горутину для обработки обновлений статуса
	go su.processStatusUpdates()

	return su
}

// processStatusUpdates обрабатывает обновления статусов асинхронно.
func (su *StatusUpdater) processStatusUpdates() {
	for update := range su.statusChan {
		su.wg.Add(1)
		go func(update statusUpdate) {
			defer su.wg.Done()
			err := su.messageRepo.UpdateMessageStatus(context.Background(), int64(update.status), update.messageID)
			if err != nil {
				log.Printf("Error updating message status: %v", err)
			}
		}(update)
	}
}

// UpdateStatus отправляет обновление статуса в канал.
func (su *StatusUpdater) UpdateStatus(messageID int64, status int64) {
	su.statusChan <- statusUpdate{messageID: messageID, status: status}
}

// Close закрывает канал и ожидает завершения всех обновлений статуса.
func (su *StatusUpdater) Close() {
	close(su.statusChan)
	su.wg.Wait()
}

type statusUpdate struct {
	messageID int64
	status    int64
}
