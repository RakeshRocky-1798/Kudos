package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"kleos/db/shedlock_db"
	"time"
)

type LockerDb struct {
	LockTime int `json:"lock_time" 20 sec`
	c        *cron.Cron
}

func NewLockerDbWithLockTime(lockFor int) *LockerDb {
	return &LockerDb{LockTime: lockFor, c: cron.New()}
}

func (l LockerDb) AddFunc(name string, spec string, shedlockService *shedlock_db.ShedlockRepository, cmd func()) error {

	fmt.Println("Adding cron job", name, spec)
	if l.c == nil {
		l.c = cron.New()
	}
	err := l.c.AddFunc(spec, func() {
		if l.DoLock(name, shedlockService) {
			defer shedlockService.Unlock(name, "")
			cmd()
		}
	})

	return err
}

func (l LockerDb) DoLock(name string, shedlockService *shedlock_db.ShedlockRepository) bool {
	s := shedlockService.Insert(name, LocalHostName(), l.LockTime)
	if true == s {
		ticker := time.NewTicker(1 * time.Second)

		go func() {
			for {
				select {
				case t := <-ticker.C:
					{
						fmt.Println("Ticking at", t)
						shedlockService.Unlock(name, "ticker")
						ticker.Stop()
						return
					}
				}
			}
		}()

	}
	return s
}

func (l LockerDb) Start() {
	l.c.Start()
}

func (l LockerDb) Stop() {
	l.c.Stop()
	l.c = nil
}
