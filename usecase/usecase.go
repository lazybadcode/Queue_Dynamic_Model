package usecase

import (
	"context"
	"log"
	"queue/database"
	"time"
)

type UsecaseInterface interface {
	GetQueue(ctx context.Context, input map[string]interface{}) ([]map[string]interface{}, error)
	CreateQueue(ctx context.Context, idCard string, mobileNo string, input map[string]interface{}) (map[string]interface{}, error)
	UpdateQueue(ctx context.Context, ids string, newDate string, newSlot int, input map[string]interface{}) (map[string]interface{}, error)
	DeleteQueue(ctx context.Context, ids string) (map[string]interface{}, error)
	Batch()
	SendSmsAllToday()
}

type Usecase struct {
	db     *database.DB
	config *Config
	log    *log.Logger
}

type Config struct {
	QueuePerSlot           int           `mapstructure:"queue_per_slot"`
	QueuePerSlotDayOff     int           `mapstructure:"queue_per_slot_day_off"`
	QueuePerSlotSpecialDay int           `mapstructure:"queue_per_slot_special_day"`
	SlotDuration           time.Duration `mapstructure:"slot_duration"`
	StartTime              time.Duration `mapstructure:"start_time"`
	CloseTime              time.Duration `mapstructure:"close_time"`
	MaxDayForQueue         int           `mapstructure:"max_day_for_queue"`
	MaxRetryReserve        int           `mapstructure:"max_retry_reserve"`

	Batch Batch `mapstructure:"batch"`
}

func New(db *database.DB, conf *Config) UsecaseInterface {
	return &Usecase{
		db:     db,
		config: conf,
		log:    log.Default(),
	}
}
