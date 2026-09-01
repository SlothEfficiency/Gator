In order to run gator you need postgres and Go installed.

To install gator, run 'go install' in the roo of this directory.

The .gatorconfig.json should be in your home directory and look something like this:
{"db_url":"postgres://username:passwort@localhost:5432/gator?sslmode=disable","current_user_name":"currently_active_user_of_cli"}

