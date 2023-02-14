FROM golang:1.18.1

WORKDIR /home
COPY ./pkg /home

RUN cd /home && go build -o library-vivek-sing

CMD ["/home/library-vivek-sing"]
