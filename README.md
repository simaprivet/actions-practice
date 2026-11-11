go + postgres в docker compose + github actions

запуск вручную: cp .env.example .env && docker compose up -d --build
открыть в браузере: http://localhost:8081

статус контейнеров: docker compose ps
проверка приложения и бд: curl http://localhost:8081/health
данные через api: curl http://localhost:8081/api/items
количество записей через api: curl http://localhost:8081/api/items/count
данные прямо из бд: docker compose exec db psql -U ci_user -d ci_db -c "SELECT * FROM items;"
логи: docker compose logs

остановить (данные остаются): docker compose down
остановить и удалить данные: docker compose down -v

github actions:
workflow: .github/workflows/deploy.yml
запуск: push в main или Actions -> Deploy and Test -> Run workflow
runner: self-hosted, запущен у меня в WSL (~/actions-runner/run.sh)
секреты: Settings -> Secrets and variables -> Actions -> DB_USER, DB_PASSWORD, DB_NAME

по критериям:
self-hosted runner - runs-on: self-hosted, runner в Settings -> Actions -> Runners со статусом Idle/Active
docker compose - шаг Start stack: docker compose up -d, шаг Verify: проверяет что app и db running
секреты - ${{ secrets.* }} записываются в .env в шаге Create .env from secrets
артефакты - deployment-report (отчёт тестов, .env без значений) и docker-logs (логи, ps, images)
тест 1 - /health: приложение отвечает и пингует бд
тест 2 - сравнивает количество и названия записей из бд (psql) с тем что отдаёт веб (/api/items)
