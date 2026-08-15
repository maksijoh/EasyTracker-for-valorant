# EasyTracker for Valorant

EasyTracker — локальное Windows-приложение, которое получает ранг игрока через
HenrikDev API и показывает его в Discord Rich Presence. Оно отображает текущий
ранг, RR и изменение рейтинга после последнего матча. После выхода из Valorant
панель восстанавливается автоматически, без перезапуска программы.

[Скачать последнюю версию](https://github.com/maksijoh/EasyTracker-for-valorant/releases/latest)

## Быстрый запуск

1. Скачайте `EasyTracker-v1.1.0-windows-amd64.exe` со страницы Releases.
2. Запустите Discord Desktop.
3. Запустите скачанный EXE двойным кликом.
4. При первом запуске заполните окно настройки и нажмите **Сохранить и запустить**.
5. Оставьте EasyTracker запущенным вместе с Discord.

Устанавливать Go, открывать терминал и вручную создавать `.env` для релизной
версии не требуется.

## Первоначальная настройка

Окно первого запуска попросит:

- **HenrikDev API key** — получить ключ можно в
  [HenrikDev Dashboard](https://api.henrikdev.xyz/dashboard/api-keys);
- **Discord Application ID** — находится в
  [Discord Developer Portal](https://discord.com/developers/applications) в
  разделе приложения **General Information**;
- **Riot Name** — часть Riot ID перед `#`;
- **Riot Tag** — часть Riot ID после `#`;
- **Регион** — `eu`, `na`, `latam`, `br`, `ap` или `kr`.

Настройки и API key сохраняются только локально в `.env` рядом с EXE. Этот файл
игнорируется Git и не должен публиковаться.

Чтобы снова открыть окно настройки:

```powershell
.\EasyTracker-v1.1.0-windows-amd64.exe --configure
```

## Иконка Discord

Владелец Discord-приложения должен один раз загрузить
`assets/valorant-presence.png` в **Rich Presence → Art Assets → Add Image** и
назвать asset точно `valorant`. После загрузки Discord может обновлять asset
несколько минут.

## Диагностика

EasyTracker пишет журнал работы в `EasyTracker.log` рядом с EXE. Если панель не
появляется, проверьте, что Discord Desktop и EasyTracker запущены от одного
пользователя и с одинаковым уровнем прав.

Локальный JSON-эндпоинт истории матчей:

```text
http://localhost:8080/api/data?name=ваше_имя&tag=ваш_тег&region=eu
```

## Сборка из исходников

```powershell
.\build.ps1
```

Готовый EXE появится в `dist/`. Для разработки также можно использовать
`go run .`.

## License

[MIT](LICENSE)

## Disclaimer

This project is not affiliated with or endorsed by Riot Games. VALORANT is a
trademark of Riot Games, Inc.
