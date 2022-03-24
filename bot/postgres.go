package bot

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"go-transaction-bot/config"
	"strconv"
	"time"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres() *Postgres {
	postgres := new(Postgres)
	err, pool := NewClient(context.Background())
	if err != nil {
		logrus.Fatal(err.Error())
	}
	postgres.pool = pool
	return postgres
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	i := 1
	for j := attempts; j > 0; {
		logrus.Info("Trying to connect to database attempt ", i, "(", attempts, ")")
		i += 1
		if err = fn(); err != nil {
			time.Sleep(delay)
			j--
			continue
		}
		return nil
	}
	return
}

func NewClient(ctx context.Context) (error, *pgxpool.Pool) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.Config.Postgres.Username,
		config.Config.Postgres.Password,
		config.Config.Postgres.Host,
		config.Config.Postgres.Port,
		config.Config.Postgres.Database)
	var pool *pgxpool.Pool

	err := DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, 5, 5*time.Second)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return nil, pool
}

func AddmembersToDB(goBot *discordgo.Session, psql *Postgres) {
	logrus.Info("Adding members to db")
	q := `INSERT INTO sellbuyinfo (userid, usernickname, selled, bought) VALUES ($1, $2, 0, 0);`
	after := ""
	for i := 0; i > -1; i++ {
		members, err := goBot.GuildMembers(config.Config.Discord.GuildID, after, 1000)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(i, len(members))
		if len(members) == 0 {
			break
		}
		for _, mem := range members {
			userid, err := strconv.ParseUint(mem.User.ID, 10, 64)
			if err != nil {
				fmt.Println(err.Error())
			}
			psql.pool.Exec(context.Background(), q, userid, mem.User.Username+"#"+mem.User.Discriminator)
		}
		after = members[len(members)-1].User.ID
	}

}
