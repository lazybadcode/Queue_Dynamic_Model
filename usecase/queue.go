package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/vektah/gqlparser/v2/ast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"queue/model"
	"queue/utils"
	"time"
)

func (u *Usecase) GetQueue(ctx context.Context, input map[string]interface{}) ([]map[string]interface{}, error) {
	filter := mapToFilter(input)

	//opCtx := graphql.GetOperationContext(ctx) //TODO optimize query for next phase
	////project := graphql.CollectAllFields(ctx)
	//project := collectFields(opCtx.Operation.SelectionSet)
	//project = utils.RemoveFirstLayer(project)
	result, err := u.db.FindAllQueue(ctx, filter, nil)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}
	return result, nil
}

func mapToFilter(input map[string]interface{}) bson.M {
	filter := bson.M{}
	for k, v := range input {
		if v == nil {
			continue
		}
		if vv, ok := v.(string); ok && vv == "" {
			continue
		}
		if k == "_id" || k == "userId" {
			v = utils.ObjectIDFromHex(v.(string))
		}
		if k == "user" {
			continue //TODO cannot filter user, because query logic is: query queue and lookup user from queue
			//if vv, ok := v.(map[string]interface{}); ok {
			//	out := mapToFilter(vv)
			//	for kk, outv := range out {
			//		filter[fmt.Sprintf("%s.%s", "user", kk)] = outv
			//	}
			//	continue
			//}
		}
		filter[k] = v
	}

	return filter
}

func collectFields(selectionSet ast.SelectionSet) []string {
	var fields []string

	for _, selection := range selectionSet {
		switch field := selection.(type) {
		case *ast.Field:
			if len(field.SelectionSet) > 0 {
				subFields := collectFields(field.SelectionSet)
				for _, subField := range subFields {
					fields = append(fields, fmt.Sprintf("%s.%s", field.Name, subField))
				}
			} else {
				fields = append(fields, field.Name)
			}

		}
	}

	return fields
}

func (u *Usecase) CreateQueue(ctx context.Context, idCard string, mobileNo string, input map[string]interface{}) (map[string]interface{}, error) {
	user, err := u.db.FindOneUserByIdCard(ctx, idCard)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}
	if user == nil {
		mUser := map[string]interface{}{
			"idCard":      idCard,
			"mobileNo":    mobileNo,
			"updatedTime": time.Now(),
			"createdTime": time.Now(),
		}
		if inputUser, ok := input["user"].(map[string]interface{}); ok {
			updateUserDetails(mUser, inputUser)
		}

		userID, err := u.db.InsertUser(ctx, mUser)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}
		user, err = u.db.FindOneUser(ctx, userID)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}
	}
	// add logic queue manage here
	//var q *model.Queue
	var qid *primitive.ObjectID
	for i := 0; i < u.config.MaxRetryReserve; i++ { //retry when queue is empty and no error (queue empty when insert and duplicate key field date and no, prevent race condition case)
		qid, err = u.reserveQueue(ctx, user, input)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}
		if qid != nil { // exit when got queue
			break
		}
	}

	if qid == nil {
		u.log.Printf("[ERROR] Err %v", errors.New("cannot insert queue"))
		return nil, errors.New("cannot insert queue")
	}

	//return nil, nil
	return u.db.FindOneQueue(ctx, *qid)
}

func (u *Usecase) UpdateQueue(ctx context.Context, ids string, newDate string, newSlot int, input map[string]interface{}) (map[string]interface{}, error) {
	id, err := primitive.ObjectIDFromHex(ids)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}

	queuem, err := u.db.FindOneQueue(ctx, id)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}
	if queuem == nil {
		u.log.Printf("[ERROR] Err %v", errors.New("queue not found"))
		return nil, errors.New("queue not found")
	}
	var queue model.Queue
	model.FromMap(queuem, &queue)

	if newDate != "" {
		_, err = time.Parse("20060102", newDate)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, errors.New("invalid date format")
		}
	}
	if newSlot > 0 && newSlot >= u.slotPerDay() {
		u.log.Printf("[ERROR] Err %v", fmt.Errorf("slot more than limit %v", u.slotPerDay()-1))
		return nil, fmt.Errorf("slot more than limit %v", u.slotPerDay()-1)
	}

	var newNo int
	switch {
	case len(newDate) == 0 && newSlot == 0:
		// do nothing
	case len(newDate) > 0 && newSlot == 0: // fix date, auto select slot
		newSlot, newNo, err = u.validateAndCalculateNewSlot(ctx, newDate, 0)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}

	case len(newDate) == 0 && newSlot > 0: // fix slot for same date
		newSlot, newNo, err = u.validateAndCalculateNewSlot(ctx, queue.Date, newSlot)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}

	case len(newDate) > 0 && newSlot > 0: // fix date and slot
		newSlot, newNo, err = u.validateAndCalculateNewSlot(ctx, newDate, newSlot)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}
	}

	mQueue := make(map[string]interface{})
	if newDate != "" {
		mQueue["date"] = newDate
	}
	if newSlot > 0 {
		mQueue["slot"] = newSlot
	}
	if newNo > 0 {
		mQueue["no"] = newNo
	}
	updateQueueDetails(mQueue, input)
	err = u.db.UpdateQueue(ctx, id, mQueue)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}

	//update user
	mUser, ok := input["user"].(map[string]interface{})
	if ok && len(mUser) > 0 {
		err = u.db.UpdateUser(ctx, queue.UserID, mUser)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}
	}

	qs, err := u.db.FindOneQueue(ctx, queue.ID)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}
	if qs == nil || len(qs) == 0 {
		u.log.Printf("[ERROR] Err %v", errors.New("internal server error 1"))
		return nil, errors.New("internal server error 1")
	}
	return qs, nil
}

