FROM ethereum/client-go:v1.9.16@sha256:549952077ecc5054408402c45ee9de44abc6c6ca3a37c6273c5b7b6dacd6971b

WORKDIR /gethnet
COPY datadir datadir/
COPY config.toml .

ENTRYPOINT [ "geth" ]
CMD [\
"--dev",\
"--datadir", "/gethnet/datadir", \
"--mine", \
"--ipcdisable", \
"--dev.period", "2", \
"--unlock", "0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f", \
"--password", "/run/secrets/node_password", \
"--config", "/gethnet/config.toml" \
]