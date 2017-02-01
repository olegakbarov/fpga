# Posgres

Start & stop manually:
```
pg_ctl -D /usr/local/var/postgres -l /usr/local/var/postgres/server.log start

pg_ctl -D /usr/local/var/postgres stop -s -m fast
```

```
psql -f init.sql
```

Links
https://blog.codeship.com/unleash-the-power-of-storing-json-in-postgres/
http://coussej.github.io/2016/02/16/Handling-JSONB-in-Go-Structs/
http://coussej.github.io/2016/01/14/Replacing-EAV-with-JSONB-in-PostgreSQL/
https://github.com/jmoiron/sqlx


# Go


// e := Entity{Id:1}
//
// err = db.QueryRow("SELECT name, description, properties FROM entity WHERE id = $1",
//               e.Id).Scan(&e.Name, &e.Description, &e.Properties)
//
// fmt.Printf("%+v\n", e)

// LISTEN TO EVENTS
// http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/
