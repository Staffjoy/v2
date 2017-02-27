# Account Migrations

In development, create a new migration in this folder with:

```
migrate -url=$COMPANY_MYSQL_CONFIG -path=$STAFFJOY/company/migrations/ create <migration_file_name>
```

Migrations are automatically applied in dev. Migrations must be manually executed in staging and production. 

For additional documentation, please see the [migrate documentation](https://github.com/mattes/migrate).
