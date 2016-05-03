# Sample Database
This is a seed database to help you get up and going.

Well actually - I should say "get much closer" to be up and going, as there are still some undocumented steps to 
work out for unsuspecting devs trying to build from this repo in its current state.

Anyway ...

It contains an almost-real dataset, at least the sites / machines and machine tools are correct for a particular real life site. 

To seed your Postgres database (as user postgres)

$ createdb -U postgres cmms
$ psql -U postgres cmms < {pathname to this dir}/cmms.sql

done !

Now, make sure that you configure your config.json file in the running dist directory to point to this database, ie :

{
	"DataSourceName": "user=postgres password=postgres dbname=cmms sslmode=disable",
	...
}

Porting to other databases would be possible, but the pain would be great, since in this application I am 
leaning very heavily on the excellent DAT toolkit, which is currently PostgreSQL only.

see : https://github.com/mgutz/dat