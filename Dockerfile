FROM golang:bookworm

ADD . /opt/wallet/src/
RUN cd /opt/wallet/src/ && go build -o /opt/wallet/wallet-server ./cmd/wallet/main.go 
RUN cp /opt/wallet/src/cmd/wallet/config.env /opt/wallet/
RUN rm -fR /opt/wallet/src
RUN useradd -m wallet
RUN chown -R wallet:wallet /opt/wallet

USER wallet

ENTRYPOINT [ "/opt/wallet/wallet-server"]
