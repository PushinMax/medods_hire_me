package main

import (
	"medods_hire_me/internal/blacklist"
	"medods_hire_me/internal/handler"
	"medods_hire_me/internal/mailer"
	"medods_hire_me/internal/repository"
	"medods_hire_me/internal/server"
	"medods_hire_me/internal/service"

	"log"
	"os"
	"os/signal"
	"syscall"

	"context"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("DB_USER"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}
	
	repository := repository.New(db)
	mailer := mailer.New()
	blacklist := blacklist.New()
	service := service.New(repository, mailer, blacklist)
	handler := handler.New(service)

	
	server := new(server.Server)

	go func() {
		if err := server.Run(viper.GetString("server.port"), handler.Init()); err != nil {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	<-ch
	_ = server.Shutdown(context.Background())
	_ = db.Close()

}




func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".") 
    viper.SetConfigType("yml")
	return viper.ReadInConfig()
}