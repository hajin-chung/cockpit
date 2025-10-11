FROM debian:bookworm

RUN apt update
RUN apt install -y axel ffmpeg

WORKDIR /cockpit

COPY ./build/cockpit .

EXPOSE 4000

CMD ["./cockpit"]