func (u *Usecase) DeleteQueue(ctx context.Context, ids string) (map[string]interface{}, error) {
	id, err := primitive.ObjectIDFromHex(ids)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}

	q, err := u.db.FindOneQueue(ctx, id)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}
	if q == nil {
		u.log.Printf("[ERROR] Err %v", errors.New("queue not found"))
		return nil, errors.New("queue not found")
	}

	err = u.db.DeleteQueue(ctx, id)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return nil, err
	}
	return q, nil
}

func (u *Usecase) reserveQueue(ctx context.Context, user *model.User, input map[string]interface{}) (*primitive.ObjectID, error) {
	slotPerDay := u.slotPerDay()
	for i := 0; i < u.config.MaxDayForQueue; i++ {
		day := utils.BeginOfDay(time.Now().AddDate(0, 0, i+1))
		maxQueuePerSlot, err := u.maxQueuePerSlot(ctx, day)
		slots, maxNo, err := u.db.ListSlotByDay(ctx, day)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}
		mSlot := make(map[int]int)
		slot := -1
		for _, v := range slots {
			mSlot[v.Slot] = v.Count
		}

		for j := 0; j < slotPerDay; j++ {
			if mSlot[j] < maxQueuePerSlot {
				slot = j
				break
			}
		}
		if slot < 0 {
			fmt.Printf("slot for day %v is full\n", day)
			continue
		}

		mQueue := map[string]interface{}{
			"userId":      user.ID,
			"no":          maxNo + 1,
			"date":        day.Add(u.config.StartTime + (time.Duration(slot) * u.config.SlotDuration)).Format("20060102"),
			"slot":        slot,
			"updatedTime": time.Now(),
			"createdTime": time.Now(),
		}
		updateQueueDetails(mQueue, input)

		id, err := u.db.InsertQueue(ctx, mQueue)
		if err != nil {
			u.log.Printf("[ERROR] Err %v", err)
			return nil, err
		}
		if id.IsZero() {
			return nil, nil
		}
		return &id, nil
	}

	u.log.Printf("[ERROR] Err %v", errors.New("no queue available"))
	return nil, errors.New("no queue available")
}

func (u *Usecase) validateAndCalculateNewSlot(ctx context.Context, date string, newSlot int) (int, int, error) {
	day, err := time.Parse("20060102", date)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return 0, 0, errors.New("invalid date format")
	}
	maxQueuePerSlot, err := u.maxQueuePerSlot(ctx, day)
	slots, maxNo, err := u.db.ListSlotByDay(ctx, day)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return 0, 0, err
	}
	mSlot := make(map[int]int)
	for _, v := range slots {
		mSlot[v.Slot] = v.Count
	}
	if newSlot > 0 {
		if mSlot[newSlot] >= maxQueuePerSlot {
			u.log.Printf("[ERROR] Err %v", "slot is full")
			return 0, 0, fmt.Errorf("slot is full")
		}
		return newSlot, maxNo + 1, nil
	}

	//case no fix new slot
	slot := -1
	slotPerDay := u.slotPerDay()
	for j := 0; j < slotPerDay; j++ {
		if mSlot[j] < maxQueuePerSlot {
			slot = j
			break
		}
	}
	if slot < 0 {
		u.log.Printf("[ERROR] Err %v", fmt.Errorf("slot for day %v is full\n", day))
		return 0, 0, fmt.Errorf("slot for day %v is full\n", day)
	}
	return slot, maxNo + 1, nil
}

func (u *Usecase) maxQueuePerSlot(ctx context.Context, t time.Time) (int, error) {
	var qPerSlot = u.config.QueuePerSlot
	sd, err := u.db.QuerySpecialDays(ctx, t)
	if err != nil {
		u.log.Printf("[ERROR] Err %v", err)
		return 0, err
	}
	if sd != nil {
		qPerSlot = u.config.QueuePerSlotSpecialDay
	}
	if t.Weekday() == time.Sunday || t.Weekday() == time.Saturday {
		qPerSlot = u.config.QueuePerSlotDayOff
	}

	return qPerSlot, nil
}

func (u *Usecase) slotPerDay() int {
	return int((u.config.CloseTime - u.config.StartTime) / u.config.SlotDuration)
}

func updateUserDetails(m, inputUser map[string]interface{}) {
	for k, v := range inputUser {
		if _, exists := m[k]; !exists {
			m[k] = v
		}
	}
}

func updateQueueDetails(m, inputUser map[string]interface{}) {
	for k, v := range inputUser {
		if k == "user" { //for model user
			continue
		}
		if _, exists := m[k]; !exists {
			m[k] = v
		}
	}
}
