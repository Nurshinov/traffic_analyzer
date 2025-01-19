package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"traffic_analyzer/internal/tcp_packets"
)

type Client struct {
	db *sql.DB
}

func New() *Client {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}

	return &Client{db: db}
}

func (c *Client) Insert(address string, all, retransmitted int) error {
	if _, err := c.db.Exec("insert into IPStat (address, all_, retransmitted) values ($1, $2, $3)",
		address, all, retransmitted); err != nil {
		return err
	}

	return nil
}

func (c *Client) Update(address string, all, retransmitted int) error {
	if _, err := c.db.Exec("update IPStat set all_ = $1, retransmitted = $2 where address = $3",
		all, retransmitted, address); err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(address string) *tcp_packets.IPstat {
	IPstat := &tcp_packets.IPstat{}

	row := c.db.QueryRow("select * "+
		"from IPStat where address = $1", address)

	if err := row.Scan(&IPstat.IP, &IPstat.All,
		&IPstat.Retransmitted); err != nil {
		return nil
	}

	return IPstat
}

func (c *Client) Close() {
	c.db.Close()
}

func (c *Client) GetAll() []*tcp_packets.IPstat {
	ipSl := make([]*tcp_packets.IPstat, 0)

	rows, err := c.db.Query("select * from IPStat")
	if err != nil {
		log.Println(err)

		return nil
	}

	for rows.Next() {
		i := &tcp_packets.IPstat{}

		err := rows.Scan(&i.IP, &i.All, &i.Retransmitted)
		if err != nil {
			log.Println(err)

			continue
		}

		ipSl = append(ipSl, i)
	}

	return ipSl
}
