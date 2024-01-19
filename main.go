package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cron "github.com/robfig/cron/v3"
)

func main() {
	postgres, err := NewPostgres()
	if err != nil {
		log.Panic(err)
	}

	err = postgres.Init()
	if err != nil {
		log.Println(err)
		return
	}

	// set scheduler berdasarkan zona waktu sesuai kebutuhan
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	// stop scheduler tepat sebelum fungsi berakhir
	defer scheduler.Stop()

	// set task yang akan dijalankan scheduler
	// gunakan crontab string untuk mengatur jadwal
	scheduler.AddFunc("0 03 * * *", func() { BirthDayScheduler(postgres) })

	// start scheduler
	go scheduler.Start()

	// trap SIGINT untuk trigger shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func SendAutomail(automailType string) {
	// ... instruksi untuk mengirim automail berdasarkan automailType
	fmt.Printf(time.Now().Format("2006-01-02 15:04:05") + " SendAutomail " + automailType + " telah dijalankan.\n")
}

func BirthDayScheduler(p *Postgres) {
	users, err := p.GetBirthdayData()
	if err != nil {
		log.Println(err)
		return
	}

	for _, user := range users {
		user.GeneratePromoCode(p)
	}

	err = p.GeneratePromoCode(users)
	if err != nil {
		log.Println(err)
		return
	}

	blastingEmail(users)

}

func (u *User) GeneratePromoCode(p *Postgres) error {
	promoCode := fmt.Sprintf("%s%s", u.Username, u.BirthDate.Format("02"))
	validUntil := time.Date(u.BirthDate.Year(), u.BirthDate.Month(), u.BirthDate.Day(), 23, 59, 59, 0, u.BirthDate.Location())

	promoId, err := p.GetPromo()
	if err != nil {
		return err
	}
	sp := SpecialPromo{
		PromoCode:  promoCode,
		ValidUntil: validUntil,
		Id_Promo:   promoId,
	}
	u.Promo = sp

	return nil
}

func blastingEmail(users []User) {
	for _, user := range users {
		fmt.Println("Email ke ", user.Email, "telah dikirim")
		time.Sleep(3 * time.Second)
	}
}
