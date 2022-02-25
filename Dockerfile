FROM golang:1.16.14-alpine as build

RUN apk add --update --no-cache bash git
WORKDIR /kubernetes-pod-monitor

COPY . .

RUN go build -o kubernetes-pod-monitor

FROM golang:1.16.14-alpine

WORKDIR /kubernetes-pod-monitor

RUN apk add --update --no-cache py3-pip
RUN pip install awscli==1.22.55

RUN adduser --disabled-password  --gecos "" shivam
USER shivam

COPY --from=build /kubernetes-pod-monitor/kubernetes-pod-monitor .
COPY --from=build /kubernetes-pod-monitor/config ./config

EXPOSE 80

CMD ./kubernetes-pod-monitor
