# Research Service Program

---

- Full description

Hello, this is a micro service that works with JSON data. Everything is deployed on a PostgreSQL database, the database code is included in the project. All server actions are recorded in logs. There are also backups for data recovery, of which there are two files. One file is created when loading the cache from the database; after this loading, the program will create a backup file in json format. And the most interesting file is endBackup.json, and it is interesting because this file is written only if the server is turned off, or if there are any errors. And since we have two files, almost two recovery points, data loss is minimized. Yes, I know it didn’t work out with the tests, but that’s it for now! And I would like to point out that the manual testing was successful! But the automatic one failed!

---

- Briefly

In short, everything is written in the go language, with the exception of only the site with its frontend. There is auto-start of all this stuff, quite good and detailed logs, and error handling.

---

- Команды для правильной работы,
Постройте образы и запустите контейнеры:
```shell
docker-compose up --build -d
```

- Проверьте работу приложения и Nginx через браузер или curl:
```shell
curl http://localhost:80
curl http://localhost:8080
docker-compose down
docker-compose up --build -d
```

- Проверьте изменения через браузер или curl:
```shell
curl http://localhost:80
```

---
