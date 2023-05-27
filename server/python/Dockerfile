FROM python:3.11.3-slim-buster
ENV PYTHONUNBUFFERED 1
RUN mkdir /grpc_server
WORKDIR /grpc_server
ADD . /grpc_server/
RUN pip install --upgrade pip
RUN pip install -r requirements.txt
CMD ["python", "server.py"]