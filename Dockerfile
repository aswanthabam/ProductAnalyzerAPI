FROM python:3.11

WORKDIR /code

COPY ./ /code/

RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

EXPOSE 8000

CMD ["python", "main.py"]