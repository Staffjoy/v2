# Account Migrations

In development, create a new migration in this folder with:

```
migrate -url=$ACCOUNT_MYSQL_CONFIG -path=$STAFFJOY/account/migrations/ create <migration_file_name>
```

Migrations are automatically applied in dev. Migrations must be manually executed in staging and production. 

For additional documentation, please see the [migrate documentation](https://github.com/mattes/migrate).
