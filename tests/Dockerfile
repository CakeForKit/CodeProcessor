FROM python:3.11-alpine

WORKDIR /app

# Копируем тесты и зависимости
COPY tests/tests3.py .
COPY requirements.txt .

# Устанавливаем зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Запускаем тесты
CMD ["pytest", "tests3.py", "-v", "--tb=short"]
