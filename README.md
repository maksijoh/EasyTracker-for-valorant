# Valorant Discord Rich Presence

Локальное Go-приложение получает ранг Valorant через HenrikDev API и раз в
минуту обновляет Discord Rich Presence текущего пользователя.

## Настройка

Скопируйте отслеживаемый шаблон `.env.example` в локальный `.env`:

```powershell
Copy-Item .env.example .env
```

Затем заполните `.env` своими значениями:

```env
ACCESS_TOKEN=ваш_ключ_HenrikDev
DISCORD_APP_ID=1538271778651111514
DISCORD_LARGE_IMAGE=valorant
RIOT_NAME=ваше_имя
RIOT_TAG=ваш_тег
RIOT_REGION=eu
```

Запустите Discord Desktop, затем приложение:

```powershell
go run .
```

Discord и приложение должны быть запущены от одного пользователя и с
одинаковым уровнем прав. Веб-эндпоинт истории матчей доступен по адресу:

```text
http://localhost:8080/api/data?name=ваше_имя&tag=ваш_тег&region=eu
```

## Иконка Discord

Загрузите `assets/valorant-presence.png` в Discord Developer Portal для своего
приложения: **Rich Presence → Art Assets → Add Image**. Имя asset должно быть
точно `valorant`. После загрузки подождите несколько минут и перезапустите
приложение.

Панель повторно отправляется Discord каждые 15 секунд, поэтому после выхода из
матча она восстанавливается автоматически без перезапуска программы.
