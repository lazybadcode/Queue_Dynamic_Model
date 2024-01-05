package usecase

import (
	"context"
	"fmt"
	"queue/model"
	"queue/utils"
	"time"
)

type Batch struct {
	SendSmsTime time.Duration `mapstructure:"send_sms_time"`
}

func (u *Usecase) Batch() {
	u.config.Batch.SendSmsTime.Minutes()
	ticker := time.NewTicker(1 * time.Second) // Check time every minute
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			tt := time.Time{}.Add(u.config.Batch.SendSmsTime)
			if t.Hour() == tt.Hour() && t.Minute() == tt.Minute() && t.Second() == 0 {
				go u.SendSmsAllToday()
			}
		}
	}

	//u.SendSmsAllToday() // use this line instead for force run when start service
}

func (u *Usecase) SendSmsAllToday() {
	u.log.Printf("[INFO] Start SendSmsAllToday")
	defer u.log.Printf("[INFO] Stop SendSmsAllToday")

	start := utils.BeginOfDay(time.Now())
	stop := start.AddDate(0, 0, 1)
	queue, err := u.db.FindQueueInCreatedTimeRange(context.TODO(), start, stop)
	if err != nil {
		u.log.Printf("[ERROR] FindQueueInCreatedTimeRange error %v", err)
		return
	}

	for _, qm := range queue {
		var q model.Queue
		model.FromMap(qm, &q)

		if q.User == nil {
			u.log.Printf("[ERROR] queue id %v no have user data", q.ID)
			continue
		}

		slotStart := u.config.StartTime + u.config.SlotDuration*time.Duration(q.Slot)
		slotEnd := slotStart + u.config.SlotDuration

		msg := fmt.Sprintf("เรียน %s, คุณได้รับหมายเลขคิวที่ %d สำหรับวันที่ %s. กรุณามาตามเวลาในช่องเวลา %s - %s. หากมีคำถามโปรดติดต่อกลับ Call Center ที่หมายเลขโทรศัพท์: %s",
			q.User.Name,
			q.No,
			q.Date,
			time.Time{}.Add(slotStart).Format("15:04"),
			time.Time{}.Add(slotEnd).Format("15:04"),
			"0123456789")
		mobileNo := q.User.MobileNo
		u.log.Printf("[INFO] Send SMS message %v for mobile no %v", msg, mobileNo)
		//TODO Connect SMS server to send message here
	}
}
