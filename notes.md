```
pg_ctl -D /usr/local/var/postgres -l /usr/local/var/postgres/server.log start

pg_ctl -D /usr/local/var/postgres stop -s -m fast

or

brew services start postgresql
```

```
psql -f init.sql
```

# Links

https://blog.codeship.com/unleash-the-power-of-storing-json-in-postgres/

http://coussej.github.io/2016/02/16/Handling-JSONB-in-Go-Structs/

http://coussej.github.io/2016/01/14/Replacing-EAV-with-JSONB-in-PostgreSQL/

https://github.com/jmoiron/sqlx


### For curl'ing

- title "JavaFuck2018"
- added_by "yyQ14wjI"
- start_date "2017-10-19 11:00:00+03"
- end_date "2017-10-19 11:00:00+03"
- description "kool yo"
- picture "pic url"
- country "USA"
- city "Los-Angeles"
- adress "Hollywood, 1"
- category "big data"
- min_price 100
- max_price 1000
- facebook_account "na"
- youtube_account "na"
- twitter_account  "na"
- tickets_available TRUE
- discount_program FALSE
- details "vse ok"
- speakers "{'speakers': ['bred fizpatrik']}"
- sponsors "{'sponsors': ['google']}"
- verified  TRUE

# Misc

// e := Entity{Id:1}
//
// err = db.QueryRow("SELECT name, description, properties FROM entity WHERE id = $1",
//               e.Id).Scan(&e.Name, &e.Description, &e.Properties)
//
// fmt.Printf("%+v\n", e)

// LISTEN TO EVENTS
// http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/


